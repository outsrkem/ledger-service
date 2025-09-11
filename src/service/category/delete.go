package category

import (
	"context"
	"errors"
	"ledger/src/models"
	"ledger/src/pkg/answer"
	"ledger/src/service/instance"
	"ledger/src/slog"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
)

func DeleteCategory() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		klog.Info("CreateCategory")
		userId := c.GetString("userId")
		categoryId := c.Param("kid")
		// 约束：有子分类的时候不能删除主类
		instanceId, err := instance.GetInstanceId(userId)
		if err != nil {
			klog.Errorf("get instance id failed %v", err)
			c.JSON(http.StatusBadRequest,
				answer.ResBody(answer.EcodeInvalidRequestError,
					"Failed to get your account info. Try again later.", ""))
			return
		}

		// TODO 系统策略提示无权限删除
		categ, err := models.FindCategoryById(instanceId, strToInt64(categoryId))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				klog.Errorf("category not found %v", err)
				c.JSON(http.StatusBadRequest,
					answer.ResBody(answer.EcodeError,
						"category not found.", nil))
				return
			}

			klog.Errorf("find category by id failed %v", err)
			c.JSON(http.StatusBadRequest,
				answer.ResBody(answer.EcodeError,
					"Failed to get your account info. Try again later.", nil))
			return
		}

		// 系统分类不能删除
		if categ.InstanceId == "" {
			klog.Error("system category cannot be deleted")
			c.JSON(http.StatusForbidden,
				answer.ResBody(answer.EcodeError,
					"system category cannot be deleted.", nil))
			return
		}

		_, err = models.DeleteCategory(instanceId, []int64{categ.Kid})
		if err != nil {
			klog.Errorf("delete category failed %v", err)
			c.JSON(http.StatusBadRequest,
				answer.ResBody(answer.EcodeError,
					"delete category failed.", nil))
			return
		}

		c.JSON(http.StatusOK, answer.ResBody(answer.EcodeOK, nil, nil))
	}
}
