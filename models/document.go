package models

// Document 文档
type Document struct {
	BaseModel
	DbID    string `json:"dbID"`    // 所属数据库
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


func ListDocument(dbID string) ([]Document, error) {
	var docs []Document
	if err := DB.Find(&docs, "db_id = ? and status = ?", dbID, true).Error; err != nil {
		return nil, err
	}
	return docs, nil
}