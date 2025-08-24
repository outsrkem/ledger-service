package instance

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"ledger/src/global"
	"ledger/src/models"
	"ledger/src/pkg/common"
)

func GetInstanceId(userId string) (string, error) {
	// 1. 前置校验：用户ID非空且格式合法
	if userId == "" {
		return "", errors.New("user id is empty")
	}

	if !common.IsValid32UUID(userId) {
		return "", errors.New("invalid user ID format, must be a 32-character UUID")
	}
	// 2. 优先查询缓存：命中则直接返回（减少DB访问）
	if instanceId, ok := global.GCache.Get(userId); ok {
		return instanceId, nil
	}

	// 3. 缓存未命中：查询数据库
	instance, err := models.GetInstance(userId)
	switch {
	// 3.1 数据库未找到记录：创建新实例
	case errors.Is(err, gorm.ErrRecordNotFound):
		if err := CreateInstance(userId); err != nil {
			// 创建设置失败：返回具体错误（便于排查）
			return "", fmt.Errorf("failed to create instance for user %s: %w", userId, err)
		}
		// 关键：创建成功后，重新查询数据库获取新实例的InstanceId
		newInstance, err := models.GetInstance(userId)
		if err != nil {
			return "", fmt.Errorf("failed to get new instance after creation: %w", err)
		}
		instanceId := newInstance.InstanceId
		// 更新缓存：将新实例ID存入缓存（有效期1天）
		global.GCache.Set(userId, instanceId, 3600*24)
		return instanceId, nil

	// 3.2 数据库查询错误（如连接失败、SQL错误）：直接返回错误
	case err != nil:
		return "", fmt.Errorf("failed to query instance from DB: %w", err)

	// 3.3 数据库查询成功：更新缓存并返回
	default:
		instanceId := instance.InstanceId
		global.GCache.Set(userId, instanceId, 3600*24)
		return instanceId, nil
	}
}

// CreateInstance 为用户创建新实例（封装创建逻辑，便于复用和维护）
func CreateInstance(userId string) error {
	// 构造实例数据：使用全局工具函数生成UUID和时间戳
	newInstance := models.OrmInstance{
		InstanceId: common.CreateUuid(),      // 生成32位UUID作为实例ID
		UserId:     userId,                   // 关联用户ID
		CreateTime: common.CreateTimestamp(), // 生成毫秒级时间戳（如1724486400000）
	}

	// 调用models层插入数据库
	if err := models.InstallInstance(newInstance); err != nil {
		// 包装错误：添加上下文（用户ID、实例ID），便于日志排查
		return fmt.Errorf("db insert failed (userID: %s, instanceID: %s): %w",
			userId, newInstance.InstanceId, err)
	}
	return nil
}
