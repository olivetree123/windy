package models

import (
	"errors"
	"sort"
	"windy/entity"
)

// Document 文档
type Document struct {
	BaseModel
	TableID      string `json:"tableID"`
	Content      string `json:"content"` // 文档内容。比如：{"id":1, "name":"olivetree"}
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
func SearchDocument(tableID string, fields []string, words []string, match map[string]string, page int, pageSize int) ([]Document, int, error) {
	match2 := make(map[string]string)
	for key, value := range match {
		field, err := GetField(tableID, key)
		if err != nil {
			return nil, 0, err
		}
		match2[field.UID] = value
	}
	// 1. 根据布尔模型，过滤出所有匹配到的文档
	docs, err := GetAllMatchDoc(tableID, words, fields, match2)
	if err != nil {
		return nil, 0, err
	}
	total := len(docs)
	if (page-1)*pageSize >= total {
		return nil, 0, errors.New("invalid page")
	}
	// 2. 根据TF-IDF 模型，为每个匹配到的文档打分
	// 打分规则：单词在文档中出现的次数除以该词的词频
	var ss []entity.DocumentScore
	for _, docID := range docs {
		score, err := GetScore(docID, words, fields)
		if err != nil {
			return nil, 0, err
		}
		ss = append(ss, *entity.NewDocumentScore(docID, score))
	}
	// 3. 对打分结果进行排序，并取前十名
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Score > ss[j].Score
	})
	startPoint := (page - 1) * pageSize
	endPoint := page * pageSize
	if endPoint > len(ss) {
		endPoint = len(ss)
	}
	var rs []string
	for _, s := range ss[startPoint:endPoint] {
		rs = append(rs, s.DocID)
	}
	// 4. 获取结果
	var documents []Document
	if err = DB.Where("uid in (?) and status = ?", rs, true).Find(&documents).Error; err != nil {
		return nil, 0, err
	}
	var result []Document
	for _, docID := range rs {
		for _, document := range documents {
			if document.UID == docID {
				result = append(result, document)
				break
			}
		}
	}
	return result, total, nil
}
