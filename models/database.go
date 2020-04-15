package models

// Database 数据库
type Database struct {
	BaseModel
	Name           string `json:"name"`
	Type           int    `json:"type"`
	DataSourceDbID string `json:"dataSourceDbID"`
}

// NewDatabase 新建数据库
func NewDatabase(name string, dataSourceDbID string) (*Database, error) {
	var dbType = 1
	if dataSourceDbID != "" {
		dbType = 2
	}
	db := Database{
		Name:           name,
		Type:           dbType,
		DataSourceDbID: dataSourceDbID,
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

func GetDatabaseByName(name string) (*Database, error) {
	var database Database
	if err := DB.First(&database, "name = ? and status = ?", name, true).Error; err != nil {
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

func ListDbWithDataSource() ([]Database, error) {
	var dbs []Database
	err := DB.Find(&dbs, "status = ? and `type` = ?", 1, 2).Error
	if err != nil {
		return nil, err
	}
	return dbs, nil
}
