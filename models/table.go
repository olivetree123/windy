package models

import "time"

type Table struct {
	BaseModel
	Name            string    `json:"name"`
	DbID            string    `json:"dbID"`
	DataSourceDbID  string    `json:"dataSourceDbID"`
	PrimaryField    string    `json:"primaryField"`
	UpdateTimeField string    `json:"updateTimeField"`
	LastUpdateTime  time.Time `json:"lastUpdateTime"`
}

// Fields 获取 DataSourceTable 的所有字段
func (table *Table) Fields() ([]Field, error) {
	var fields []Field
	if err := DB.Find(&fields, "table_id = ? and status = ?", table.UID, true).Error; err != nil {
		return nil, err
	}
	return fields, nil
}

func ListTable(dbID string) ([]Table, error) {
	var tables []Table
	if err := DB.Find(&tables, "db_id = ? and status = ?", dbID, true).Error; err != nil {
		return tables, nil
	}
	return tables, nil
}

func ListTableByDataSourceDB(dataSourceDbID string) ([]Table, error) {
	var tables []Table
	if err := DB.Find(&tables, "data_source_db_id = ? and status = ?", dataSourceDbID, true).Error; err != nil {
		return tables, nil
	}
	return tables, nil
}

func GetTable(tableID string) (*Table, error) {
	var table Table
	if err := DB.First(&table, "uid = ? and status = ?", tableID, true).Error; err != nil {
		return nil, err
	}
	return &table, nil
}
