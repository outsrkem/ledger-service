package instance

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"ledger/src/global"
	"ledger/src/models"
	"ledger/src/pkg/common"

	"gorm.io/gorm"
)

const (
	InstanceValidity = 3600 * 2               // 实例ID缓存有效期(秒)
	retryMax         = 3                      // 最大重试次数
	retryDelay       = 100 * time.Millisecond // 重试延迟
)

// 用户级别的锁映射，避免并发创建
var userLocks = struct{ sync.Map }{}

// GetInstanceId 获取用户的实例ID，确保每个用户只有一个实例
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
	// 3.1 数据库未找到记录：尝试创建新实例
	case errors.Is(err, gorm.ErrRecordNotFound):
		// 使用应用层锁机制确保并发安全
		instanceId, err := createInstanceSafely(userId)
		if err != nil {
			return "", fmt.Errorf("failed to create instance for user %s: %w", userId, err)
		}
		return instanceId, nil

	// 3.2 数据库查询错误（如连接失败、SQL错误）：直接返回错误
	case err != nil:
		return "", fmt.Errorf("failed to query instance from DB: %w", err)

	// 3.3 数据库查询成功：更新缓存并返回
	default:
		instanceId := instance.InstanceId
		global.GCache.Set(userId, instanceId, InstanceValidity)
		return instanceId, nil
	}
}

// createInstanceSafely 安全地为用户创建实例，避免重复创建
func createInstanceSafely(userId string) (string, error) {
	// 获取用户专属锁，确保同一用户并发请求时只有一个能执行创建逻辑
	lockKey := fmt.Sprintf("instance_lock_%s", userId)

	// 尝试获取锁，已存在则等待
	var lock *sync.Mutex
	for {
		if value, loaded := userLocks.LoadOrStore(lockKey, &sync.Mutex{}); loaded {
			// 锁已存在，获取已有的锁
			lock = value.(*sync.Mutex)
			break
		}
		// 刚创建的锁可能被其他goroutine释放，短暂等待后重试
		time.Sleep(10 * time.Millisecond)
	}

	// 加锁确保原子操作
	lock.Lock()
	defer func() {
		lock.Unlock()
		// 释放锁资源
		userLocks.Delete(lockKey)
	}()

	// 双重检查：再次确认实例是否已存在
	existingInstance, err := models.GetInstance(userId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", fmt.Errorf("failed to check instance existence: %w", err)
	}

	// 如果实例已存在（无错误），直接返回
	if err == nil {
		global.GCache.Set(userId, existingInstance.InstanceId, InstanceValidity)
		return existingInstance.InstanceId, nil
	}

	// 准备创建新实例
	newInstance := models.OrmInstance{
		InstanceId: common.CreateUuid(),
		UserId:     userId,
		CreateTime: common.CreateTimestamp(),
	}

	// 尝试创建实例，处理可能的唯一约束冲突
	var lastErr error
	for i := 0; i < retryMax; i++ {
		// 最后一次尝试不使用重试
		if i == retryMax-1 {
			if err := models.InstallInstance(newInstance); err != nil {
				lastErr = fmt.Errorf("db insert failed (userID: %s, instanceID: %s): %w",
					userId, newInstance.InstanceId, err)
				break
			}
			// 创建成功，更新缓存并返回
			global.GCache.Set(userId, newInstance.InstanceId, InstanceValidity)
			return newInstance.InstanceId, nil
		}

		// 带重试的尝试
		if err := models.InstallInstance(newInstance); err != nil {
			lastErr = err
			// 短暂延迟后重试
			time.Sleep(retryDelay)
			continue
		}

		// 创建成功
		global.GCache.Set(userId, newInstance.InstanceId, InstanceValidity)
		return newInstance.InstanceId, nil
	}

	return "", lastErr
}
