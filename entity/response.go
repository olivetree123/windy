package entity

import (
	"encoding/json"
	"time"
)

// DocumentEntity 文档
type DocumentEntity struct {
	UID       string      `json:"uid"`
	Content   interface{} `json:"content"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

// IndexEntity 索引
type IndexEntity struct {
	UID       string          `json:"uid"`
	Word      string          `json:"word"`
	Document  *DocumentEntity `json:"document"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

func NewDocumentEntity(uid string, content string, createdAt time.Time, updatedAt time.Time) (*DocumentEntity, error) {
	var docContent interface{}
	if err := json.Unmarshal([]byte(content), &docContent); err != nil {
		return nil, err
	}
	document := DocumentEntity{
		UID:       uid,
		Content:   docContent,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	return &document, nil
}

//
//func NewIndexEntity(uid string, word string, documentID string, createdAt time.Time, updatedAt time.Time) (*IndexEntity, error) {
//	doc, err := models.GetDocument(documentID)
//	if err != nil {
//		return nil, err
//	}
//	documentEntity, err := NewDocumentEntity(doc.UID, doc.Content, doc.CreatedAt, doc.UpdatedAt)
//	if err != nil {
//		return nil, err
//	}
//	idx := IndexEntity{
//		UID:       uid,
//		Word:      word,
//		Document:  documentEntity,
//		CreatedAt: createdAt,
//		UpdatedAt: updatedAt,
//	}
//	return &idx, nil
//}
