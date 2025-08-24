package route

import (
	"context"
	"encoding/json"
	"fmt"
	"ledger/src/config"
	"ledger/src/pkg/uuid"
	"ledger/src/slog"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hiuias/uias-sdk-go"
)

func RequestId() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		xRequestId := string(c.GetHeader("X-Request-Id"))
		if xRequestId == "" {
			xRequestId = uuid.V4UUID()
			c.Response.Header.Set("X-Request-Id", xRequestId)
			klog.Warnf("request id is empty, Set a new request id: %s", xRequestId)
		}
		c.Set("xRequestId", xRequestId)
		c.Next(ctx)
	}
}

func RequestRecorder() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		start := time.Now()
		c.Next(ctx)
		stop := time.Now()
		latency := stop.Sub(start)
		klog.Infof("|%14s | %d |%7s %s",
			latency, c.Response.StatusCode(), string(c.Request.Method()), c.Request.URI().String())
	}
}

const xRequestIdKey = "X-Request-Id"
const xAuthTokenKey = "X-Auth-Token"

func apc(action string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		client, err := uias.NewClientBuilder().
			WithEndpoint(config.AppCfg.Ledger.Uias.Endpoint).
			WithTimeout(10 * time.Second).
			WithSkipTlsVerify(false).
			Build()
		if err != nil {
			fmt.Printf("Failed to create client: %v\n", err)
			return
		}
		type actionRaw struct {
			Uias struct {
				Action string `json:"action"`
			} `json:"uias"`
		}

		var raw actionRaw
		raw.Uias.Action = action
		rawBody, err := json.Marshal(raw)
		if err != nil {
			klog.Errorf("Error marshaling audit log: %v", err)
			c.Abort()
			return
		}

		resp, err := client.VerifyAction(
			ctx,
			c.Request.Header.Get(xRequestIdKey),
			c.Request.Header.Get(xAuthTokenKey),
			rawBody,
		)
		if err != nil {
			klog.Warnf("Permission verification error, upstream exception: %v", err)
			c.JSON(500, resp)
			c.Abort()
			return
		}

		if resp.Payload.Authentication != 1 {
			// 没有权限，返回403和上游返回体，便于查看问题
			klog.Warnf("Permission denial. result: %+v", resp)
			c.JSON(403, resp)
			c.Abort()
			return
		}

		klog.Info("Permission is granted, and the operation is authorized.")
		c.Set("domainId", resp.Payload.User.Domain.Id)
		c.Set("domainName", resp.Payload.User.Domain.Name)
		c.Set("userId", resp.Payload.User.Id)
		c.Set("account", resp.Payload.User.Name.Account)
		klog.Debug("end check action")
		c.Next(ctx)
	}
}
