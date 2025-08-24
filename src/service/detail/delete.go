package detail

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"ledger/src/models"
	"ledger/src/pkg/answer"
	"ledger/src/service/instance"
	"ledger/src/slog"
	"net/http"
)

func DeleteTransaction() func(ctx context.Context, c *app.RequestContext) {
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

		type Query struct {
			TransactionId int64 `path:"id"`
		}

		var q Query
		if err := c.BindPath(&q); err != nil {
			klog.Errorf("bind json failed %v", err)
			c.JSON(http.StatusBadRequest,
				answer.ResBody(answer.EcodeInvalidRequestError,
					"query error.", nil))
			return
		}
		// 查询
		res, err := models.FindTransactionForUser(instanceId, q.TransactionId)
		if err != nil {
			klog.Errorf("delete transaction failed  %v", err)
			c.JSON(http.StatusBadRequest,
				answer.ResBody(answer.EcodeInvalidRequestError,
					"delete transaction failed.", nil))
			return
		}
		// 没查到
		if len(res) == 0 {
			klog.Errorf("transaction not found %v", res)
			c.JSON(http.StatusNotFound,
				answer.ResBody(answer.EcodeInvalidRequestError,
					"transaction not found.", nil))
			return
		}
		// 删除
		if err := models.DeleteTransaction(instanceId, q.TransactionId); err != nil {
			klog.Errorf("delete transaction failed  %v", err)
			c.JSON(http.StatusBadRequest,
				answer.ResBody(answer.EcodeInvalidRequestError,
					"delete transaction failed.", nil))
			return
		}
		// 删除成功
		c.JSON(http.StatusOK, answer.ResBody(answer.EcodeOK, nil, nil))
	}

}
