package httpx

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), CorrelationMiddleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	return r
}

func CorrelationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cid := c.GetHeader("X-Correlation-Id")
		if cid == "" {
			cid = uuid.New().String()
		}
		c.Set("correlation_id", cid)
		c.Writer.Header().Set("X-Correlation-Id", cid)
		c.Next()
	}
}

func RegisterHealth(r *gin.Engine, live, ready func() error) {
	r.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/health/ready", func(c *gin.Context) {
		if ready != nil {
			if err := ready(); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})
}

func Run(r *gin.Engine, port string) error {
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return srv.ListenAndServe()
}

func CorrelationID(c *gin.Context) string {
	if v, ok := c.Get("correlation_id"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return uuid.New().String()
}
