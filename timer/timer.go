package timer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
	"windy/config"
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

func parseResultToMap(columns []string, data []interface{}) map[string]interface{} {
	length := len(columns)
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
	return result
}

// Timer 定时器，每 60s 从数据库获取新的数据
func Timer() {
	duration := time.Duration(time.Second * 20)
	t := time.NewTicker(duration)
	defer t.Stop()
	for {
		<-t.C
		dbs, err := models.ListDataSourceDB()
		if err != nil {
			log.Logger.Error(err)
			continue
		}
		for _, db := range dbs {
			tables, err := models.ListTableByDataSourceDB(db.UID)
			if err != nil {
				log.Logger.Error(err)
				continue
			}
			dbURL := fmt.Sprintf("%s:%s@(%s:%d)/%s", db.UserName, db.Password, db.Host, db.Port, db.DBName)
			conn, err := gorm.Open("mysql", dbURL)
			if err != nil {
				log.Logger.Error(err)
				continue
			}
			for _, table := range tables {
				fields, err := table.Fields()
				if err != nil {
					log.Logger.Error(err)
					continue
				}
				// 无名函数，获取所有列的名称
				fieldsStr := func() string {
					var fs []string
					for _, field := range fields {
						fs = append(fs, field.Name)
					}
					return strings.Join(fs, ",")
				}()
				err = models.DB.Transaction(func(tx *gorm.DB) error {
					now := time.Now()
					sql := fmt.Sprintf("SELECT %s FROM %s where %s > '%s' and %s <= '%s'", fieldsStr, table.Name, table.UpdateTimeField, table.LastUpdateTime, table.UpdateTimeField, now)
					//log.Logger.Info(sql)
					rows, err := conn.Raw(sql).Rows()
					if err != nil {
						return err
					}
					columns, err := rows.Columns()
					if err != nil {
						return err
					}
					length := len(columns)
					if err = tx.Model(&models.Table{}).Where("uid = ? and status = ?", table.UID, true).Update("last_update_time", now).Error; err != nil {
						return err
					}
					for rows.Next() {
						data := makeResultReceiver(length)
						if err = rows.Scan(data...); err != nil {
							log.Logger.Error(err)
							return err
						}
						result := parseResultToMap(columns, data)
						// 考虑各个字段的类型
						for _, field := range fields {
							name := field.Name
							value, found := result[field.Name]
							if !found {
								continue
							}
							if field.Type == config.FieldTypeInt {
								result[name], err = strconv.Atoi(value.(string))
							} else if field.Type == config.FieldTypeBool {
								result[name], err = strconv.ParseBool(value.(string))
							} else if field.Type == config.FieldTypeFloat {
								result[name], err = strconv.ParseFloat(value.(string), 32)
							} else if field.Type == config.FieldTypeTime {
								result[name], err = time.Parse("2006-01-02 15:04:05", value.(string))
							}
							if err != nil {
								return err
							}
						}
						result2, err := json.Marshal(result)
						if err != nil {
							return err
						}
						primaryValue := result[table.PrimaryField].(string)
						content := string(result2)
						// 检查文档是否已存在
						var doc2 models.Document
						err = tx.First(&doc2, "table_id = ? and primary_value = ?", table.UID, primaryValue).Error
						if err == nil {
							// 已存在，要更新文档，删除索引，并重建索引
							if err := tx.Model(&models.Document{}).Where("uid = ? and status = ?", doc2.UID, true).Update("content", content).Error; err != nil {
								return err
							}
							if err := tx.Where("doc_id = ?", doc2.UID).Delete(&models.Index{}).Error; err != nil {
								return err
							}
							for _, field := range fields {
								if field.Name == table.PrimaryField || field.Type != config.FieldTypeString || field.UnSearch {
									continue
								}
								value, found := result[field.Name]
								if !found {
									continue
								}
								words := index.SplitWord(value.(string))
								for _, word := range words {
									idx := models.Index{
										Word:    word.Content,
										DbID:    table.DbID,
										TableID: table.UID,
										DocID:   doc2.UID,
										FieldID: field.UID,
										Count:   word.Count,
										// Position: strings.Join(position2, ","),
										Position: "",
									}
									if err := tx.Create(&idx).Error; err != nil {
										return err
									}
								}
							}
							return nil
						}
						if err != gorm.ErrRecordNotFound {
							return err
						}
						// 文档不存在，创建文档和索引
						doc := models.Document{
							TableID:      table.UID,
							Content:      content,
							PrimaryValue: primaryValue,
						}
						if err = tx.Create(&doc).Error; err != nil {
							return err
						}
						for _, field := range fields {
							if field.Name == table.PrimaryField || field.Type != config.FieldTypeString || field.UnSearch {
								continue
							}
							value, found := result[field.Name]
							if !found {
								continue
							}
							words := index.SplitWord(value.(string))
							for _, word := range words {
								idx := models.Index{
									Word:    word.Content,
									DbID:    table.DbID,
									TableID: table.UID,
									DocID:   doc.UID,
									FieldID: field.UID,
									Count:   word.Count,
									// Position: strings.Join(position2, ","),
									Position: "",
								}
								if err := tx.Create(&idx).Error; err != nil {
									return err
								}
							}
						}
					}
					return nil
				})
				if err != nil {
					log.Logger.Error(err)
				}
			}
		}
	}
}
