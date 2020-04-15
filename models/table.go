package models

import "time"

type Table struct {
	BaseModel
	Name            string    `json:"name"`
	DbID            string    `json:"dbID"`
	PrimaryField    string    `json:"primaryField"`
	UpdateTimeField string    `json:"updateTimeField"`
	LastUpdateTime  time.Time `json:"lastUpdateTime"`
}

// Fields 获取 Table 的所有列
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

func GetTable(tableID string) (*Table, error) {
	var table Table
	if err := DB.First(&table, "uid = ? and status = ?", tableID, true).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

func GetTableByName(dbID string, name string) (*Table, error) {
	var table Table
	if err := DB.First(&table, "db_id = ? and name = ? and status = ?", dbID, name, true).Error; err != nil {
		return nil, err
	}
	return &table, nil
}
