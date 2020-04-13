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
	TableID  string `json:"tableID"`  // 所属的表
	DocID    string `json:"docID"`    // 所属文档
	FieldID  string `json:"fieldID"`  // word 所属的字段
	Count    int    `json:"count"`    // 出现的次数
	Position string `json:"position"` // 出现的位置，多个位置组成的列表
}

// NewIndex 新建索引
func NewIndex(dbID string, tableID string, docID string, fieldID string, word string, position []int) (*Index, error) {
	// var position2 []string
	// for _, pos := range position {
	// 	position2 = append(position2, strconv.Itoa(pos))
	// }
	index := Index{
		Word:    word,
		DbID:    dbID,
		TableID: tableID,
		DocID:   docID,
		FieldID: fieldID,
		Count:   1,
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

// GetAllMatchDoc 获取所有匹配的文档
func GetAllMatchDoc(dbID string, tableID string, words []string, fields []string) ([]string, error) {
	var docs []string
	sql := "select doc_id, count(1) as count1 from `index` where db_id = ? and table_id = ? and status = ? and word in (?) and field_id in (?)  group by doc_id"
	rows, err := DB.Raw(sql, dbID, tableID, true, words, fields).Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var docID string
		var count int
		if err = rows.Scan(&docID, &count); err != nil {
			return nil, err
		}
		docs = append(docs, docID)
	}
	return docs, nil
}

// GetScore 给文档的匹配度打分
func GetScore(documentID string, words []string, fields []string) (float32, error) {
	var indexes []Index
	if err := DB.Find(&indexes, "doc_id = ? and word in (?) and field_id in (?)", documentID, words, fields).Error; err != nil {
		return 0, err
	}
	var scoreFinal float32
	for _, idx := range indexes {
		score := float32(idx.Count) / float32(index.GetWordFreq(idx.Word))
		scoreFinal += score
	}
	return scoreFinal, nil
}
