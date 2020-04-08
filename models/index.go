package models

import (
	"github.com/jinzhu/gorm"
	"windy/index"
)

// Index 索引，word 在文档中出现的位置和次数
type Index struct {
	BaseModel
	Word     string `json:"word"`     // 词语
	DbID     string `json:"dbID"`     // 所属数据库
	DocID    string `json:"docID"`    // 所属文档
	Count    int    `json:"count"`    // 出现的次数
	Position string `json:"position"` // 出现的位置，多个位置组成的列表
}

// NewIndex 新建索引
func NewIndex(dbID string, docID string, word string, position []int) (*Index, error) {
	// var position2 []string
	// for _, pos := range position {
	// 	position2 = append(position2, strconv.Itoa(pos))
	// }
	index := Index{
		Word:  word,
		DbID:  dbID,
		DocID: docID,
		Count: 1,
		// Position: strings.Join(position2, ","),
		Position: "",
	}
	err := DB.Create(&index).Error
	if err != nil {
		return nil, err
	}
	return &index, nil
}

// GetIndex 获取索引
func GetIndex(dbID string, docID string, word string) (*Index, error) {
	var index Index
	err := DB.Find(&index, "db_id = ? and doc_id = ? and word = ? and status = ?", dbID, docID, word, true).Error
	if err != nil {
		return nil, err
	}
	return &index, nil
}

func (index *Index) AddCount() error {
	index.Count++
	err := DB.Save(index).Error
	if err != nil {
		return err
	}
	return nil
}

// CreateIndexForWords 为多个 word 批量创建索引
func CreateIndexForWords(dbID string, docID string, words []index.Word) error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		for _, word := range words {
			var idx Index
			err := tx.First(&idx, "db_id = ? and doc_id = ? and word = ? and status = ?", dbID, docID, word, true).Error
			if err == nil {
				idx.Count++
				if err = tx.Save(&idx).Error; err != nil {
					return err
				}
				continue
			}
			if err != gorm.ErrRecordNotFound {
				return err
			}
			idx = Index{
				Word:     word.Content,
				DbID:     dbID,
				DocID:    docID,
				Count:    word.Count,
				Position: "",
			}
			if err = tx.Create(&idx).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err

}
