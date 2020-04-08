package entity

import (
	"encoding/json"
	"time"
	"windy/config"
	"windy/models"
)

// DocumentEntity 文档
type DocumentEntity struct {
	UID       string      `json:"uid"`
	Content   interface{} `json:"content"`
	Format    string      `json:"format"`
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

func NewDocumentEntity(uid string, content string, format string, createdAt time.Time, updatedAt time.Time) (*DocumentEntity, error) {
	var docContent interface{}
	if format == config.FormatJson {
		if err := json.Unmarshal([]byte(content), &docContent); err != nil {
			return nil, err
		}
	} else {
		docContent = content
	}
	document := DocumentEntity{
		UID:       uid,
		Format:    format,
		Content:   docContent,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	return &document, nil
}

func NewIndexEntity(uid string, word string, documentID string, createdAt time.Time, updatedAt time.Time) (*IndexEntity, error) {
	doc, err := models.GetDocument(documentID)
	if err != nil {
		return nil, err
	}
	documentEntity, err := NewDocumentEntity(doc.UID, doc.Content, doc.Format, doc.CreatedAt, doc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	idx := IndexEntity{
		UID:       uid,
		Word:      word,
		Document:  documentEntity,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	return &idx, nil
}
