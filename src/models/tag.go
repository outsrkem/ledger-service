package models

type OrmTag struct {
	Kid        int64  `gorm:"column:kid;primaryKey"` // 主键ID
	InstanceId string `gorm:"column:instance_id"`
	Name       string `gorm:"column:name"`
	Color      string `gorm:"column:color"` // 标签颜色（十六进制，如#FF0000）
	CreateTime int64  `gorm:"column:create_time"`
	UpdateTime int64  `gorm:"column:update_time"`
}

func (OrmTag) TableName() string {
	return TableNameTag
}

func InstallTag(tag OrmTag) error {
	return db.Create(&tag).Error
}
