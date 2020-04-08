package models

import (
	"strings"
	"time"
)

// DataSourceDB 数据源，目前仅支持 MySQL
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
	DataSourceDBID  string    `json:"dataSourceDBID"`
	Name            string    `json:"name"`
	Fields          string    `json:"fields"`
	PrimaryKey      string    `json:"primaryKey"`
	IndexDBID       string    `json:"indexDBID"`
	UpdateTimeField string    `json:"updateTimeField"`
	LastUpdateTime  time.Time `json:"lastUpdateTime"`
}

// DataSourceDocument 数据源表与 document 的关系
type DataSourceDocument struct {
	BaseModel
	DataSourceTableID string `json:"dataSourceTableID"`
	PrimaryValue      string `json:"primaryValue"`
	DocumentID        string `json:"documentId"`
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
	if err := DB.First(&dataSource, "name = ? and status = ?", name, true).Error; err != nil {
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
func NewDataSourceTable(dataSourceDBID string, name string, fields []string, primaryKey string, indexDBID string, updateTimeField string) (*DataSourceTable, error) {
	fieldsStr := strings.Join(fields, ",")
	table := DataSourceTable{
		DataSourceDBID:  dataSourceDBID,
		Name:            name,
		Fields:          fieldsStr,
		PrimaryKey:      primaryKey,
		IndexDBID:       indexDBID,
		UpdateTimeField: updateTimeField,
	}
	if err := DB.Create(&table).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

// GetDataSourceTable 根据名称获取数据源表
func GetDataSourceTable(dbID string, name string) (*DataSourceTable, error) {
	var table DataSourceTable
	if err := DB.First(&table, "data_source_db_id = ? and name = ? and status = ?", dbID, name, true).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

// ListDataSourceTable 获取数据源表
func ListDataSourceTable(dbID string) ([]DataSourceTable, error) {
	var tables []DataSourceTable
	if err := DB.Find(&tables, "data_source_db_id = ? and status = ?", dbID, true).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

// NewDataSourceDocument 新建数据源文档
func NewDataSourceDocument(dataSourceTableID string, documentID string, primaryValue string) (*DataSourceDocument, error) {
	doc := DataSourceDocument{
		DataSourceTableID: dataSourceTableID,
		PrimaryValue:      primaryValue,
		DocumentID:        documentID,
	}
	if err := DB.Create(&doc).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

// GetDataSourceDocument 根据主键的值获取数据源文档
func GetDataSourceDocument(dataSourceTableID string, primaryValue string) (*DataSourceDocument, error) {
	var doc DataSourceDocument
	if err := DB.First(&doc, "data_source_table_id = ? and primary_value = ?", dataSourceTableID, primaryValue).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}
