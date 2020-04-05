package main

import (
	"windy/handlers"

	"github.com/olivetree123/coco"
)

func main() {
	// index.LoadIndex()
	c := coco.NewCoco()
	c.AddRouter("POST", "/windy/db/create", handlers.DBCreateHandler)
	c.AddRouter("GET", "/windy/db/list", handlers.DBListHandler)
	c.AddRouter("POST", "/windy/doc/create", handlers.DocCreateHandler)
	c.AddRouter("GET", "/windy/doc/list", handlers.DocListHandler)
	c.AddRouter("GET", "/windy/doc/search", handlers.DocSearchHandler)
	c.AddRouter("POST", "/windy/datasource/db/create", handlers.DataSourceDBCreateHandler)
	c.AddRouter("GET", "/windy/datasource/db/list", handlers.DataSourceDBListHandler)
	c.AddRouter("POST", "/windy/datasource/table/create", handlers.DataSourceTableCreateHandler)
	c.AddRouter("GET", "/windy/datasource/table/list", handlers.DataSourceTableListHandler)
	c.Run("0.0.0.0", 5000)
}
