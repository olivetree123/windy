package models

type Field struct {
	BaseModel
	TableID  string `json:"tableID"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	UnSearch bool   `json:"unSearch"`
}

func NewField(tableID string, name string, tp string, unSearch bool) (*Field, error) {
	field := Field{
		TableID:  tableID,
		Name:     name,
		Type:     tp,
		UnSearch: unSearch,
	}
	if err := DB.Create(&field).Error; err != nil {
		return nil, err
	}
	return &field, nil
}

func GetField(tableID string, name string) (*Field, error) {
	var field Field
	if err := DB.First(&field, "table_id = ? and name = ? and status = ?", tableID, name, true).Error; err != nil {
		return nil, err
	}
	return &field, nil
}

func ListField(tableID string) ([]Field, error) {
	var fields []Field
	if err := DB.Find(&fields, "table_id = ? and status = ?", tableID, true).Error; err != nil {
		return nil, err
	}
	return fields, nil
}
