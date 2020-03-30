package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
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
func (param *DocCreateParam) Load(request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(param)
	if err != nil {
		log.Error(err)
		return coco.ErrorResponse(100)
	}
}

// DocCreateHandler 创建文档
func DocCreateHandler(c *coco.Coco) coco.Result {
	var param DocCreateParam
	param.Load(c.Request)
	status, err := param.Validate()
	if !status {
		log.Error(err)
		return coco.ErrorResponse(100)
	}
	doc, err := models.NewDocument(dbID, content)
	if err != nil {
		log.Info(err)
		return coco.ErrorResponse(100)
	}
	return coco.APIResponse(doc)
}

// DocListHandler 获取文档列表，可以根据关键字进行过滤
func DocListHandler(c *coco.Coco) coco.Result {
	dbID := c.Params.ByName("db_id")
	query := c.Params.ByName("query")
	var words []string
	if query != "" {
		words = index.SplitWord(query)
	}
	indexes, err := models.FindIndex(dbID, words)
	if err != nil {
		log.Info(err)
		return coco.ErrorResponse(100)
	}
	return coco.APIResponse(indexes)
}
