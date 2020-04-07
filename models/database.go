package models

// Database 数据库
type Database struct {
	BaseModel
	Name string `json:"name"`
}

// NewDatabase 新建数据库
func NewDatabase(name string) (*Database, error) {
	db := Database{
		Name: name,
	}
	err := DB.Create(&db).Error
	if err != nil {
		return nil, err
	}
	return &db, nil
}

// GetDatabase 获取数据库
func GetDatabase(uid string) (*Database, error) {
	var database Database
	if err := DB.First(&database, "uid = ? and status = ?", uid, true).Error; err != nil {
		return nil, err
	}
	return &database, nil
}

// ListDatabase 数据库列表
func ListDatabase() ([]Database, error) {
	var dbs []Database
	err := DB.Find(&dbs, "status = ?", 1).Error
	if err != nil {
		return nil, err
	}
	return dbs, nil
}
