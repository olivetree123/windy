package timer

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
	"windy/log"
	"windy/models"
)

// Timer 定时器，每 60s 从数据库获取新的数据
func Timer() {
	duration := time.Duration(time.Second * 60)
	t := time.NewTicker(duration)
	defer t.Stop()
	for {
		<-t.C
		dbs, err := models.ListDataSourceDB()
		if err != nil {
			log.Error(err)
			continue
		}
		for _, db := range dbs {
			tables, err := models.ListDataSourceTable(db.UID)
			if err != nil {
				log.Error(err)
				continue
			}
			dbURL := fmt.Sprintf("%s:%s@(%s:%d)/%s", db.UserName, db.Password, db.Host, db.Port, db.DBName)
			conn, err := gorm.Open("mysql", dbURL)
			if err != nil {
				log.Error(err)
				continue
			}
			for _, table := range tables {
				fields := "(" + table.Fields + ")"
				sql := fmt.Sprintf("SLECT %s FROM %s where updated > %s", fields, table.Name, "2020-04-01")
				rows, err := conn.Raw(sql).Rows()
				if err != nil {
					log.Error(err)
					continue
				}

			}
		}
	}
}
