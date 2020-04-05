package models

import (
	"strings"
)

// DataSourceDB 数据源
type DataSourceDB struct {
	BaseModel
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UserName string `json:"username"`
	Password string `json:"-"`
	DBName   string `json:"dbName"`
}

// DataSourceTable 数据源的表
type DataSourceTable struct {
	BaseModel
	DbID   string `json:"dbId"`
	Name   string `json:"name"`
	Fields string `json:"fields"`
}

// NewDataSourceDB 新建数据源
func NewDataSourceDB(name string, host string, port int, userName string, password string, dbName string) (*DataSourceDB, error) {
	dataSource := DataSourceDB{
		Name:     name,
		Host:     host,
		Port:     port,
		UserName: userName,
		Password: password,
		DBName:   dbName,
	}
	if err := DB.Create(&dataSource).Error; err != nil {
		return nil, err
	}
	return &dataSource, nil
}

// GetDataSourceDB 获取数据源
func GetDataSourceDB(name string) (*DataSourceDB, error) {
	var dataSource DataSourceDB
	if err := DB.First(&dataSource, "where name = ? and status = ?", name, true).Error; err != nil {
		return nil, err
	}
	return &dataSource, nil
}

// ListDataSourceDB 获取所有数据源列表
func ListDataSourceDB() ([]DataSourceDB, error) {
	var dataSourceList []DataSourceDB
	if err := DB.Find(&dataSourceList, "status = ?", true).Error; err != nil {
		return nil, err
	}
	return dataSourceList, nil
}

// NewDataSourceTable 新建数据源的表
func NewDataSourceTable(dbID string, name string, fields []string) (*DataSourceTable, error) {
	fieldsStr := strings.Join(fields, ",")
	table := DataSourceTable{
		DbID:   dbID,
		Name:   name,
		Fields: fieldsStr,
	}
	if err := DB.Create(&table).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

// GetDataSourceTable 根据名称获取数据源表
func GetDataSourceTable(dbID string, name string) (*DataSourceTable, error) {
	var table DataSourceTable
	if err := DB.First(&table, "db_id = ? and name = ? and status = ?", dbID, name, true).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

// ListDataSourceTable 获取数据源表
func ListDataSourceTable(dbID string) ([]DataSourceTable, error) {
	var tables []DataSourceTable
	if err := DB.Find(&tables, "db_id = ? and status = ?", dbID, true).Error; err != nil {
		return nil, err
	}
	return tables, nil
}
