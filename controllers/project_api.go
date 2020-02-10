package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"nomni/utils/api"
	"rtc-api/models"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
)

type ProjectApiController struct {
}

func (d ProjectApiController) Init(g echoswagger.ApiGroup) {
	g.SetSecurity("Authorization")
	g.GET("", d.GetAll).
		AddParamQueryNested(SearchInput{}).
		AddParamQuery("", "name", "go-api", false).
		AddParamQuery("", "depth", "-1:all child,0:no child,1: 1 child", false).
		AddParamQuery("", "simple", "true", false)
	g.GET("/filterDbNames", d.GetDeleteDbNames).
		AddParamQueryNested(DatabaseDto{})
	g.GET("/isCirculerReference", d.IsCirculerReference).
		AddParamQuery("", "rawId", "1", true).
		AddParamQuery("", "subIds", "2,3,4", true)
	g.GET("/:id", d.GetById).
		AddParamPath("", "id", "1").
		AddParamQuery("", "depth", "-1:all child,0:no child,1: 1 child", false)
	g.POST("", d.Create).
		AddParamBody(ProjectDto{}, "project", "new project", true)
	g.PUT("/:id", d.Update).
		AddParamPath("", "id", "1").
		AddParamBody(ProjectDto{}, "project", "update project", true)
	g.DELETE("/:id", d.Delete).
		AddParamPath("", "id", "14")
}
func (d ProjectApiController) IsCirculerReference(c echo.Context) error {
	var v struct {
		RawId  int    `query:"rawId"`
		SubIds string `query:"subIds"`
	}
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.ParameterParsingError(err))
	}
	if v.RawId == 0 { // 1.add
		return ReturnApiSucc(c, http.StatusOK, false)
	}

	//2.edit
	if len(v.SubIds) ==0{
		return ReturnApiSucc(c, http.StatusOK, false)
	}
	projects, err := models.Project{}.GetAllReal(c.Request().Context())
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, err)
	}
	p := d.ProjectById(v.RawId, projects)
	if p == nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.InvalidParamError("rawId", fmt.Sprint(v.RawId), errors.New("not found")))
	}

	ids, err := d.strToInts(v.SubIds)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.InvalidParamError("subIds", fmt.Sprint(v.SubIds), err))
	}
	ps := d.ProjectByIds(ids, projects)
	boolDto := &BoolDto{}
	for _, project := range ps {
		d.isCirculerReference(p.Id, project, projects, boolDto)
		if boolDto.Flag == true {
			return ReturnApiSucc(c, http.StatusOK, true)
		}
	}
	return ReturnApiSucc(c, http.StatusOK, false)
}

func (d ProjectApiController) strToInts(str string) ([]int, error) {
	idStrs := strings.Split(str, ",")
	ids := make([]int, 0)
	for _, idStr := range idStrs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, int(id))
	}
	return ids, nil
}

func (d ProjectApiController) isCirculerReference(rawId int, project *models.Project, projects []*models.Project, boolDto *BoolDto) {
	if len(project.SubIds) != 0 {
		if models.ContainInt(project.SubIds, rawId) { //To prevent circular references
			boolDto.Flag = true
			return
		}
		subProjects := d.ProjectByIds(project.SubIds, projects)
		for k, v := range subProjects {
			if len(v.SubIds) != 0 {
				d.isCirculerReference(rawId, subProjects[k], projects, boolDto)
			}
		}
	}
}

func (d ProjectApiController) GetAll(c echo.Context) error {
	if len(c.QueryParam("name")) != 0 {
		return d.GetByName(c)
	}
	if len(c.QueryParam("simple")) != 0 {
		return d.GetAllSimple(c)
	}
	var v SearchInput
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.ParameterParsingError(err))
	}
	if v.MaxResultCount == 0 {
		v.MaxResultCount = DefaultMaxResultCount
	}
	totalCount, items, err := models.Project{}.GetAll(c.Request().Context(),
		v.Sortby, v.Order, v.SkipCount, v.MaxResultCount, v.Like)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, err)
	}
	if len(items) == 0 {
		return ReturnApiFail(c, http.StatusBadRequest, api.RtcServiceHasNotFoundError())
	}
	var depth int
	if len(c.QueryParam("depth")) != 0 {
		depthInt64, err := strconv.ParseInt(c.QueryParam("depth"), 10, 64)
		if err != nil {
			return ReturnApiFail(c, http.StatusBadRequest, api.InvalidParamError("depth", c.QueryParam("depth"), err))
		}
		depth = int(depthInt64)
	}
	var projects []*models.Project
	if depth != 0 {
		projects, err = models.Project{}.GetAllReal(c.Request().Context())
		if err != nil {
			return ReturnApiFail(c, http.StatusInternalServerError, err)
		}
	}
	imageAccount, producer, consumer, err := d.extraInfo(c.Request().Context())
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, err)
	}
	for k := range items {
		if status, err := d.getAllWithChild(c, items[k], depth, projects, imageAccount, producer, consumer); err != nil {
			return ReturnApiFail(c, status, err)
		}
	}
	return ReturnApiSucc(c, http.StatusOK, api.ArrayResult{
		TotalCount: totalCount,
		Items:      items,
	})

}

