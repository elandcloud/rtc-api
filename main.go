package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var Version string

func main() {
	serviceName, flag := (Flag{}).Init()
	if StringPointCheck(serviceName) == false {
		return
	}

	c, err := Config{}.LoadEnv(*serviceName, flag)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = writeLocal(c); err != nil {
		fmt.Println(err)
		return
	}

	if err = (Compose{}).Exec(c); err != nil {
		fmt.Println(err)
		return
	}
}

func writeLocal(c *FullDto) (err error) {
	if scope == LOCAL.String() {
		return
	}
	if err = composeWriteYml(c); err != nil {
		return
	}
	if err = (Nginx{}).WriteConfig(c.Project, c.Port.EventBroker); err != nil {
		return
	}
	if err = (ProjectInfo{}).WriteSql(c.Project); err != nil {
		return
	}

	if err = (Config{}).WriteYml(c); err != nil {
		return
	}
	return
}

func composeWriteYml(c *FullDto) (err error) {
	viper := viper.New()
	p := ProjectInfo{}
	d := Compose{}
	if p.ShouldKafka(c.Project) {
		d.setComposeKafkaEland(viper, c.Port.Kafka, c.Port.KafkaSecond, c.Port.Zookeeper, c.Ip)
	}
	if p.ShouldDb(c.Project, MYSQL) {
		d.setComposeMysql(viper, c.Port.Mysql)
	}
	if p.ShouldDb(c.Project, REDIS) {
		d.setComposeRedis(viper, c.Port.Redis)
	}

	if p.ShouldEventBroker(c.Project) {
		streamNames := p.StreamList(c.Project)
		EventBroker{}.SetEventBroker(viper, c.Port.EventBroker, streamNames)
	}
	d.setComposeApp(viper, c.Project)
	d.setComposeNginx(viper, c.Project.ServiceName, c.Port.Nginx)
	d.WriteYml(viper)
	return
}
