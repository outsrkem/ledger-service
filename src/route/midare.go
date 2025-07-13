package route

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/app"
	"io"
	"ledger/src/config"
	"ledger/src/pkg/answer"
	"ledger/src/pkg/uuid"
	"ledger/src/slog"
	"math"
	"net/http"
	"time"
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
		xRequestId := c.Request.Header.Get(xRequestIdKey)
		klog.Infof("start check action, [%s] [%s]", xRequestId, action)

		xAuthToken := c.Request.Header.Get(xAuthTokenKey)
		if xAuthToken == "" {
			klog.Error(xAuthTokenKey + " is empty.")
			c.JSON(http.StatusForbidden, answer.ResBody(answer.EcodeInvalidTokenError, xAuthTokenKey+" is empty.", ""))
			c.Abort()
			return
		}

		type actionRaw struct {
			Uias struct {
				Action string `json:"action"`
			} `json:"uias"`
		}
		var raw actionRaw
		raw.Uias.Action = action
		rawJson, err := json.Marshal(raw)
		if err != nil {
			klog.Errorf("Error marshaling audit log: %v", err)
			c.Abort()
			return
		}

		// url := "https://uias-devops.endpoint.outsrkem.top:30078/v1/uias/action/check"
		url := config.AppCfg.Ledger.Uias.Endpoint + "/v1/uias/action/check"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(rawJson))
		if err != nil {
			klog.Errorf("Error creating request: %v", err)
			c.JSON(http.StatusForbidden, answer.ResBody(answer.EcodeInvalidTokenError, "Internal service error.", ""))
			c.Abort()
			return
		}

		req.Header.Set(xRequestIdKey, xRequestId)
		req.Header.Set(xAuthTokenKey, xAuthToken)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: config.AppCfg.Ledger.Uias.SkipTlsVerify}}
		client := &http.Client{Timeout: 30 * time.Second, Transport: transport}
		resp, err := client.Do(req)
		if err != nil {
			klog.Errorf("Error sending req log: %v", err)
			c.JSON(http.StatusForbidden, answer.ResBody(answer.EcodeInvalidTokenError, "Internal service error.", ""))
			c.Abort()
			return
		}
		defer func() {
			if resp != nil {
				if err := resp.Body.Close(); err != nil {
					klog.Errorf("Close request failed: %v", err)
				}
			}
		}()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			klog.Error("io.ReadAll", err)
			c.JSON(http.StatusForbidden, answer.ResBody(answer.EcodeInvalidTokenError, "Internal service error.", ""))
			c.Abort()
			return
		}

		var result map[string]interface{}
		if resp.StatusCode != http.StatusOK {
			klog.Errorf("Request failed with status code %d", resp.StatusCode)
			klog.Infof("check action url: %s", url)
			c.JSON(resp.StatusCode, result)
			c.Abort()
			return
		}

		if err := json.Unmarshal(body, &result); err != nil {
			klog.Warn("json Unmarshal err: ", err)
			klog.Error(result)
			c.JSON(http.StatusForbidden, answer.ResBody(answer.EcodeInvalidTokenError, "Internal service error.", ""))
			c.Abort()
			return
		}

		klog.Debugf("result %v", result)
		// 尝试获取 payload 数据
		payload, ok := result["payload"].(map[string]interface{})
		if !ok {
			klog.Warn("Failed to get payload from response")
			c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeUpstreamResponseError, "Internal service error.", ""))
			return
		}

		// 尝试获取 authentication 数据
		authentication, ok := payload["authentication"].(float64)
		if !ok {
			klog.Warn("Failed to get authentication from payload")
			c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeUpstreamResponseError, "Internal service error.", ""))
			return
		}

		// 防御性检查：确保是整数值,上游始终返回 authentication 是整数
		if math.Mod(authentication, 1) != 0 {
			klog.Warnf("Unexpected non-integer value for authentication: %f", authentication)
			c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeUpstreamResponseError, "Invalid authentication value", nil))
			return
		}
		// 转换为整数并进行全权限比较
		if int(authentication) != 1 {
			// 没有权限，返回403和上游返回体，便于查看问题
			klog.Warnf("Permission denial. result: %+v", result)
			c.JSON(403, result)
			c.Abort()
			return
		}

		// 尝试获取 user 数据
		user, ok := payload["user"].(map[string]interface{})
		if !ok {
			klog.Warn("Failed to get userinfo from payload")
			c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeUpstreamResponseError, "Internal service error.", ""))
			return
		}
		userId, ok := user["id"].(string)
		if !ok {
			klog.Warn("Failed to get user ID from userinfo")
			c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeUpstreamResponseError, "Internal service error.", ""))
			return
		}
		// 安全获取嵌套的account字段
		name, ok := user["name"].(map[string]interface{})
		if !ok {
			klog.Warn("Failed to get user name structure")
			c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeUpstreamResponseError, "Invalid user information", ""))
			return
		}
		account, ok := name["account"].(string)
		if !ok {
			klog.Warn("Failed to get user account")
			c.JSON(http.StatusInternalServerError, answer.ResBody(answer.EcodeUpstreamResponseError, "Invalid user information", ""))
			return
		}

		klog.Info("Permission is granted, and the operation is authorized.")
		c.Set("userId", userId)
		c.Set("account", account)
		klog.Debug("end check action")
		c.Next(ctx)
	}
}
