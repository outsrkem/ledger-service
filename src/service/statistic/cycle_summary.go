package statistic

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/shopspring/decimal"
	"ledger/src/models"
	"ledger/src/pkg/answer"
	"ledger/src/pkg/constants"
	"ledger/src/slog"
	"math"
	"net/http"
	"time"
)

type QueryArgs struct {
	Ct    string `query:"ct"`    // 周期类型 month,year
	Td    string `query:"td"`    // 目标节点  2026-07 / 2026
	Tn    string `query:"tn"`    // 靠前个数topNum
	From  string `query:"from"`  // 自定义起点 YYYY-MM-DD
	To    string `query:"to"`    // 自定义结束 YYYY-MM-DD
	Limit string `query:"limit"` // 近N个周期值
}

// ResolveCycleRange 解析周期时间范围,优先 from/to，没有则根据 ct+td 自动生成
func ResolveCycleRange(ct, td, from, to string) (start, end int64, name string, dayCnt int, err error) {
	loc := time.Local
	var a, b time.Time
	// 自定义周期优先
	if from != "" && to != "" {
		a, err = time.ParseInLocation(time.DateOnly, from, loc)
		if err != nil {
			return
		}
		b, err = time.ParseInLocation(time.DateOnly, to, loc)
		if err != nil {
			return
		}
		// to日期的次日0点（左闭右开标准结束边界）
		b = b.AddDate(0, 0, 1)

		name = from + "~" + to
		dayCnt = int(b.Sub(a).Hours() / 24) // 不用+1，因为b已经是次日零点

		start = a.UnixMilli()
		end = b.UnixMilli() // end = 下一天0点毫秒，查询用 < end
		return
	}

	// 固定月度
	if ct == "month" {
		a, err = time.ParseInLocation("2006-01", td, loc)
		if err != nil {
			return
		}
		// 下个月1号 00:00:00
		b = time.Date(a.Year(), a.Month()+1, 1, 0, 0, 0, 0, loc)

		name = td + "月"
		dayCnt = int(b.Sub(a).Hours() / 24)
		start = a.UnixMilli()
		end = b.UnixMilli()
		return
	}

	// 固定年度
	if ct == "year" {
		a, err = time.ParseInLocation("2006", td, loc)
		if err != nil {
			return
		}
		// 次年1月1日 00:00:00
		b = time.Date(a.Year()+1, 1, 1, 0, 0, 0, 0, loc)

		name = td + "年"
		dayCnt = int(b.Sub(a).Hours() / 24)
		start = a.UnixMilli()
		end = b.UnixMilli()
		return
	}

	err = errors.New("cycle type error")
	return
}
func Round2(f float64) float64 {
	return math.Round(f*100) / 100
}

// CycleSummary 月度汇总
// 月度 ct=month&td=2026-07
// 年度 ct=year&td=2026
// from=2026-05-01&to=2026-07-15
func CycleSummary() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		var req QueryArgs
		if err := c.BindAndValidate(&req); err != nil {
			klog.Error("BindAndValidate ", err)
			c.JSON(http.StatusBadRequest, answer.ResBody(answer.EcodeInvalidRequestError, err.Error(), nil))
			return
		}
		//1. 解析时间区间（优先自定义from/to，否则根据 ct+td 自动计算）
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
		result, err := dao.CycleSummary(from, to)
		if err != nil {
			klog.Errorf("统计收支失败: %v", err)
			c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeError, "统计失败", nil))
			return
		}
		type S01StatResp struct {
			Income  decimal.Decimal `json:"income"`  // 总收入
			Expense decimal.Decimal `json:"expense"` // 总支出
		}
		payload := S01StatResp{
			Income:  result.Income,
			Expense: result.Expense,
		}
		c.JSON(http.StatusOK, answer.ResBody(answer.EcodeOK, nil, payload))
	}
}
