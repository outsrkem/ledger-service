package category

import (
	"context"
	"ledger/src/slog"

	"github.com/cloudwego/hertz/pkg/app"
)

func DeleteCategory() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		klog.Info("CreateCategory")
		// 约束：有子分类的时候不能删除主类
	}
}
