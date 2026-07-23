package transaction

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/shopspring/decimal"
	"ledger/src/models"
	"ledger/src/pkg/answer"
	"ledger/src/service/instance"
	"ledger/src/slog"
	"net/http"
	"strconv"
)

type ReqDetailItem struct {
	Name     *string          `json:"name,omitempty"`     // 物品名称
	Quantity *decimal.Decimal `json:"quantity,omitempty"` // 数量
	Unit     *string          `json:"unit,omitempty"`     // 单位
	Price    *decimal.Decimal `json:"price,omitempty"`    // 单价
	Total    *decimal.Decimal `json:"total,omitempty"`    // 总价
}

type ReqUpdateTransaction struct {
	Cid     *int64           `json:"cid,omitempty"`      // 类型ID
	OccTime *string          `json:"occ_time,omitempty"` // 发生的时间ISO 8601
	Amount  *decimal.Decimal `json:"amount,omitempty"`   // 金额 0.0000
	Remark  *string          `json:"remark,omitempty"`   // 备注
	Total   *int8            `json:"total,omitempty"`    // 是否计入本月收支，1：计入；0：不计入
	Detail  *[]ReqDetailItem `json:"detail,omitempty"`
}

// BuildTransUpdateMap 只提取【主交易表字段】，自动剔除detail
func BuildTransUpdateMap(req *ReqUpdateTransaction) (map[string]any, error) {
	m := make(map[string]any)
	if req.Cid != nil {
		m["category_id"] = *req.Cid
	}
	if req.OccTime != nil {
		m["occ_time"] = *req.OccTime
	}
	if req.Amount != nil {
		m["amount"] = *req.Amount
	}
	if req.Remark != nil {
		m["remark"] = *req.Remark
	}
	if req.Total != nil {
		m["total"] = *req.Total
	}
	return m, nil
}

// ConvertToOrmDetail 将请求明细转为 OrmDetail 数组
func ConvertToOrmDetail(tid int64, items []ReqDetailItem) ([]*models.OrmDetail, error) {
	var list []*models.OrmDetail
	for _, v := range items {
		row := models.OrmDetail{
			TransactionId: tid,
		}
		if v.Name != nil {
			row.Name = *v.Name
		}
		if v.Quantity != nil {
			row.Quantity = *v.Quantity
		}
		if v.Unit != nil {
			row.Unit = *v.Unit
		}
		if v.Price != nil {
			row.Price = *v.Price
		}
		if v.Total != nil {
			row.Total = *v.Total
		}
		list = append(list, &row)
	}

	// 关键：items 是空切片时，强制返回非nil空切片
	if list == nil {
		list = make([]*models.OrmDetail, 0)
	}
	return list, nil
}
func UpdateTransaction() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)

		var req ReqUpdateTransaction
		if err := c.BindJSON(&req); err != nil {
			klog.Errorf("bind json failed %v", err)
			c.JSON(http.StatusBadRequest,
				answer.ResBody(answer.EcodeInvalidRequestError,
					"Failed to parse request data. Check JSON format.", nil))
			return
		}

		// 1. 获取路由参数 tid /transaction/:id
		tidStr := c.Param("id")
		tid, err := strconv.ParseInt(tidStr, 10, 64)
		if err != nil {
			klog.Error(err)
			c.JSON(http.StatusBadRequest,
				answer.ResBody(answer.EcodeInvalidRequestError, "invalid transaction id", nil))
			return
		}

		userId := c.GetString("userId")
		instanceId, err := instance.GetInstanceId(userId)
		if err != nil {
			klog.Errorf("get instance id failed %v", err)
			c.JSON(http.StatusBadRequest,
				answer.ResBody(answer.EcodeInvalidRequestError,
					"Failed to get your account info. Try again later.", nil))
			return
		}
		fmt.Println(req)
		// 2. 构建主表更新map
		transMap, err := BuildTransUpdateMap(&req)
		if err != nil {
			klog.Error("build trans update map error", "err", err)
			c.JSON(http.StatusInternalServerError,
				answer.ResBody(answer.EcodeError, "build update data failed", nil))
			return
		}

		// 3. 处理明细
		var detailList []*models.OrmDetail
		if req.Detail != nil {
			detailList, err = ConvertToOrmDetail(tid, *req.Detail)
			if err != nil {
				klog.Error("convert detail error", "err", err)
				c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeError, "parse detail failed", nil))
				return
			}
		}

		// 校验：主表无更新 且 detailList == nil（前端没有传detail字段）→ 无任何更新内容
		if len(transMap) == 0 && detailList == nil {
			c.JSON(http.StatusBadRequest,
				answer.ResBody(answer.EcodeInvalidRequestError, "no fields to update", nil))
			return
		}

		// 4. 初始化DAO
		dao, err := models.NewDBModel(instanceId)
		if err != nil {
			klog.Errorf("NewDBModel error %s", err)
			c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeError, err.Error(), nil))
			return
		}

		// 5. 调用DAO，只传3个参数，匹配签名
		err = dao.UpdateTransaction(tid, transMap, detailList)
		if err != nil {
			klog.Error("dao UpdateTransaction failed", "err", err)
			c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeError, "update failed", nil))
			return
		}

		c.JSON(http.StatusOK, answer.ResBody(answer.EcodeOK, "success", nil))
	}
}
