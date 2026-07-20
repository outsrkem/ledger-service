package statistic

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/shopspring/decimal"
	"ledger/src/models"
	"ledger/src/pkg/answer"
	"ledger/src/pkg/constants"
	"ledger/src/slog"
	"net/http"
	"time"
)

// DayTrend 日趋势
func DayTrend() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		var req QueryArgs
		if err := c.BindAndValidate(&req); err != nil {
			klog.Error("BindAndValidate ", err)
			c.JSON(http.StatusBadRequest, answer.ResBody(answer.EcodeInvalidRequestError, err.Error(), nil))
			return
		}

		from, to, name, dayCnt, err := ResolveCycleRange(req.Ct, req.Td, req.From, req.To)
		if err != nil {
			klog.Error("ResolveCycleRange err", err)
			c.JSON(http.StatusBadRequest, answer.ResBody(answer.EcodeInvalidRequestError, "时间参数错误", nil))
			return
		}
		fmt.Println(from)
		fmt.Println(to)
		fmt.Println(name)
		fmt.Println(dayCnt)
		//2. 获取DB实例
		dao, err := models.NewDBModel(c.GetString(constants.InstanceIdKey))
		if err != nil {
			klog.Errorf("NewDBModel error %s", err)
			c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeError, err.Error(), nil))
			return
		}

		result, err := dao.DayTrend(from, to)
		if err != nil {
			klog.Errorf("统计收支失败: %v", err)
			c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeError, "统计失败", nil))
			return
		}

		dataMap := make(map[string]*models.DayStat)
		for _, item := range result {
			dataMap[item.Day] = item
		}

		// 接口输出专用结构体
		type DayStat struct {
			Day     string `json:"day"`
			Income  string `json:"income"`
			Expense string `json:"expense"`
		}
		var fullList []DayStat

		startTime := time.UnixMilli(from)
		endTime := time.UnixMilli(to)
		current := startTime
		for current.Before(endTime) {
			dayStr := current.Format("2006-01-02")
			inc := decimal.Zero
			exp := decimal.Zero
			if exist, ok := dataMap[dayStr]; ok {
				inc = exist.Income
				exp = exist.Expense
			}
			// 强制保留2位小数,末尾不足自动补 0
			ds := DayStat{
				Day:     dayStr,
				Income:  inc.Round(2).StringFixed(2),
				Expense: exp.Round(2).StringFixed(2),
			}
			fullList = append(fullList, ds)
			current = current.AddDate(0, 0, 1)
		}

		c.JSON(http.StatusOK, answer.ResBody(answer.EcodeOK, nil, fullList))
	}
}

// MonthTrend 按月统计，补齐空白月份、强制两位小数
func MonthTrend() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		var req QueryArgs
		if err := c.BindAndValidate(&req); err != nil {
			klog.Error("BindAndValidate ", err)
			c.JSON(http.StatusBadRequest, answer.ResBody(answer.EcodeInvalidRequestError, err.Error(), nil))
			return
		}

		from, to, name, dayCnt, err := ResolveCycleRange(req.Ct, req.Td, req.From, req.To)
		if err != nil {
			klog.Error("ResolveCycleRange err", err)
			c.JSON(http.StatusBadRequest, answer.ResBody(answer.EcodeInvalidRequestError, "时间参数错误", nil))
			return
		}
		fmt.Println(from)
		fmt.Println(to)
		fmt.Println(name)
		fmt.Println(dayCnt)
		// 获取DB实例
		dao, err := models.NewDBModel(c.GetString(constants.InstanceIdKey))
		if err != nil {
			klog.Errorf("NewDBModel error %s", err)
			c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeError, err.Error(), nil))
			return
		}
		result, err := dao.MonthTrend(from, to)
		if err != nil {
			klog.Errorf("按月统计收支失败: %v", err)
			c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeError, "统计失败", nil))
			return
		}

		// 1. 数据库结果存入map，快速匹配
		dataMap := make(map[string]*models.DayStat)
		for _, item := range result {
			dataMap[item.Day] = item
		}

		// 输出结构体，和日趋势格式完全统一
		type OutputStat struct {
			Day     string `json:"day"`
			Income  string `json:"income"`
			Expense string `json:"expense"`
		}
		var fullList []OutputStat

		startTime := time.UnixMilli(from)
		endTime := time.UnixMilli(to)
		current := startTime

		// 逐月循环生成连续年月 2006-01
		for current.Before(endTime) {
			monthStr := current.Format("2006-01")
			inc := decimal.Zero
			exp := decimal.Zero
			if exist, ok := dataMap[monthStr]; ok {
				inc = exist.Income
				exp = exist.Expense
			}
			// 强制保留两位小数
			item := OutputStat{
				Day:     monthStr,
				Income:  inc.Round(2).StringFixed(2),
				Expense: exp.Round(2).StringFixed(2),
			}
			fullList = append(fullList, item)
			// 月份+1
			current = current.AddDate(0, 1, 0)
		}

		c.JSON(http.StatusOK, answer.ResBody(answer.EcodeOK, nil, fullList))
	}
}

// YearTrend 按年统计，补齐空白年份、强制两位小数
func YearTrend() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		var req QueryArgs
		if err := c.BindAndValidate(&req); err != nil {
			klog.Error("BindAndValidate ", err)
			c.JSON(http.StatusBadRequest, answer.ResBody(answer.EcodeInvalidRequestError, err.Error(), nil))
			return
		}

		from, to, name, dayCnt, err := ResolveCycleRange(req.Ct, req.Td, req.From, req.To)
		if err != nil {
			klog.Error("ResolveCycleRange err", err)
			c.JSON(http.StatusBadRequest, answer.ResBody(answer.EcodeInvalidRequestError, "时间参数错误", nil))
			return
		}
		fmt.Println(from)
		fmt.Println(to)
		fmt.Println(name)
		fmt.Println(dayCnt)
		// 获取DB实例
		dao, err := models.NewDBModel(c.GetString(constants.InstanceIdKey))
		if err != nil {
			klog.Errorf("NewDBModel error %s", err)
			c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeError, err.Error(), nil))
			return
		}
		result, err := dao.YearTrend(from, to)
		if err != nil {
			klog.Errorf("按年统计收支失败: %v", err)
			c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeError, "统计失败", nil))
			return
		}

		// 1. 数据库结果存入map
		dataMap := make(map[string]*models.DayStat)
		for _, item := range result {
			dataMap[item.Day] = item
		}

		// 输出结构体，和日/月格式统一
		type OutputStat struct {
			Day     string `json:"day"`
			Income  string `json:"income"`
			Expense string `json:"expense"`
		}
		var fullList []OutputStat

		startTime := time.UnixMilli(from)
		endTime := time.UnixMilli(to)
		current := startTime

		// 逐年循环生成完整年份 2006
		for current.Before(endTime) {
			yearStr := current.Format("2006")
			inc := decimal.Zero
			exp := decimal.Zero
			if exist, ok := dataMap[yearStr]; ok {
				inc = exist.Income
				exp = exist.Expense
			}
			// 强制两位小数输出
			item := OutputStat{
				Day:     yearStr,
				Income:  inc.Round(2).StringFixed(2),
				Expense: exp.Round(2).StringFixed(2),
			}
			fullList = append(fullList, item)
			// 年份+1
			current = current.AddDate(1, 0, 0)
		}

		c.JSON(http.StatusOK, answer.ResBody(answer.EcodeOK, nil, fullList))
	}
}
