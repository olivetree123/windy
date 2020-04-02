package main

import (
	"windy/handlers"

	"github.com/olivetree123/coco"
)

func main() {
	// index.LoadIndex()
	c := coco.NewCoco()
	c.AddRouter("POST", "/windy/db", handlers.DBCreateHandler)
	c.AddRouter("GET", "/windy/db/list", handlers.DBListHandler)
	c.AddRouter("POST", "/windy/doc", handlers.DocCreateHandler)
	c.AddRouter("GET", "/windy/doc/list", handlers.DocListHandler)
	c.AddRouter("GET", "/windy/doc/search", handlers.DocSearchHandler)
	c.Run("0.0.0.0", 5000)
}
