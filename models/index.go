package models

import (
	"strconv"
	"strings"
)

// Index 索引，word 在文档中出现的位置和次数
type Index struct {
	BaseModel
	Word     string `json:"word"`     // 词语
	DbID     string `json:"db_id"`    // 所属数据库
	DocID    string `json:"doc_id"`   // 所属文档
	Count    int    `json:"count"`    // 出现的次数
	Position string `json:"position"` // 出现的位置，多个位置组成的列表
}

// NewIndex 新建索引
func NewIndex(dbID string, docID string, word string, count int, position []int) (*Index, error) {
	var position2 []string
	for _, pos := range position {
		position2 = append(position2, strconv.Itoa(pos))
	}
	index := Index{
		Word:     word,
		DbID:     dbID,
		DocID:    docID,
		Count:    count,
		Position: strings.Join(position2, ","),
	}
	err := DB.Create(&index).Error
	if err != nil {
		return nil, err
	}
	return &index, nil
}

// FindIndex 查找索引
func FindIndex(dbID string, words []string) ([]Index, error) {
	var indexes []Index
	query := DB.Where("db_id = ? and status = ?", dbID, true)
	if len(words) > 0 {
		query = query.Where("word in (?)", words)
	}
	err := query.Find(&indexes).Error
	if err != nil {
		return nil, err
	}
	return indexes, nil
}
