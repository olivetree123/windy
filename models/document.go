package models

// Document 文档
type Document struct {
	BaseModel
	DbID    string `json:"db_id"`   // 所属数据库
	Content string `json:"content"` // 文档内容
}

// NewDocument 新建文档
func NewDocument(dbID string, content string) (*Document, error) {
	doc := Document{
		DbID:    dbID,
		Content: content,
	}
	err := DB.Create(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}
