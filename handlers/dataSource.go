package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"windy/log"
	"windy/models"

	"github.com/olivetree123/coco"
)

// DataSourceCreateParam 创建数据库参数
type DataSourceDBCreateParam struct {
	Name     string
	Host     string
	Port     int
	UserName string
	Password string
	DBName   string
}

// Validate 参数验证
func (param *DataSourceDBCreateParam) Validate() (bool, error) {
	if param.Name == "" {
		return false, errors.New("name should not be null")
	}
	if param.Host == "" {
		return false, errors.New("host should not be null")
	}
	if param.Port == 0 {
		return false, errors.New("port should not be null")
	}
	if param.UserName == "" {
		return false, errors.New("username should not be null")
	}
	if param.Password == "" {
		return false, errors.New("password should not be null")
	}
	return true, nil
}

// Load 加载参数
func (param *DataSourceDBCreateParam) Load(request *http.Request) error {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(param)
	if err != nil {
		return err
	}
	return nil
}

// DataSourceDBCreateHandler 新建数据源
func DataSourceDBCreateHandler(c *coco.Coco) coco.Result {
	var param DataSourceDBCreateParam
	param.Load(c.Request)
	status, err := param.Validate()
	if !status {
		log.Logger.Error(err)
		return coco.ErrorResponse(100)
	}
	if dataSource, _ := models.GetDataSourceDB(param.Name); dataSource != nil {
		log.Logger.Error("dataSource exists, name=", param.Name)
		return coco.ErrorResponse(100)
	}
	dataSourceDB, err := models.NewDataSourceDB(param.Name, param.Host, param.Port, param.UserName, param.Password, param.DBName)
	if err != nil {
		log.Logger.Error(err)
		return coco.ErrorResponse(100)
	}
	return coco.APIResponse(dataSourceDB)
}

// DataSourceDBListHandler 获取数据源列表
func DataSourceDBListHandler(c *coco.Coco) coco.Result {
	dataSourceDBList, err := models.ListDataSourceDB()
	if err != nil {
		log.Logger.Error(err)
		return coco.ErrorResponse(100)
	}
	return coco.APIResponse(dataSourceDBList)
}
