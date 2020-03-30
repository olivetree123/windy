package main

import (
	"windy/handlers"

	"github.com/olivetree123/coco"
)

func main() {
	// index.LoadIndex()
	c := coco.NewCoco()
	c.AddRouter("POST", "/windy/db", handlers.DBCreateHandler)
	c.AddRouter("GET", "/windy/db/list", handlers.DBCreateHandler)
	c.AddRouter("POST", "/windy/doc", handlers.DocCreateHandler)
	c.AddRouter("POST", "/windy/doc/list", handlers.DocListHandler)
	c.Run("0.0.0.0", 5000)
}
