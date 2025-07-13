package category

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"ledger/src/slog"
)

func DeleteCategory() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		klog.Info("CreateCategory")
	}
}
