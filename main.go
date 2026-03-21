package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	ginpkg "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"stockflow/component/postgres"
	"stockflow/component/ratelimit"
	rediscmp "stockflow/component/redis"
	"stockflow/middleware"

	inventorygin "stockflow/module/inventory/transport/gin"
	ordergin "stockflow/module/order/transport/gin"
	paymentgin "stockflow/module/payment/transport/gin"
	productgin "stockflow/module/product/transport/gin"
	usergin "stockflow/module/user/transport/gin"
	warehousegin "stockflow/module/warehouse/transport/gin"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is required")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, postgres.Config{
		DSN: dsn,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	redisClient, err := rediscmp.NewClient(ctx, rediscmp.Config{
		Addr: redisAddr,
		DB:   0,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	r := ginpkg.Default()

	globalLimiter := ratelimit.NewRedisLimiter(redisClient, "rl")
	r.Use(middleware.RateLimit(globalLimiter, 100, time.Minute))

	r.GET("/health", func(c *ginpkg.Context) {
		if err := pool.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, ginpkg.H{
				"status": "down",
				"error":  err.Error(),
			})
			return
		}

		if err := redisClient.Ping(c.Request.Context()).Err(); err != nil {
			c.JSON(http.StatusServiceUnavailable, ginpkg.H{
				"status": "down",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, ginpkg.H{"status": "ok"})
	})

	productgin.RegisterRoutes(r, pool)
	warehousegin.RegisterRoutes(r, pool)
	inventorygin.RegisterRoutes(r, pool)
	ordergin.RegisterRoutes(r, pool)
	paymentgin.RegisterRoutes(r, pool)
	usergin.RegisterRoutes(r, pool)

	log.Printf("server is running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
