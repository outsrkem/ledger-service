package route

import (
	"context"
	"ledger/src/service/category"
	"ledger/src/service/tag"
	"ledger/src/service/transaction"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func Middleware(h *server.Hertz) {
	h.Use(RequestId())
	h.Use(RequestRecorder())
}

func helloWorld() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, utils.H{"message": "Hello World."})
	}
}

func AppRoute(h *server.Hertz) {
	h.HEAD("")
	h.GET("/", helloWorld())

	h.POST("/v1/category/main", apc("ledger:category:create"), category.CreateCategoryMain())           // 添加主分类 √
	h.POST("/v1/category/:categoryId/sub", apc("ledger:category:create"), category.CreateCategorySub()) // 添加子分类 √
	h.DELETE("/v1/category/:kid", apc("ledger:category:delete"), category.DeleteCategory())             // 删除分类
	h.GET("/v1/category", apc("ledger:category:list"), category.SelectCategory())                       // 查询分类 √
	h.PATCH("/v1/category/:id", apc("ledger:category:update"), category.UpdateCategory())               // 修改分类

	h.POST("/v1/transactions", apc("ledger:transaction:create"), transaction.CreateTransaction())       // 添加记账流水 √
	h.GET("/v1/transactions", apc("ledger:transaction:list"), transaction.SelectTransaction())          // 查询记账流水 √
	h.PATCH("/v1/transactions/:id", apc("ledger:transaction:update"), helloWorld())                     // 修改
	h.DELETE("/v1/transactions/:id", apc("ledger:transaction:delete"), transaction.DeleteTransaction()) // 删除 √

	h.POST("/v1/tag", apc("ledger:tag:create"), tag.CreateTag()) // 创建标签 √

	h.GET("/v1/finances/statistics", apc(""), helloWorld()) // 统计收支
}
