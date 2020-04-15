package models

// DataSourceDB 数据源，目前仅支持 MySQL
type DataSourceDB struct {
	BaseModel
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UserName string `json:"username"`
	Password string `json:"-"`
	DBName   string `json:"dbName"`
}

// NewDataSourceDB 新建数据源
func NewDataSourceDB(name string, host string, port int, userName string, password string, dbName string) (*DataSourceDB, error) {
	dataSource := DataSourceDB{
		Name:     name,
		Host:     host,
		Port:     port,
		UserName: userName,
		Password: password,
		DBName:   dbName,
	}
	if err := DB.Create(&dataSource).Error; err != nil {
		return nil, err
	}
	return &dataSource, nil
}

func GetDataSource(uid string) (*DataSourceDB, error) {
	var dataSource DataSourceDB
	if err := DB.First(&dataSource, "uid = ? and status = ?", uid, true).Error; err != nil {
		return nil, err
	}
	return &dataSource, nil
}

// GetDataSourceByName 获取数据源
func GetDataSourceByName(name string) (*DataSourceDB, error) {
	var dataSource DataSourceDB
	if err := DB.First(&dataSource, "name = ? and status = ?", name, true).Error; err != nil {
		return nil, err
	}
	return &dataSource, nil
}

// ListDataSourceDB 获取所有数据源列表
func ListDataSourceDB() ([]DataSourceDB, error) {
	var dataSourceList []DataSourceDB
	if err := DB.Find(&dataSourceList, "status = ?", true).Error; err != nil {
		return nil, err
	}
	return dataSourceList, nil
}