func (d ProjectApiController) GetById(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.InvalidParamError("id", c.Param("id"), err))
	}
	has, project, err := models.Project{}.GetById(c.Request().Context(), int(id))
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, err)
	}
	if has == false {
		return ReturnApiFail(c, http.StatusBadRequest, api.RtcServiceHasNotFoundError())
	}
	if status, err := d.getWithChild(c, project); err != nil {
		return ReturnApiFail(c, status, err)
	}

	return ReturnApiSucc(c, http.StatusOK, project)
}

func (d ProjectApiController) GetDeleteDbNames(c echo.Context) error {
	var v DatabaseDto
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.ParameterParsingError(err))
	}
	if err := c.Validate(v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.ParameterParsingError(err))
	}

	projects, err := models.Project{}.GetAllReal(c.Request().Context())
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, err)
	}
	tempNames := make([]string, 0)
	names := strings.Split(v.DbNames, ",")
	for _, project := range projects {
		for k, dbNames := range project.Setting.Databases {
			if k == "mysql" {
				for _, dbName := range dbNames {
					if v.TenantName == project.TenantName &&
						v.Namespace == project.Namespace &&
						v.Id != project.Id &&
						ContainsString(names, dbName) {
						tempNames = append(tempNames, dbName)
					}
				}
			}
		}
	}
	dbNames := make([]string, 0)
	for _, v := range names {
		if !ContainsString(tempNames, v) {
			dbNames = append(dbNames, v)
		}
	}
	v.DbNames = strings.Join(dbNames, ",")
	return ReturnApiSucc(c, http.StatusOK, v)
}

func (d ProjectApiController) GetAllSimple(c echo.Context) error {
	project, err := models.Project{}.GetAllSimple(c.Request().Context())
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, err)
	}
	return ReturnApiSucc(c, http.StatusOK, project)
}
func (d ProjectApiController) GetByName(c echo.Context) error {
	name := c.QueryParam("name")
	has, project, err := models.Project{}.GetByName(c.Request().Context(), name)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, err)
	}
	if has == false {
		return ReturnApiFail(c, http.StatusBadRequest, api.RtcServiceHasNotFoundError())
	}

	if status, err := d.getWithChild(c, project); err != nil {
		return ReturnApiFail(c, status, err)
	}
	return ReturnApiSucc(c, http.StatusOK, project)
}
func (ProjectApiController) ProjectByIds(ids []int, projects []*models.Project) []*models.Project {
	pFilters := make([]*models.Project, 0)
	for _, id := range ids {
		for _, p := range projects {
			if id == p.Id {
				pFilters = append(pFilters, p)
			}
		}
	}
	return pFilters
}
func (ProjectApiController) ProjectById(id int, projects []*models.Project) *models.Project {
	for _, p := range projects {
		if id == p.Id {
			return p
		}
	}
	return nil
}
func (ProjectApiController) ProjectByName(name string, projects []*models.Project) *models.Project {
	for _, p := range projects {
		if name == p.Name {
			return p
		}
	}
	return nil
}

func (ProjectApiController) ProjectByNames(names []string, projects []*models.Project) []*models.Project {
	pFilters := make([]*models.Project, 0)
	for _, name := range names {
		for _, p := range projects {
			if name == p.Name {
				pFilters = append(pFilters, p)
			}
		}
	}
	return pFilters
}

func (d ProjectApiController) loopGet(c echo.Context, rawId int, project *models.Project, projects []*models.Project, depth int) {
	if len(project.SubIds) != 0 {
		if models.ContainInt(project.SubIds, rawId) { //To prevent circular references
			return
		}
		subProjects := d.ProjectByIds(project.SubIds, projects)
		project.Children = subProjects
		if depth == 1 {
			return
		}

		for k, v := range project.Children {
			if len(v.SubIds) != 0 {
				d.loopGet(c, rawId, project.Children[k], projects, depth)
			}
		}
	}
}

