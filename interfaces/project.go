package interfaces

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/application"
	"github.com/opensourceways/xihe-server/domain/entity"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Users struct defines the dependencies that will be used
type Project struct {
	app *application.ProjectApp
}

//Users constructor
func NewProject(repo repository.ProjectRepository) *Project {
	app := application.NewProjectAPP(repo)
	return &Project{
		app: app,
	}
}

// @Summary Save
// @Description Save
// @Tags  Project
// @Param	body		body 	entity.Project	true		"email username phone"
// @Accept json
// @Produce json
// @Router /v1/project/save [post]
func (p *Project) Save(c *gin.Context) {
	var item entity.Project
	err := c.ShouldBindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ExportData(http.StatusBadRequest, "error input data ", err.Error()))
		return
	}
	result, err := p.app.Save(&item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ExportData(http.StatusInternalServerError, "save error ", err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.ExportData(util.CodeStatusNormal, "ok", result))
}

// @Summary Update
// @Description Update
// @Tags  Project
// @Param	id		path 	string	true		"  Project id "
// @Param	body		body 	entity.Project	true		"email username phone"
// @Accept json
// @Produce json
// @Router /v1/project/update/{id} [put]
func (p *Project) Update(c *gin.Context) {
	var item entity.Project
	err := c.ShouldBindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ExportData(http.StatusBadRequest, "error input data ", err.Error()))
		return
	}
	idstr := c.Param("id")
	docID, err := primitive.ObjectIDFromHex(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ExportData(http.StatusBadRequest, "ObjectIDFromHex error ", err.Error()))
		return
	}
	idFilter := bson.M{}
	idFilter["_id"] = docID
	result, err := p.app.Update(idFilter, &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ExportData(http.StatusInternalServerError, "Update error ", err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.ExportData(util.CodeStatusNormal, "ok", result))
}

// @Summary Query
// @Description Query
// @Tags  Project
// @Param	title		query 	string	false		"title"
// @Param	offset		query 	int	false		"  offset  default :0"
// @Param	limit		query 	int	false		"  limit , default :10 "
// @Param	order		query 	string	false		"  order by create time, default :desc "
// @Accept json
// @Produce json
// @Router /v1/project/query [get]
func (p *Project) Query(c *gin.Context) {
	idFilter := bson.M{}
	title := c.Query("title")
	if len(title) > 0 {
		idFilter["title"] = title
	}
	order := c.Query("order")
	if len(order) > 0 {
		order = "desc"
	}
	offset, _ := strconv.ParseInt(c.Query("offset"), 10, 64)
	limit, _ := strconv.ParseInt(c.Query("limit"), 10, 64)
	if limit <= 0 {
		limit = 10
	}
	result, err := p.app.Query(idFilter, offset, limit, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ExportData(http.StatusInternalServerError, "Query error ", err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.ExportData(util.CodeStatusNormal, "ok", result))
}

// @Summary GetSingleOne
// @Description GetSingleOne
// @Tags  Project
// @Param	id		path 	string	false		"id"
// @Accept json
// @Produce json
// @Router /v1/project/getSingleOne/{id} [get]
func (p *Project) GetSingleOne(c *gin.Context) {
	idFilter := bson.M{}
	idStr := c.Param("id")
	docID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ExportData(http.StatusBadRequest, "id error ", err.Error()))
		return
	}
	idFilter["_id"] = docID
	result, err := p.app.Get(idFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ExportData(http.StatusInternalServerError, "Get error ", err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.ExportData(http.StatusOK, "ok", result))

}
