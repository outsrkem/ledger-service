package route

import (
	"context"
	"ledger/src/service/category"
	"ledger/src/service/statistic"
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

	h.GET("/v1/bill/:id", apc("ledger:transaction:list"), transaction.BillDetails()) // 账单详情ledger:bill:get

	h.POST("/v1/tag", apc("ledger:tag:create"), tag.CreateTag()) // 创建标签 √

	h.GET("/v1/bill/statistic", apc("ledger:transaction:list"), statistic.StatOve()) // 统计收支,年，月，自定义

	//	/v1/stat/cycle-summary?cycleType=month&targetDate=2026-07 月度汇总 ct=month&td=2026-07
	//	/v1/stat/chart/day?startDate=2026-07-01&endDate=2026-07-31 当月每日趋势 from=2026-07-01&to=2026-07-31
	//	/v1/stat/category?cycleType=month&targetDate=2026-07 分类饼图 ct=month&td=2026-07
	//	/v1/stat/chart/month?limit=6 近6月对比柱状
	//	/v1/stat/expense-rank?cycleType=month&targetDate=2026-07&topNum=10 支出分类排行 ct=month&td=2026-07&tn=10

	h.GET("/v1/stat/cycle-summary", apc("ledger:transaction:list"), statistic.CycleSummary()) // √ 月度汇总
	h.GET("/v1/stat/chart/day", apc("ledger:transaction:list"), statistic.DayTrend())         // √ 日收支趋势
	h.GET("/v1/stat/chart/month", apc("ledger:transaction:list"), statistic.MonthTrend())     // √ 月收支趋势
	h.GET("/v1/stat/chart/year", apc("ledger:transaction:list"), statistic.YearTrend())       // √ 年收支趋势
	//h.GET("/v1/stat/category", apc("ledger:transaction:list")) // 收支分类构成
	//h.GET("/v1/stat/chart/month", apc("ledger:transaction:list")) // 近 6 个月收支
	//h.GET("/v1/stat/expense-rank", apc("ledger:transaction:list")) // 本月支出分类 TOP 排行
}