func (d ProjectApiController) Create(c echo.Context) error {
	var v models.Project
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.ParameterParsingError(err))
	}
	if err := c.Validate(v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.ParameterParsingError(err))
	}
	v.Name = models.Project{}.GetName(v.TenantName, v.Namespace, v.Service)
	has, _, err := models.Project{}.GetByName(c.Request().Context(), v.Name)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, err)
	}
	if has {
		return ReturnApiFail(c, http.StatusBadRequest, api.RtcServiceHasExistError())
	}
	affectedRow, err := v.Create(c.Request().Context())
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, err)
	}
	if affectedRow == int64(0) {
		return ReturnApiFail(c, http.StatusBadRequest, api.NotCreatedError())
	}
	return ReturnApiSucc(c, http.StatusCreated, v)
}

func (d ProjectApiController) Update(c echo.Context) error {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.InvalidParamError("id", c.Param("id"), err))
	}

	var v models.Project
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.ParameterParsingError(err))
	}
	if err := c.Validate(v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.ParameterParsingError(err))
	}
	v.Id = int(id)
	v.Name = models.Project{}.GetName(v.TenantName, v.Namespace, v.Service)
	has, p1, err := models.Project{}.GetById(c.Request().Context(), v.Id)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, err)
	}
	if has == false {
		return ReturnApiFail(c, http.StatusBadRequest, api.RtcServiceHasNotExistError())
	}
	nameExist, p2, err := models.Project{}.GetByName(c.Request().Context(), v.Name)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, err)
	}
	if nameExist == true {
		name1 := models.Project{}.GetName(p1.TenantName, p1.Namespace, p1.Service)
		name2 := models.Project{}.GetName(p2.TenantName, p2.Namespace, p2.Service)
		if name1 != name2 { //Name is a unique
			return ReturnApiFail(c, http.StatusBadRequest, api.RtcServiceHasExistError())
		}
	}
	affectedRow, err := v.Update(c.Request().Context(), v.Id)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, err)
	}
	if affectedRow == int64(0) {
		return ReturnApiFail(c, http.StatusBadRequest, api.NotUpdatedError())
	}
	return ReturnApiSucc(c, http.StatusOK, v)
}

func (d ProjectApiController) getWithChild(c echo.Context, project *models.Project) (int, error) {
	var depth int
	if len(c.QueryParam("depth")) != 0 {
		depthInt64, err := strconv.ParseInt(c.QueryParam("depth"), 10, 64)
		if err != nil {
			return http.StatusBadRequest, api.InvalidParamError("depth", c.QueryParam("depth"), err)
		}
		depth = int(depthInt64)
	}
	if depth != 0 {
		items, err := models.Project{}.GetAllReal(c.Request().Context())
		if err != nil {
			return http.StatusInternalServerError, err
		}
		d.loopGet(c, project.Id, project, items, depth)
	}
	imageAccount, producer, consumer, err := d.extraInfo(c.Request().Context())
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if err := (ProjectOwner{}).Reload(c.Request().Context(), project, imageAccount, producer, consumer); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
func (d ProjectApiController) getAllWithChild(c echo.Context, project *models.Project, depth int, items []*models.Project, imageAccount []models.ImageAccount, producer *models.Project, consumer *models.Project) (int, error) {

	if depth != 0 {
		d.loopGet(c, project.Id, project, items, depth)
	}
	if err := (ProjectOwner{}).Reload(c.Request().Context(), project, imageAccount, producer, consumer); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
func (d ProjectApiController) extraInfo(ctx context.Context) ([]models.ImageAccount, *models.Project, *models.Project, error) {
	imageAccount, err := models.ImageAccount{}.GetAll(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	_, producer, err := models.Project{}.GetByName(ctx, "event-broker-kafka")
	if err != nil {
		return nil, nil, nil, err
	}
	_, consumer, err := models.Project{}.GetByName(ctx, "event-kafka-consumer")
	if err != nil {
		return nil, nil, nil, err
	}
	return imageAccount, producer, consumer, nil
}
func (d ProjectApiController) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.InvalidParamError("id", c.Param("id"), err))
	}
	has, v, err := models.Project{}.GetById(c.Request().Context(), int(id))
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, err)
	}
	if has == false {
		return ReturnApiFail(c, http.StatusBadRequest, api.RtcServiceHasNotExistError())
	}
	ids, err := models.Project{}.GetParentIds(c.Request().Context(), int(id))
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, err)
	}
	if len(ids) != 0 {
		return ReturnApiFail(c, http.StatusBadRequest, api.RtcServiceNotAllowDeleteError(ids))
	}
	affectedRow, err := models.Project{}.Delete(c.Request().Context(), int(id))
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, err)
	}
	if affectedRow == int64(0) {
		return ReturnApiFail(c, http.StatusBadRequest, api.NotDeletedError())
	}
	return ReturnApiSucc(c, http.StatusOK, v)
}
