package migrator

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 测试用SQLite内存数据库（每次测试全新实例）
func newTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	return db
}

// 测试升级流程
func TestUpMigrationA(t *testing.T) {
	db := newTestDB(t)
	tableName := "auto_migrations"
	require.NoError(t, InitTable(db, tableName))
	rootDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取工作目录失败: %v", err)
	}
	fmt.Println("rootDir:", rootDir)
	// 构建绝对路径
	scriptsDir := filepath.Join(rootDir, "test_scripts")
	// 验证目录是否存在
	if _, err := os.Stat(scriptsDir); os.IsNotExist(err) {
		t.Fatalf("脚本目录不存在: %s", scriptsDir)
	}

	// 使用本地文件提供者
	m := New(db, WithTableName(tableName))
	m.SetScriptProvider(NewLocalFileProvider(scriptsDir)) // 修改这里

	// 初始状态：无脚本执行
	hasPending, err := m.HasPending()
	require.NoError(t, err)
	assert.True(t, hasPending)

	latest, err := m.LatestVersion()
	require.NoError(t, err)
	assert.Empty(t, latest)

	// 执行升级
	require.NoError(t, m.Up())

	// 验证结果
	hasPending, err = m.HasPending()
	require.NoError(t, err)
	assert.False(t, hasPending)

	latest, err = m.LatestVersion()
	require.NoError(t, err)
	assert.Equal(t, "v1.0.1", latest)

	// 验证表是否创建
	var count int64
	db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='test'").Scan(&count)
	assert.Equal(t, int64(1), count) // 补全断言：验证test表存在
}

//go:embed test_scripts/*.sql
var testScripts embed.FS

func TestUpMigrationWithEmbedMigrationsDir(t *testing.T) {
	db := newTestDB(t)
	tableName := "auto_migrations"
	require.NoError(t, InitTable(db, tableName))

	// 验证嵌入资源是否存在
	_, err := os.ReadDir("test_scripts")
	if err != nil {
		t.Fatalf("嵌入脚本目录不存在: %v", err)
	}
	fmt.Println("-----------1----------")

	// 使用嵌入提供者
	m := New(db, WithTableName(tableName), WithEmbedMigrationsDir("test_scripts", testScripts))

	// 初始状态：无脚本执行
	hasPending, err := m.HasPending()
	require.NoError(t, err)
	assert.True(t, hasPending)
	fmt.Println("-----------2----------")

	latest, err := m.LatestVersion()
	require.NoError(t, err)
	assert.Empty(t, latest)

	// 执行升级
	require.NoError(t, m.Up())

	// 验证结果
	hasPending, err = m.HasPending()
	require.NoError(t, err)
	assert.False(t, hasPending)

	latest, err = m.LatestVersion()
	require.NoError(t, err)
	assert.Equal(t, "v1.0.1", latest)

	// 验证表是否创建
	var count int64
	db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='test'").Scan(&count)
	assert.Equal(t, int64(1), count) // 补全断言：验证test表存在
}

func TestNewDefault(t *testing.T) {
	db := newTestDB(t)
	tableName := "auto_migrations"
	m := NewDefault(db, &Config{
		TableName:      tableName,
		ScriptProvider: NewEmbeddedScriptProvider("test_scripts", testScripts),
	})

	// 初始状态：无脚本执行
	hasPending, err := m.HasPending()
	require.NoError(t, err)
	assert.True(t, hasPending)
	fmt.Println("-----------2----------")

	latest, err := m.LatestVersion()
	require.NoError(t, err)
	assert.Empty(t, latest)

	// 执行升级
	require.NoError(t, m.Up())

	// 验证结果
	hasPending, err = m.HasPending()
	require.NoError(t, err)
	assert.False(t, hasPending)

	latest, err = m.LatestVersion()
	require.NoError(t, err)
	assert.Equal(t, "v1.0.1", latest)

	// 验证表是否创建
	var count int64
	db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='test'").Scan(&count)
	assert.Equal(t, int64(1), count) // 验证test表存在
}
