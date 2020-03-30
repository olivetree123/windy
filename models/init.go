package models

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/twinj/uuid"
)

// DB 数据库连接对象
var DB *gorm.DB

// DBDir 数据库文件位置
var DBDir = "/Users/gao/windy/"

// BaseModel 基础模型
type BaseModel struct {
	UID       uint      `gorm:"primary_key" json:"uid"`
	Status    bool      `gorm:"default:true" json:"status"` // 0 删除，1 正常
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// BeforeCreate 写前处理
func (model *BaseModel) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("UID", uuid.NewV4().String())
	return nil
}

// Remove 移除
func (model *BaseModel) Remove() bool {
	model.Status = false
	return true
}

// InitDB 数据库初始化
func InitDB() {
	_, err := os.Stat(DBDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(DBDir, 0755)
		if err != nil {
			panic(err)
		}
	}
	dbPath := filepath.Join(DBDir, "windy.db")
	DB, err = gorm.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	DB.SingularTable(true)
	if !DB.HasTable(Database{}) {
		DB.CreateTable(Database{})
	}
	if !DB.HasTable(Document{}) {
		DB.CreateTable(Document{})
	}
	if !DB.HasTable(Index{}) {
		DB.CreateTable(Index{})
	}
}

func init() {
	InitDB()
}
