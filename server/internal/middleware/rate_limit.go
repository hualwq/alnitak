package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/resp"
)

// TokenBucket 令牌桶结构
type TokenBucket struct {
	capacity   int        // 桶容量
	rate       int        // 令牌产生速率（个/秒）
	tokens     int        // 当前令牌数量
	lastRefill time.Time  // 上次补充令牌的时间
	mu         sync.Mutex // 互斥锁
}

// NewTokenBucket 创建新的令牌桶
func NewTokenBucket(capacity, rate int) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		rate:       rate,
		tokens:     capacity,
		lastRefill: time.Now(),
	}
}

// Take 尝试获取令牌
func (tb *TokenBucket) Take() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// 补充令牌
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	tokensToAdd := int(elapsed.Seconds() * float64(tb.rate))

	if tokensToAdd > 0 {
		tb.tokens = min(tb.capacity, tb.tokens+tokensToAdd)
		tb.lastRefill = now
	}

	// 检查是否有可用令牌
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RateLimit 限流中间件
func RateLimit() gin.HandlerFunc {
	// 创建令牌桶
	bucket := NewTokenBucket(
		global.Config.Security.RateLimit.Capacity,
		global.Config.Security.RateLimit.Rate,
	)

	return func(c *gin.Context) {
		// 如果限流未启用，直接通过
		if !global.Config.Security.RateLimit.Enabled {
			c.Next()
			return
		}

		// 尝试获取令牌
		if bucket.Take() {
			c.Next()
		} else {
			// 限流，返回错误
			resp.FailWithMessage(c, "请求过于频繁，请稍后再试")
			c.Abort()
		}
	}
}

// EmailRateLimit 邮箱验证码限流中间件
func EmailRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取邮箱地址
		var req struct {
			Email string `json:"email" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			resp.FailWithMessage(c, "请求参数有误")
			c.Abort()
			return
		}

		// 检查邮箱限流
		key := "email_rate_limit:" + req.Email
		count := global.Redis.Get(key)

		if count != "" {
			// 已经达到限流
			resp.FailWithMessage(c, "邮箱验证码发送过于频繁，请稍后再试")
			c.Abort()
			return
		}

		// 设置限流标记
		timeout := time.Duration(global.Config.Security.RateLimit.Timeout) * time.Second
		global.Redis.Set(key, "1", timeout)

		c.Next()
	}
}
