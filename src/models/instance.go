package models

import "ledger/src/database/mysql"

type OrmInstance struct {
	Kid        int64  `gorm:"column:kid;primaryKey"`
	InstanceId string `gorm:"column:instance_id"`
	UserId     string `gorm:"column:user_id"`
	CreateTime int64  `gorm:"column:create_time"`
}

func (OrmInstance) TableName() string {
	return SetTableName("instance")
}
func GetInstance(userId string) (OrmInstance, error) {
	var instance OrmInstance
	r := mysql.OrmDB.Model(&OrmInstance{}).Where("user_id = ?", userId).First(&instance)
	return instance, r.Error
}

func InstallInstance(data OrmInstance) error {
	return mysql.OrmDB.Create(&data).Error
}
