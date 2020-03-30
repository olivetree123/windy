package handlers

import (
	"windy/index"
	"windy/log"
	"windy/models"

	"github.com/olivetree123/coco"
)

// DocCreateHandler 创建文档
func DocCreateHandler(c *coco.Coco) coco.Result {
	dbID := c.Params.ByName("db_id")
	content := c.Params.ByName("content")
	if dbID == "" || content == "" {
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
