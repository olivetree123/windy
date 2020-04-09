package models

import (
	"windy/log"
)

// Document 文档
type Document struct {
	BaseModel
	TableID      string `json:"tableID"`
	Content      string `json:"content"` // 文档内容
	PrimaryValue string `json:"primaryValue"`
	//Format  string `json:"format"`  // 文档格式，支持 string 和 json，默认为 string。v2：只需要支持 json
}

// NewDocument 新建文档
func NewDocument(tableID string, content string) (*Document, error) {
	doc := Document{
		TableID: tableID,
		Content: content,
	}
	err := DB.Create(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// UpdateDocument 更新文档
func UpdateDocument(documentID string, content string) error {
	err := DB.Model(&Document{}).Where("uid = ?", documentID).Update("content", content).Error
	return err
}

// GetDocument 根据 ID 获取文档
func GetDocument(documentID string) (*Document, error) {
	var document Document
	if err := DB.First(&document, "uid = ? and status = ?", documentID, true).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

// ListDocument 获取文档列表
func ListDocument(dbID string) ([]Document, error) {
	var docs []Document
	if err := DB.Find(&docs, "db_id = ? and status = ?", dbID, true).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

// SearchDocument 查找索引
func SearchDocument(dbID string, tableID string, fields []string, words []string) ([]Document, error) {
	var docs []string
	sql := "select doc_id, count(1) as count1 from `index` where db_id = ? and table_id = ? and status = ? and word in (?) and field_id in (?)  group by doc_id order by count1 desc limit 10"
	rows, err := DB.Raw(sql, dbID, tableID, true, words, fields).Rows()
	//rows, err := DB.Raw("select doc_id, count(1) as count1 from `index` where db_id = ? and status = ? and word in (?) group by doc_id order by count1 desc limit 10", dbID, true, words).Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var docID string
		var count int
		if err = rows.Scan(&docID, &count); err != nil {
			return nil, err
		}
		log.Logger.Info(docID, "\t", count)
		docs = append(docs, docID)
	}
	var documents []Document
	if err = DB.Where("uid in (?) and status = ?", docs, true).Find(&documents).Error; err != nil {
		return nil, err
	}
	var result []Document
	for _, docID := range docs {
		for _, document := range documents {
			if document.UID == docID {
				result = append(result, document)
				break
			}
		}
	}
	return result, nil
}
