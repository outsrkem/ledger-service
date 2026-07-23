package transaction

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

// Detail 账目明细
type Detail struct {
	Name     string          `json:"name"`     // 物品名称（如“牛奶”）
	Quantity decimal.Decimal `json:"quantity"` // 数量（如：2件商品、2次服务）
	Unit     string          `json:"unit"`     // 单位
	Price    decimal.Decimal `json:"price"`    // 单价
	Total    decimal.Decimal `json:"total"`    // 总价（数目 × 单价，与原 amount 绝对值匹配）
}

// ReqTransaction 请求数据
type ReqTransaction struct {
	Cid     int64           `json:"cid"`      // 类型ID
	OccTime string          `json:"occ_time"` // 发生的时间ISO 8601, 20220909T18:47:66+0800
	Amount  decimal.Decimal `json:"amount"`   // 金额 0.0000
	Remark  string          `json:"remark"`   // 备注
	Total   int8            `json:"total"`    // 是否计入本月收支，1：计入；0：不计入
	Detail  []*Detail       `json:"detail"`
}

// 转换时间字符串为毫秒时间戳
func timeToMillisecond(timeStr string) int64 {
	// https://pkg.go.dev/time#pkg-constants
	layout := "2006-01-02T15:04:05-0700"
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return 0
	}

	return t.UnixMilli()
}

func CreateTransaction() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		userId := c.GetString("userId")
		var ReqData ReqTransaction
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
			OccAt:      timeToMillisecond(ReqData.OccTime),
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
					Unit:     val.Unit,
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
		type Category struct {
			Name string `json:"name,omitempty"`
			Id   int64  `json:"id,omitempty"`
		}
		type resp struct {
			Kid        int64           `json:"id"`
			CategoryId int64           `json:"category_id"`
			Category   Category        `json:"category"`
			Amount     decimal.Decimal `json:"amount"`
			OccTime    int64           `json:"occ_time"`
			Remark     string          `json:"remark"`
			UpdateTime int64           `json:"update_time"`
			CreateTime int64           `json:"create_time"`
		}
		data := make([]*resp, len(result))
		for k, v := range result {
			item := &resp{
				Kid:        v.Kid,
				CategoryId: v.CategoryId,
				Category: Category{
					Id: v.CategoryId,
				},
				Amount:     v.Amount,
				OccTime:    timeToMillisecond(v.OccTime),
				Remark:     v.Remark,
				UpdateTime: v.UpdateTime,
				CreateTime: v.CreateTime,
			}
			data[k] = item
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

type BillDetailResp struct {
	Cid     int64           `json:"cid"`      // 对应 CategoryId
	OccTime string          `json:"occ_time"` // 时间字符串原样返回
	Amount  decimal.Decimal `json:"amount"`   // 主账单金额
	Remark  string          `json:"remark"`   // 备注
	Detail  []DetailItem    `json:"detail"`   // 明细数组
}

// DetailItem 账单明细项
type DetailItem struct {
	Name     string          `json:"name"`     // 物品名称
	Quantity decimal.Decimal `json:"quantity"` // 数目
	Price    decimal.Decimal `json:"price"`    // 单价
	Total    decimal.Decimal `json:"total"`    // 总价 = 数量 × 单价
	Unit     string          `json:"unit"`     // 单位
}

// BillDetails 账单详情
func BillDetails() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		klog.Debug("BillDetails")

		// 1. 获取用户实例ID
		userId := c.GetString("userId")
		instanceId, err := instance.GetInstanceId(userId)
		if err != nil {
			klog.Errorf("get instance id failed: %v", err)
			c.JSON(http.StatusBadRequest,
				answer.ResBody(answer.EcodeInvalidRequestError,
					"Failed to get your account info. Try again later.", nil))
			return
		}

		// 2. 绑定路径参数 billId
		type Query struct {
			BillId int64 `path:"id"` // 统一大写，避免绑定失败
		}
		var q Query
		if err := c.BindPath(&q); err != nil {
			klog.Errorf("bind path failed: %v", err)
			c.JSON(http.StatusBadRequest,
				answer.ResBody(answer.EcodeInvalidRequestError,
					"invalid query parameter", nil))
			return
		}

		// 3. 查询主账单
		billList, err := models.GetBillById(instanceId, q.BillId)
		if err != nil {
			klog.Errorf("get bill failed: %v", err)
			c.JSON(http.StatusInternalServerError,
				answer.ResBody(answer.EcodeError,
					"get bill info error", nil))
			return
		}
		if len(billList) == 0 {
			c.JSON(http.StatusOK,
				answer.ResBody(answer.EcodeOK, "bill not found", nil))
			return
		}
		bill := billList[0] // 单条账单

		// 4. 查询账单明细
		detailList, err := models.GetDetail(q.BillId)
		if err != nil {
			klog.Errorf("get bill detail failed: %v", err)
			c.JSON(http.StatusInternalServerError,
				answer.ResBody(answer.EcodeError,
					"get bill detail error", nil))
			return
		}

		// 5. 组装明细结构
		detailItems := make([]DetailItem, 0, len(detailList))
		for _, d := range detailList {
			detailItems = append(detailItems, DetailItem{
				Name:     d.Name,
				Quantity: d.Quantity,
				Price:    d.Price,
				Total:    d.Total,
				Unit:     d.Unit,
			})
		}

		// 6. 组装最终返回体（完全匹配你要求的JSON）
		respData := BillDetailResp{
			Cid:     bill.CategoryId,
			OccTime: bill.OccTime,
			Amount:  bill.Amount,
			Remark:  bill.Remark,
			Detail:  detailItems,
		}

		// 7. 返回成功响应
		c.JSON(http.StatusOK, answer.ResBody(answer.EcodeOK, "", respData))
	}
}
