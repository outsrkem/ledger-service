package statistic

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

func RecentSixMonthChartTrend() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		//klog := slog.FromContext(c)
		//now := time.Now()
	}
}
