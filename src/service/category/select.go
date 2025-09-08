package category

import (
	"context"
	"ledger/src/models"
	"ledger/src/pkg/answer"
	"ledger/src/service/instance"
	"ledger/src/slog"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
)

type RespPayload struct {
	Id       int64          `json:"id"`                 // 主键ID
	Name     string         `json:"name"`               // 分类名称
	Children []*RespPayload `json:"children,omitempty"` // 子类
}

// ConvertToTree 将OrmCategory列表转换为树形Payload结构
func ConvertToTree(categories []*models.CategoryWithRelResult) []*RespPayload {
	// 创建ID到Payload指针的映射
	nodeMap := make(map[int64]*RespPayload)
	var roots []*RespPayload

	// 先将所有节点转换为Payload指针并存入映射
	for _, cat := range categories {
		node := &RespPayload{
			Id:       cat.Kid,
			Name:     cat.Name,
			Children: []*RespPayload{},
		}
		nodeMap[cat.Kid] = node
	}

	// 建立父子关系
	for _, cat := range categories {
		currentNode := nodeMap[cat.Kid]

		if cat.Pid == 0 {
			roots = append(roots, currentNode)
		} else {
			// 找到父节点并添加到其子节点列表
			if parentNode, exists := nodeMap[cat.Pid]; exists {
				parentNode.Children = append(parentNode.Children, currentNode)
			}
		}
	}
	return roots
}

// SelectCategory 查询分类
func SelectCategory() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)

		type Query struct {
			Direction int8 `query:"direction"`
		}
		var query Query
		err := c.BindAndValidate(&query)
		if err != nil {
			klog.Errorf("param bind failed: %v", err)
			c.JSON(http.StatusBadRequest,
				answer.ResBody(answer.EcodeInvalidRequestError,
					"Failed to validate request parameters. Please check your input format.", ""))
			return
		}

		if query.Direction != 1 && query.Direction != 2 {
			klog.Warnf("invalid direction: %d (valid: 1,2)", query.Direction)
			c.JSON(http.StatusBadRequest,
				answer.ResBody(answer.EcodeInvalidRequestError,
					"Invalid direction parameter. Allowed values are 1 (income) or 2 (expense).", ""))
			return
		}

		instanceId, err := instance.GetInstanceId(c.GetString("userId"))
		if err != nil {
			klog.Errorf("get instance id failed: %v", err)
			c.JSON(http.StatusInternalServerError,
				answer.ResBody(answer.EcodeError,
					"Failed to retrieve your account information. Please try again later.", ""))
			return
		}

		klog.Debugf("got instance id: %s", instanceId)

		category, err := models.GetCategoryWithRel(instanceId, query.Direction)
		if err != nil {
			klog.Errorf("query category failed (dir: %d): %v", query.Direction, err)
			c.JSON(http.StatusInternalServerError,
				answer.ResBody(answer.EcodeError,
					"Failed to load categories. Please try again later.", ""))
			return
		}

		payload := map[string]interface{}{
			"items": ConvertToTree(category),
		}

		klog.Debugf("category query success (count: %d)", len(category))
		c.JSON(http.StatusOK, answer.ResBody(answer.EcodeOK, nil, payload))
	}
}
