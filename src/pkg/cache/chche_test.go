package cache

import (
	"fmt"
	"testing"
)

var (
	// InstanceCache 全局实例ID缓存（项目启动时初始化）
	instanceCache = NewInstanceCache()
)

// GetUserInstance 获取用户的实例ID（优先读缓存，缓存未命中则查库并更新缓存）
func GetUserInstance(userID string) (string, error) {
	fmt.Println("GetUserInstance")
	// 1. 先查缓存
	if instanceID, ok := instanceCache.Get(userID); ok {
		fmt.Println("缓存命中")
		return instanceID, nil
	}
	fmt.Println("先查缓存")
	// 2. 缓存未命中，查数据库（这里替换为你的实际查库逻辑）
	instanceID := userID
	fmt.Println("缓存未命中")
	// 3. 查库成功，更新缓存（设置1小时过期，可根据业务调整）
	instanceCache.Set(userID, instanceID, 3600)
	fmt.Println(" 查库成功，更新缓存")
	return instanceID, nil
}

func TestCache(t *testing.T) {
	fmt.Println(GetUserInstance("asdasdsad111111111111111"))
	fmt.Println("----------------------")
	fmt.Println(GetUserInstance("asdasdsad111111111111111"))

	fmt.Println(GetUserInstance("yyyyyyyyyyyyyyyyyyyyyyyyyyyy"))
	fmt.Println("----------------------")
	fmt.Println(GetUserInstance("yyyyyyyyyyyyyyyyyyyyyyyyyyyy"))
}
