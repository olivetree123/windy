package timer

import (
	"encoding/json"
	"fmt"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
	"windy/index"
	"windy/log"
	"windy/models"
)

func makeResultReceiver(length int) []interface{} {
	result := make([]interface{}, 0, length)
	for i := 0; i < length; i++ {
		var current interface{}
		current = struct{}{}
		result = append(result, &current)
	}
	return result
}

// Timer 定时器，每 60s 从数据库获取新的数据
func Timer() {
	duration := time.Duration(time.Second * 20)
	t := time.NewTicker(duration)
	defer t.Stop()
	for {
		log.Logger.Info("timer...")
		<-t.C
		dbs, err := models.ListDataSourceDB()
		if err != nil {
			log.Logger.Error(err)
			continue
		}
		for _, db := range dbs {
			tables, err := models.ListDataSourceTable(db.UID)
			if err != nil {
				log.Logger.Error(err)
				continue
			}
			dbURL := fmt.Sprintf("%s:%s@(%s:%d)/%s", db.UserName, db.Password, db.Host, db.Port, db.DBName)
			conn, err := gorm.Open("mysql", dbURL)
			//conn, err := sql.Open("mysql", dbURL)
			if err != nil {
				log.Logger.Error(err)
				continue
			}
			for _, table := range tables {
				sql := fmt.Sprintf("SELECT %s FROM %s order by create_time desc limit 20", table.Fields, table.Name)
				log.Logger.Info(sql)
				rows, err := conn.Raw(sql).Rows()
				if err != nil {
					log.Logger.Error(err)
					continue
				}
				columns, err := rows.Columns()
				if err != nil {
					log.Logger.Error(err)
					continue
				}
				length := len(columns)
				for rows.Next() {
					data := makeResultReceiver(length)
					if err = rows.Scan(data...); err != nil {
						log.Logger.Error(err)
						continue
					}
					result := make(map[string]interface{})
					for i := 0; i < length; i++ {
						k := columns[i]
						v := *(data[i]).(*interface{})
						result[k] = v
					}
					for key, value := range result {
						vType := reflect.TypeOf(value)
						switch vType.String() {
						case "[]uint8":
							v := string(value.([]uint8))
							result[key] = v
						default:
							log.Logger.Info("type = ", vType.String())
						}
					}
					result2, err := json.Marshal(result)
					if err != nil {
						log.Logger.Error(err)
						continue
					}
					primaryValue := result[table.PrimaryKey].(string)
					content := string(result2)
					// 检查文档是否已存在
					doc2, err := models.GetDataSourceDocument(table.UID, primaryValue)
					if err == nil {
						// 已存在，要更新文档，删除索引，并重建索引
						err = models.DB.Transaction(func(tx *gorm.DB) error {
							if err := tx.Model(&models.Document{}).Where("uid = ? and status = ?", doc2.DocumentID, true).Update("content", content).Error; err != nil {
								return err
							}
							if err := tx.Where("doc_id = ?", doc2.DocumentID).Delete(&models.Index{}).Error; err != nil {
								return err
							}
							for key, value := range result {
								if key == table.PrimaryKey {
									continue
								}
								words := index.SplitWord(value.(string))
								for _, word := range words {
									idx := models.Index{
										Word:  word,
										DbID:  table.IndexDBID,
										DocID: doc2.DocumentID,
										Count: 1,
										// Position: strings.Join(position2, ","),
										Position: "",
									}
									if err := tx.Create(&idx).Error; err != nil {
										return err
									}
								}
							}
							return nil
						})

						continue
					}
					if err != gorm.ErrRecordNotFound {
						log.Logger.Error(err)
						continue
					}
					// 文档不存在，创建文档和索引
					err = models.DB.Transaction(func(tx *gorm.DB) error {
						doc := models.Document{
							DbID:    db.UID,
							Content: content,
						}
						err := tx.Create(&doc).Error
						if err != nil {
							return err
						}

						for key, value := range result {
							if key == table.PrimaryKey {
								continue
							}
							words := index.SplitWord(value.(string))
							for _, word := range words {
								idx := models.Index{
									Word:  word,
									DbID:  table.IndexDBID,
									DocID: doc.UID,
									Count: 1,
									// Position: strings.Join(position2, ","),
									Position: "",
								}
								if err := tx.Create(&idx).Error; err != nil {
									return err
								}
							}
						}

						doc2 := models.DataSourceDocument{
							DataSourceTableID: table.UID,
							PrimaryValue:      primaryValue,
							DocumentID:        doc.UID,
						}
						if err := tx.Create(&doc2).Error; err != nil {
							return err
						}
						return nil
					})
					if err != nil {
						log.Logger.Error(err)
						continue
					}
				}
			}
		}
	}
}
