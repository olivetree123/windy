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

// DataSourceTableCreateParam 创建数据源表参数
type DataSourceTableCreateParam struct {
	DataSourceDBID  string
	Name            string   // table name
	Fields          []string // 需要进行搜索的字段
	PrimaryKey      string   // 主键字段
	IndexDBID       string   // 索引所在的库
	UpdateTimeField string
}

// Validate 参数验证
func (param *DataSourceTableCreateParam) Validate() (bool, error) {
	if param.DataSourceDBID == "" {
		return false, errors.New("dataSourceDBID should not be null")
	}
	if param.Name == "" {
		return false, errors.New("name should not be null")
	}
	if len(param.Fields) == 0 {
		return false, errors.New("fields should not be null")
	}
	if param.PrimaryKey == "" {
		return false, errors.New("primaryKey should not be null")
	}
	if param.IndexDBID == "" {
		return false, errors.New("indexDBID should not be null")
	}
	if param.UpdateTimeField == "" {
		return false, errors.New("updateTimeField should not be null")
	}
	return true, nil
}

// Load 加载参数
func (param *DataSourceTableCreateParam) Load(request *http.Request) error {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(param)
	if err != nil {
		return err
	}
	return nil
}

// DataSourceTableListParam 获取数据源表的参数
type DataSourceTableListParam struct {
	DbID string
}

// Validate 参数验证
func (param *DataSourceTableListParam) Validate() (bool, error) {
	if param.DbID == "" {
		return false, errors.New("dbID should not be null")
	}
	return true, nil
}

// Load 加载参数
func (param *DataSourceTableListParam) Load(request *http.Request) error {
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

// DataSourceTableCreateHandler 设置数据源的表
func DataSourceTableCreateHandler(c *coco.Coco) coco.Result {
	var param DataSourceTableCreateParam
	param.Load(c.Request)
	status, err := param.Validate()
	if !status {
		log.Logger.Error(err)
		return coco.ErrorResponse(100)
	}
	if _, err := models.GetDatabase(param.IndexDBID); err != nil {
		log.Logger.Error(err)
		return coco.ErrorResponse(100)
	}
	if table, _ := models.GetDataSourceTable(param.DataSourceDBID, param.Name); table != nil {
		log.Logger.Error("table already exists")
		return coco.ErrorResponse(100)
	}
	table, err := models.NewDataSourceTable(param.DataSourceDBID, param.Name, param.Fields, param.PrimaryKey, param.IndexDBID, param.UpdateTimeField)
	if err != nil {
		log.Logger.Error(err)
		return coco.ErrorResponse(100)
	}
	return coco.APIResponse(table)
}

// DataSourceTableListHandler 获取数据源表的列表
func DataSourceTableListHandler(c *coco.Coco) coco.Result {
	var param DataSourceTableListParam
	param.Load(c.Request)
	status, err := param.Validate()
	if !status {
		log.Logger.Error(err)
		return coco.ErrorResponse(100)
	}
	tables, err := models.ListDataSourceTable(param.DbID)
	if err != nil {
		log.Logger.Error(err)
		return coco.ErrorResponse(100)
	}
	return coco.APIResponse(tables)
}
