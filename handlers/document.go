package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"windy/entity"
	"windy/index"
	"windy/log"
	"windy/models"

	"github.com/olivetree123/coco"
)

// DocCreateParam 创建文档参数
type DocCreateParam struct {
	DbID    string
	Content string
}

type DocSearchParam struct {
	//DbID    string
	//TableID string
	DB       string
	Table    string
	Fields   []string
	Query    string
	Match    map[string]string
	Page     int
	PageSize int
}

// Validate 参数验证
func (param *DocCreateParam) Validate() (bool, error) {
	if param.DbID == "" {
		return false, errors.New("db_id should not be null")
	}
	if param.Content == "" {
		return false, errors.New("content should not be null")
	}
	return true, nil
}

// Load 加载参数
func (param *DocCreateParam) Load(request *http.Request) error {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(param)
	if err != nil {
		return err
	}
	return nil
}

// Validate 参数验证
func (param *DocSearchParam) Validate() (bool, error) {
	if param.DB == "" {
		return false, errors.New("db should not be null")
	}
	return true, nil
}

// Load 加载参数
func (param *DocSearchParam) Load(request *http.Request) error {
	param.Match = make(map[string]string)
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(param)
	if err != nil {
		return err
	}
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	return nil
}

// DocCreateHandler 创建文档
func DocCreateHandler(c *coco.Coco) coco.Result {
	var param DocCreateParam
	param.Load(c.Request)
	status, err := param.Validate()
	if !status {
		log.Logger.Error(err)
		return coco.ErrorResponse(100)
	}
	if _, err := models.GetDatabase(param.DbID); err != nil {
		log.Logger.Error(err)
		return coco.ErrorResponse(100)
	}
	// TODO: 相同的内容是否允许重复写入？
	// 这里需要用事务
	doc, err := models.NewDocument(param.DbID, param.Content)
	if err != nil {
		log.Logger.Error(err)
		return coco.ErrorResponse(100)
	}
	words := index.SplitWord(param.Content)
	if err = models.CreateIndexForWords(param.DbID, doc.UID, words); err != nil {
		log.Logger.Error(err)
		return coco.ErrorResponse(100)
	}
	return coco.APIResponse(doc)
}

// DocListHandler 获取文档列表，可以根据关键字进行过滤
func DocListHandler(c *coco.Coco) coco.Result {
	params := c.Request.URL.Query()
	dbID := params.Get("dbID")
	if dbID == "" {
		log.Logger.Error("dbID should not be null")
		return coco.ErrorResponse(100)
	}
	docs, err := models.ListDocument(dbID)
	if err != nil {
		log.Logger.Error(err)
		return coco.ErrorResponse(100)
	}
	return coco.APIResponse(docs)
}

// DocSearchHandler 搜索文档
func DocSearchHandler(c *coco.Coco) coco.Result {
	var param DocSearchParam
	param.Load(c.Request)
	status, err := param.Validate()
	if !status {
		log.Logger.Error(err)
		return coco.ErrorResponse(100)
	}
	db, err := models.GetDatabaseByName(param.DB)
	if err != nil {
		log.Logger.Error(err)
		return coco.ErrorResponse(100)
	}
	table, err := models.GetTableByName(db.UID, param.Table)
	if err != nil {
		log.Logger.Error(err)
		return coco.ErrorResponse(100)
	}
	words := index.SplitWord(param.Query)
	var ws []string
	for _, word := range words {
		ws = append(ws, word.Value)
	}
	var fields []string
	for _, name := range param.Fields {
		field, err := models.GetField(table.UID, name)
		if err != nil {
			log.Logger.Error(err)
			return coco.ErrorResponse(100)
		}
		fields = append(fields, field.UID)
	}
	if len(fields) == 0 {
		fs, err := models.ListField(table.UID)
		if err != nil {
			log.Logger.Error(err)
			return coco.ErrorResponse(100)
		}
		for _, field := range fs {
			fields = append(fields, field.UID)
		}
	}
	documents, total, err := models.SearchDocument(table.UID, fields, ws, param.Match, param.Page, param.PageSize)
	if err != nil {
		log.Logger.Info(err)
		return coco.ErrorResponse(100)
	}
	var data []*entity.DocumentEntity
	for _, doc := range documents {
		r, err := entity.NewDocumentEntity(doc.UID, doc.Content, doc.CreatedAt, doc.UpdatedAt)
		if err != nil {
			log.Logger.Info(err)
			return coco.ErrorResponse(100)
		}
		data = append(data, r)
	}
	result := entity.NewSearchResult(data, param.Page, param.PageSize, total)
	return coco.APIResponse(result)
}
