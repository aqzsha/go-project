package healthcheck

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Deps struct {
	Redis     *redis.Client
	Externals map[string]string
}

// HealthHandler godoc
// @Summary      Check the health status of microservices
// @Description  Verifies connectivity to Redis and internal services
// @Tags         Health
// @Produce      json
// @Router       /health/check [get]
func HealthHandler(d Deps) gin.HandlerFunc {
	client := &http.Client{Timeout: 2 * time.Second}

	return func(c *gin.Context) {
		start := time.Now()
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		var redisErr error

		var wg sync.WaitGroup
		wg.Add(1)

		go func() { defer wg.Done(); redisErr = pingRedis(ctx, d.Redis) }()

		ext := make(map[string]string, len(d.Externals))
		var wgExt sync.WaitGroup
		var mu sync.Mutex

		for name, url := range d.Externals {
			wgExt.Add(1)
			go func(name, url string) {
				defer wgExt.Done()
				if err := pingHTTP(ctx, client, url); err != nil {
					mu.Lock()
					ext[name] = err.Error()
					mu.Unlock()
					return
				}
				mu.Lock()
				ext[name] = "ok"
				mu.Unlock()
			}(name, url)
		}

		wg.Wait()
		wgExt.Wait()

		ok := redisErr == nil
		code := http.StatusOK
		if !ok {
			code = http.StatusServiceUnavailable
		}

		c.JSON(code, gin.H{
			"status":    map[bool]string{true: "ok", false: "fail"}[ok],
			"redis":     errStr(redisErr),
			"externals": ext,
			"took_ms":   time.Since(start).Milliseconds(),
		})

	}
}

func pingRedis(ctx context.Context, rdb *redis.Client) error {
	if rdb == nil {
		return nil
	}
	return rdb.Ping(ctx).Err()
}

func pingHTTP(ctx context.Context, httpc *http.Client, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := httpc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		return fmtError("status %d", resp.StatusCode)
	}
	return nil
}

func errStr(err error) string {
	if err == nil {
		return "ok"
	}

	return err.Error()
}

func fmtError(format string, a ...any) error {
	return &simpleError{msg: fmt.Sprintf(format, a...)}
}

type simpleError struct{ msg string }

func (e *simpleError) Error() string { return e.msg }
