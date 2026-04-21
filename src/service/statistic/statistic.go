package statistic

import (
	"context"
	"fmt"
	"ledger/src/models"
	"ledger/src/pkg/answer"
	"ledger/src/pkg/constants"
	"ledger/src/slog"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/shopspring/decimal"
)

// StatOve 数据总览
func StatOve() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		sa := c.Query("sa") // 总览类型
		pt := c.Query("pt") // 统计周期

		// TODO 处理查询周期
		from := c.Query("from")
		to := c.Query("to")

		switch sa {
		case "s01": // 统计收支
			S01(c, pt, from, to)
		case "s02": // 分类统计
			// ...
		default:
			// ...
		}
	}
}

func S01(c *app.RequestContext, pt, from, to string) {
	klog := slog.FromContext(c)

	// 参数默认值（不传则默认：今天）
	if from == "" || to == "" {
		now := time.Now().Format("2006-01-02")
		from = now
		to = now
	}

	// ===================== 时间格式处理 =====================
	var (
		startDate, endDate string
		parseErr           error
	)

	switch pt {
	case "year":
		// 年份：YYYY → 转换为 2025-01-01 ~ 2025-12-31
		var year int
		_, parseErr = fmt.Sscanf(from, "%d", &year)
		if parseErr != nil {
			klog.Errorf("年份格式错误: %v", parseErr)
			c.JSON(http.StatusBadRequest, answer.ResBody(answer.EcodeError, "年份格式必须为 YYYY", nil))
			return
		}
		startDate = fmt.Sprintf("%d-01-01", year)
		endDate = fmt.Sprintf("%d-12-31", year) // 修正为12-31更合理

	case "month":
		t, err := time.Parse("2006-01", from)
		if err != nil {
			return
		}

		sd := t.Format("2006-01-02")
		ed := t.AddDate(0, 1, 0).Add(-time.Nanosecond).Format("2006-01-02")

		startDate, endDate = sd, ed

	case "custom":
		// 自定义：校验 YYYY-MM-DD 合法性
		startDate = from
		endDate = to
		// 校验开始日期
		_, err1 := time.Parse("2006-01-02", startDate)
		// 校验结束日期
		_, err2 := time.Parse("2006-01-02", endDate)
		if err1 != nil || err2 != nil {
			klog.Errorf("自定义日期格式错误 from=%s to=%s", from, to)
			c.JSON(http.StatusBadRequest, answer.ResBody(answer.EcodeError, "日期格式必须为 YYYY-MM-DD", nil))
			return
		}

	default:
		c.JSON(http.StatusBadRequest, answer.ResBody(answer.EcodeError, "不支持的时间类型", nil))
		return
	}
	// ==============================================================

	// 初始化数据库连接
	dao, err := models.NewDBModel(c.GetString(constants.InstanceIdKey))
	if err != nil {
		klog.Errorf("NewDBModel error %s", err)
		c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeError, err.Error(), nil))
		return
	}

	// 调用统计（传入处理后的标准日期）
	result, err := dao.CountSouZhi(pt, startDate, endDate)
	if err != nil {
		klog.Errorf("统计收支失败: %v", err)
		c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeError, "统计失败", nil))
		return
	}

	type S01StatResp struct {
		Income  decimal.Decimal `json:"income"`  // 总收入
		Expense decimal.Decimal `json:"expense"` // 总支出
	}
	// 构造返回
	payload := S01StatResp{
		Income:  result.Income,
		Expense: result.Expense,
	}

	// 返回成功响应
	c.JSON(http.StatusOK, answer.ResBody(answer.EcodeOK, nil, payload))
}
