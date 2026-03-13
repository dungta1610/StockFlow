package gin

import (
	ginpkg "github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"stockflow/module/inventory/storage"
)

func RegisterRoutes(r *ginpkg.Engine, db *pgxpool.Pool) {
	store := storage.NewSQLStore(db)

	inventories := r.Group("/inventories")
	{
		inventories.POST("/adjust", AdjustStockHandler(store))
		inventories.GET("", GetInventoryHandler(store))
		inventories.GET("/transactions", ListInventoryTransactionsHandler(store))
	}
}
