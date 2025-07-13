package route

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"ledger/src/service/category"
	"net/http"
)

func Middleware(h *server.Hertz) {
	h.Use(RequestId())
	h.Use(RequestRecorder())
}

func helloWorld() func(c context.Context, ctx *app.RequestContext) {
	return func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(http.StatusOK, utils.H{"message": "Hello World"})
	}
}

func AppRoute(h *server.Hertz) {
	h.GET("/", helloWorld())
	//income 收入  disbursement 支出
	h.POST("/api/v1/category", category.CreateCategory())       // 创建分类
	h.GET("/api/v1/category", category.SelectCategory())        // 查询分类
	h.PATCH("/api/v1/category/:id", category.UpdateCategory())  // 修改分类
	h.DELETE("/api/v1/category/:id", category.DeleteCategory()) // 删除分类

	h.POST("/api/v1/finances/transactions", helloWorld())       // 添加记账流水
	h.GET("/api/v1/finances/transactions", helloWorld())        // 查询记账流水
	h.PATCH("/api/v1/finances/transactions/:id", helloWorld())  // 修改
	h.DELETE("/api/v1/finances/transactions/:id", helloWorld()) // 删除

	h.GET("/api/v1/finances/statistics", helloWorld()) // 统计收支

}
