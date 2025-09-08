package models

type OrmInstance struct {
	Kid        int64  `gorm:"column:kid;primaryKey"`
	InstanceId string `gorm:"column:instance_id"`
	UserId     string `gorm:"column:user_id"`
	CreateTime int64  `gorm:"column:create_time"`
}

func (OrmInstance) TableName() string {
	return TableNameInstance
}
func GetInstance(userId string) (OrmInstance, error) {
	var instance OrmInstance
	r := db.Model(&OrmInstance{}).Where("user_id = ?", userId).First(&instance)
	return instance, r.Error
}

func InstallInstance(data OrmInstance) error {
	return db.Create(&data).Error
}
