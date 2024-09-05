/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package ratelimiter provides functionality for logging operation-related information.
package ratelimiter

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	redislib "github.com/opensourceways/redis-lib"
	"github.com/sirupsen/logrus"
	"github.com/throttled/throttled/v2"
	"github.com/throttled/throttled/v2/store/goredisstore.v8"
	"golang.org/x/xerrors"

	commonctl "github.com/opensourceways/xihe-server/common/controller"
	"github.com/opensourceways/xihe-server/common/domain/allerror"
)

const (
	userIdParsed          = "user_id"
	defaultClientPoolSize = 10
	defaultIdleTimeOutNum = 30
	// RequestNumPerSec represents the maximum number of requests allowed per second.
	RequestNumPerSec = 10
	// BurstNumPerSec represents the maximum number of burst requests allowed per second.
	BurstNumPerSec   = 10
	maxCASMultiplier = 100
)

var (
	limiter *rateLimiter
)

// InitRateLimiter creates a new instance of the rateLimiter struct.
func InitRateLimiter(cfg redislib.Config, rateCfg *Config) error {
	Init(rateCfg)
	// Initialize a redis client using go-redis
	var client *redis.Client
	if cfg.DBCert != "" {
		ca, err := os.ReadFile(cfg.DBCert)
		if err != nil {
			return fmt.Errorf("read cert failed")
		}

		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(ca) {
			return fmt.Errorf("new pool failed")
		}

		tlsConfig := &tls.Config{
			InsecureSkipVerify: true, // #nosec G402
			RootCAs:            pool,
		}

		client = redis.NewClient(&redis.Options{
			PoolSize:    defaultClientPoolSize, // default
			IdleTimeout: defaultIdleTimeOutNum * time.Second,
			DB:          cfg.DB,
			Addr:        cfg.Address,
			Password:    cfg.Password,
			TLSConfig:   tlsConfig,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			PoolSize:    defaultClientPoolSize, // default
			IdleTimeout: defaultIdleTimeOutNum * time.Second,
			DB:          cfg.DB,
			Addr:        cfg.Address,
			Password:    cfg.Password,
		})
	}
	// Setup store
	store, err := goredisstore.NewCtx(client, "api-rate-limit:")
	if err != nil {
		return xerrors.Errorf("init goredisstore failed: %w", err)
	}

	requestNum := RequestNumPerSec
	if config.RequestNum > 0 {
		requestNum = config.RequestNum
	}
	burstNum := BurstNumPerSec
	if config.BurstNum > 0 {
		burstNum = config.BurstNum
	}
	// Setup quota
	quota := throttled.RateQuota{
		MaxRate:  throttled.PerSec(requestNum),
		MaxBurst: burstNum,
	}
	rateLimitCtx, err := throttled.NewGCRARateLimiterCtx(store, quota)
	if err != nil {
		return xerrors.Errorf("init NewGCRARateLimiterCtx failed: %w", err)
	}
	// set max cas limit value because maxCASAttempts must be more over than requestNum
	maxCASAttempts := requestNum * maxCASMultiplier
	rateLimitCtx.SetMaxCASAttemptsLimit(maxCASAttempts)
	logrus.Infof(" ratelimit with: rate: %d burst: %d", requestNum, burstNum)

	httpRateLimiter := &throttled.HTTPRateLimiterCtx{
		RateLimiter: rateLimitCtx,
	}

	limiter = &rateLimiter{limitCli: httpRateLimiter}

	return nil
}

type rateLimiter struct {
	limitCli *throttled.HTTPRateLimiterCtx
}

// Limiter is a function that returns a pointer to the rateLimiter instance.
func Limiter() *rateLimiter {
	return limiter
}

func (rl *rateLimiter) CheckLimit(ctx *gin.Context) {
	v, ok := ctx.Get(userIdParsed)
	if !ok {
		ctx.Next()
		return
	}

	key := fmt.Sprintf("%v", v)
	limited, _, err := rl.limitCli.RateLimiter.RateLimitCtx(ctx.Request.Context(), key, 1)
	if err != nil {
		commonctl.SendError(ctx, allerror.NewOverLimit(allerror.ErrorRateLimitOver, "too many requests", err))
		ctx.Abort()
		return
	}

	if limited {
		commonctl.SendError(ctx, allerror.NewOverLimit(allerror.ErrorRateLimitOver, "too many requests",
			fmt.Errorf("over limit")))
		ctx.Abort()
	} else {
		ctx.Next()
	}
}
