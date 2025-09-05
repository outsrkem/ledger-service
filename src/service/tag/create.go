package tag

import (
	"context"
	"ledger/src/models"
	"ledger/src/pkg/answer"
	"ledger/src/pkg/common"
	"ledger/src/service/instance"
	"ledger/src/slog"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
)

type ReqTag struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// CreateTag 添加新标签
func CreateTag() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)

		userId := c.GetString("userId")
		instanceId, err := instance.GetInstanceId(userId)
		if err != nil {
			klog.Errorf("get instance id failed  %v", err)
			c.JSON(http.StatusBadRequest,
				answer.ResBody(answer.EcodeInvalidRequestError,
					"Failed to get your account info. Try again later.", ""))
			return
		}

		var ReqData ReqTag
		if err := c.BindJSON(&ReqData); err != nil {
			klog.Errorf("bind json failed %v", err)
			c.JSON(http.StatusBadRequest,
				answer.ResBody(answer.EcodeInvalidRequestError,
					"Failed to parse request data. Check JSON format.", ""))
			return
		}

		now := common.CreateTimestamp()
		data := models.OrmTag{
			InstanceId: instanceId,
			Name:       ReqData.Name,
			Color:      ReqData.Color,
			CreateTime: now,
			UpdateTime: now,
		}

		if err := models.InstallTag(data); err != nil {
			klog.Errorf("insert tag failed %v", err)
			c.JSON(http.StatusInternalServerError,
				answer.ResBody(answer.EcodeError, "Failed to create tag.", ""))
			return
		}

		c.JSON(http.StatusCreated, answer.ResBody(answer.EcodeOK, nil, nil))
	}
}
