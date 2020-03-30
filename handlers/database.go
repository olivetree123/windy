package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"windy/log"
	"windy/models"

	"github.com/olivetree123/coco"
)

// Param 参数基类
type Param interface {
	Load(request *http.Request)
	Validate()
}

// DBCreateParam 创建数据库参数
type DBCreateParam struct {
	Name string
}

// Validate 参数验证
func (param *DBCreateParam) Validate() (bool, error) {
	if param.Name == "" {
		return false, errors.New("name should not be null")
	}
	return true, nil
}

// Load 加载参数
func (param *DBCreateParam) Load(request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(param)
	if err != nil {
		log.Error(err)
		return coco.ErrorResponse(100)
	}
}

// DBCreateHandler 创建数据库
func DBCreateHandler(c *coco.Coco) coco.Result {
	var param DBCreateParam
	param.Load(c.Request)
	status, err := param.Validate()
	if !status {
		log.Error(err)
		return coco.ErrorResponse(100)
	}
	db, err := models.NewDatabase(param.Name)
	if err != nil {
		log.Error(err)
		return coco.ErrorResponse(100)
	}
	return coco.APIResponse(db)
}

// DBListHandler 数据库列表
func DBListHandler(c *coco.Coco) coco.Result {
	dbs, err := models.ListDatabase()
	if err != nil {
		log.Error(err)
		return coco.ErrorResponse(100)
	}
	return coco.APIResponse(dbs)
}
