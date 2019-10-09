package cmd

type ProjectOwner struct {
}

func (d ProjectOwner) ReLoad(p *Project) error {
	p.Owner.IsKafka = d.ShouldKafka(p)
	p.Owner.IsMysql = d.ShouldDb(p, MYSQL)
	p.Owner.IsSqlServer = d.ShouldDb(p, SQLSERVER)
	p.Owner.IsRedis = d.ShouldDb(p, REDIS)
	list := d.Database(p)
	p.Owner.DbTypes = d.DatabaseList(list)
	d.SetNames(p)
	d.SetDependLoop(p)
	d.SetStreams(p)
	p.Owner.IsStream = d.ShouldStream(p.Owner.StreamNames)
	if err := d.SetEvent(p); err != nil {
		return err
	}
	if err := d.SetDbAccount(p, list); err != nil {
		return err
	}
	if err := d.SetImageAccount(p); err != nil {
		return err
	}
	return nil
}
func (d ProjectOwner) SetImageAccount(p *Project) error {
	var err error
	p.Owner.ImageAccounts, err = Project{}.GetImageAccount()
	return err
}
func (d ProjectOwner) SetDbAccount(p *Project, list map[string][]string) error {
	var err error
	if p.Owner.IsMysql {
		p.Owner.MysqlAccount, err = Project{}.GetDbAccount(MYSQL)
		if err != nil {
			return err
		}
		p.Owner.MysqlAccount.DbNames = d.GetDbNameByType(MYSQL, list)
	}
	if p.Owner.IsSqlServer {
		p.Owner.SqlServerAccount, err = Project{}.GetDbAccount(SQLSERVER)
		if err != nil {
			return err
		}
		p.Owner.SqlServerAccount.DbNames = d.GetDbNameByType(SQLSERVER, list)
	}
	return nil
}
func (ProjectOwner) GetDbNameByType(dbType DateBaseType, list map[string][]string) []string {
	for k, v := range list {
		if dbType.String() == k {
			return v
		}
	}
	return nil
}
func (d ProjectOwner) SetEvent(p *Project) error {
	if p.Owner.IsStream {
		var err error
		p.Owner.EventProducer, err = Project{}.GetEventProducer()
		if err !=nil{
			return err
		}
		p.Owner.EventConsumer, err = Project{}.GetEventConsumer()
		if err !=nil{
			return err
		}
	}
	return nil
}

func (d ProjectOwner) ShouldKafka(p *Project) bool {
	flag := false
	d.kafka([]*Project{p}, &flag)
	return flag
}

func (d ProjectOwner) ShouldDb(p *Project, dbType DateBaseType) bool {
	list := d.Database(p)
	for k := range list {
		if dbType.String() == k {
			return true
		}
	}
	return false
}

func (d ProjectOwner) ShouldStream(streams []string) bool {
	if len(streams) != 0 {
		return true
	}
	return false
}

func (d ProjectOwner) DatabaseList(list map[string][]string) []string {
	dbTypes := make([]string, 0)
	for k := range list {
		dbTypes = append(dbTypes, k)
	}
	return dbTypes
}

func (d ProjectOwner) Database(p *Project) (list map[string][]string) {
	list = make(map[string][]string, 0)
	projects := []*Project{p}
	d.database(list, projects)
	for k, v := range list {
		list[k] = Unique(v)
	}
	return
}

func (d ProjectOwner) database(list map[string][]string, projects []*Project) {
	for _, project := range projects {
		for k, v := range project.Setting.Databases {
			if _, ok := list[k]; ok {
				list[k] = append(list[k], v...)
			} else {
				list[k] = v
			}
		}
		if len(project.Children) != 0 {
			d.database(list, project.Children)
		}
	}
	return
}

func (d ProjectOwner) kafka(projects []*Project, isKafka *bool) {
	for _, project := range projects {
		if project.Setting.IsProjectKafka {
			*isKafka = true
			return
		}
		if len(project.Children) != 0 {
			d.kafka(project.Children, isKafka)
		}
	}
	return
}

func (d ProjectOwner) SetNames(p *Project) {
	d.childNames(p, p)
}

func (d ProjectOwner) childNames(p *Project, pLoop *Project) {
	for _, v := range pLoop.Children {
		p.Owner.ChildNames = append(p.Owner.ChildNames, v.Name)
		if len(v.Children) != 0 {
			d.childNames(v, v)
		}
	}
}

func (d ProjectOwner) SetDependLoop(p *Project) {
	d.setDepend(p)
	d.setDependLoop(p, p)
}

func (d ProjectOwner) setDependLoop(p *Project, pLoop *Project) {
	for k, v := range pLoop.Children {
		pLoop.DependsOn = append(pLoop.DependsOn, v.Name)
		d.setDepend(pLoop.Children[k])
		if len(v.Children) != 0 {
			d.setDependLoop(pLoop.Children[k], pLoop.Children[k])
		}
	}
}
func (d ProjectOwner) setDepend(pLoop *Project) {
	if pLoop.Setting.IsProjectKafka {
		pLoop.DependsOn = append(pLoop.DependsOn, "kafka")
	}
	dbTypes := d.DatabaseList(pLoop.Setting.Databases)
	pLoop.DependsOn = append(pLoop.DependsOn, dbTypes...)
}

func (d ProjectOwner) SetStreams(p *Project) {
	d.setStreams(p, p)
	p.Owner.StreamNames = append(p.Owner.StreamNames, p.Setting.StreamNames...)
	p.Owner.StreamNames = Unique(p.Owner.StreamNames)
}

func (d ProjectOwner) setStreams(p *Project, pLoop *Project) {
	for _, v := range pLoop.Children {
		p.Owner.StreamNames = append(p.Owner.StreamNames, v.Setting.StreamNames...)
		if len(v.Children) != 0 {
			d.setStreams(v, v)
		}
	}
}