package handlers

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"net/http"
	"windy/log"
	"windy/models"

	"github.com/olivetree123/coco"
)

type Field struct {
	Name       string
	Type       string
	UnSearch   bool
	PrimaryKey bool
}

type TableCreateParam struct {
	DbID            string
	Name            string
	Fields          []Field
	PrimaryField    string
	UpdateTimeField string
	DataSourceDbID  string // 关联的数据源
}

func (param *TableCreateParam) Validate() (bool, error) {
	if param.Name == "" {
		return false, errors.New("name should not be null")
	}
	if param.DbID == "" {
		return false, errors.New("dbID should not be null")
	}
	if param.DataSourceDbID == "" {
		return false, errors.New("dataSourceDBID should not be null")
	}
	return true, nil
}

func (param *TableCreateParam) Load(request *http.Request) error {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(param)
	if err != nil {
		return err
	}
	return nil
}

// TableCreateHandler 创建 table 和 field
func TableCreateHandler(c *coco.Coco) coco.Result {
	var param TableCreateParam
	param.Load(c.Request)
	status, err := param.Validate()
	if !status {
		log.Logger.Error(err)
		return coco.ErrorResponse(100)
	}
	var table models.Table
	err = models.DB.Transaction(func(tx *gorm.DB) error {
		table = models.Table{
			Name:            param.Name,
			DbID:            param.DbID,
			DataSourceDbID:  param.DataSourceDbID,
			PrimaryField:    param.PrimaryField,
			UpdateTimeField: param.UpdateTimeField,
		}
		if err := tx.Create(&table).Error; err != nil {
			return err
		}
		for _, fieldParam := range param.Fields {
			field := models.Field{
				TableID:  table.UID,
				Name:     fieldParam.Name,
				Type:     fieldParam.Type,
				UnSearch: fieldParam.UnSearch,
			}
			if err := tx.Create(&field).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Logger.Error(err)
		return coco.ErrorResponse(100)
	}
	return coco.APIResponse(table)
}
