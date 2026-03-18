package gin

import (
	ginpkg "github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"stockflow/module/order/storage"
)

func RegisterRoutes(r *ginpkg.Engine, db *pgxpool.Pool) {
	store := storage.NewSQLStore(db)

	orders := r.Group("/orders")
	{
		orders.POST("", CreateOrderHandler(store))
		orders.GET("", ListOrdersHandler(store))
		orders.GET("/:id", GetOrderHandler(store))
		orders.POST("/:id/cancel", CancelOrderHandler(store))
		orders.POST("/:id/expire", ExpireOrderHandler(store))
	}
}
