package detail

import (
	"context"
	"ledger/src/models"
	"ledger/src/pkg/answer"
	"ledger/src/pkg/common"
	"ledger/src/service/instance"
	"ledger/src/slog"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/shopspring/decimal"
)

type ReqDetail struct {
	Cid     int64           `json:"cid"`      // 类型ID
	OccTime string          `json:"occ_time"` // 发生的时间ISO 8601, 20220909T18:47:66+0800
	Amount  decimal.Decimal `json:"amount"`   // 金额 0.0000
	Remark  string          `json:"remark"`   // 备注
	Total   int8            `json:"total"`    // 是否计入本月收支，1：计入；0：不计入
	Detail  []*Detail       `json:"detail"`
}

// Detail 账目明细
type Detail struct {
	Name     string          `json:"name"`     // 物品名称（如“牛奶”）
	Quantity int             `json:"quantity"` // 数目（如：2件商品、2次服务）
	Price    decimal.Decimal `json:"price"`    // 单价
	Total    decimal.Decimal `json:"total"`    // 总价（数目 × 单价，与原 amount 绝对值匹配）
	Remark   string          `json:"cremark"`  // 详细备注
}

// 转换时间字符串为毫秒时间戳
func timeToMillisecond(timeStr string) int64 {
	// 参考：https://pkg.go.dev/time#pkg-constants
	layout := "2006-01-02T15:04:05-0700"

	// 解析时间字符串
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return 0
	}

	// 返回毫秒级时间戳（Unix时间戳是秒，乘以1000得到毫秒）
	// 也可以直接使用 t.UnixMilli()（Go 1.17+支持）
	return t.Unix() * 1000
}
func CreateTransaction() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		userId := c.GetString("userId")
		var ReqData ReqDetail
		if err := c.BindJSON(&ReqData); err != nil {
			klog.Errorf("bind json failed %v", err)
			c.JSON(http.StatusBadRequest,
				answer.ResBody(answer.EcodeInvalidRequestError,
					"Failed to parse request data. Check JSON format.", ""))
			return
		}

		// 查询用户的实例ID
		instanceId, err := instance.GetInstanceId(userId)
		if err != nil {
			klog.Errorf("get instance id failed  %v", err)
			c.JSON(http.StatusBadRequest,
				answer.ResBody(answer.EcodeInvalidRequestError,
					"Failed to get your account info. Try again later.", ""))
			return
		}

		now := common.CreateTimestamp()
		dbData := &models.OrmTransaction{
			InstanceId: instanceId,
			CategoryId: ReqData.Cid,
			Amount:     ReqData.Amount,
			OccTime:    ReqData.OccTime,
			Remark:     ReqData.Remark,
			UpdateTime: now,
			CreateTime: now,
		}

		if len(ReqData.Detail) <= 0 {
			// 没有明细
			err = models.InstallTransaction(dbData)
			if err != nil {
				klog.Errorf("insert transaction failed %v", err)
				c.JSON(http.StatusInternalServerError,
					answer.ResBody(answer.EcodeError,
						"Failed to create transaction. Try again later.", ""))
				return
			}
		} else {
			detail := make([]*models.OrmDetail, 0)
			// 有明细详情
			for _, val := range ReqData.Detail {
				detail = append(detail, &models.OrmDetail{
					Name:     val.Name,
					Quantity: val.Quantity,
					Price:    val.Price,
					Total:    val.Total,
				})
			}
			err = models.InstallTransactionAndDetail(dbData, detail)
			if err != nil {
				klog.Errorf("insert transaction failed %v", err)
				c.JSON(http.StatusInternalServerError,
					answer.ResBody(answer.EcodeError,
						"Failed to create transaction. Try again later.", ""))
				return
			}
		}

		c.JSON(http.StatusCreated, answer.ResBody(answer.EcodeOK, nil, nil))
	}
}

func SelectTransaction() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		userId := c.GetString("userId")
		// 查询用户的实例ID
		instanceId, err := instance.GetInstanceId(userId)
		if err != nil {
			klog.Errorf("get instance id failed  %v", err)
			c.JSON(http.StatusBadRequest,
				answer.ResBody(answer.EcodeInvalidRequestError,
					"Failed to get your account info. Try again later.", ""))
			return
		}
		limit, offset, err := common.GetPagingQuery(c)
		if err != nil {
			klog.Error(err)
			return
		}
		count := int64(0)
		result, err := models.FindTransactionAll(instanceId, limit, offset, &count)
		if err != nil {
			klog.Error("in db select role error: ", err)
			c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeError, "Internal server error.", ""))
			return
		}
		type resp struct {
			Kid        int64           `json:"id"`
			CategoryId int64           `json:"category_id"`
			Amount     decimal.Decimal `json:"amount"`
			OccTime    int64           `json:"occ_time"`
			Remark     string          `json:"remark"`
			UpdateTime int64           `json:"update_time"`
			CreateTime int64           `json:"create_time"`
		}
		data := make([]*resp, 0)
		for _, v := range result {
			data = append(data, &resp{
				Kid:        v.Kid,
				CategoryId: v.CategoryId,
				Amount:     v.Amount,
				OccTime:    timeToMillisecond(v.OccTime),
				Remark:     v.Remark,
				UpdateTime: v.UpdateTime,
				CreateTime: v.CreateTime,
			})
		}
		pageInfo := answer.SetPageInfo(limit, offset, count)
		payload := map[string]interface{}{
			"items":     data,
			"page_info": pageInfo,
		}
		klog.Debug("account count: ", count)
		klog.Info("select accounts successfully.")
		c.JSON(http.StatusOK, answer.ResBody(answer.EcodeOK, "", payload))
	}
}
