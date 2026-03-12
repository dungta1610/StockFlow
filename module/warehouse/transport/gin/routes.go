package gin

import (
	ginpkg "github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"stockflow/module/warehouse/storage"
)

func RegisterRoutes(r *ginpkg.Engine, db *pgxpool.Pool) {
	store := storage.NewSQLStore(db)

	warehouses := r.Group("/warehouses")
	{
		warehouses.POST("", CreateWarehouseHandler(store))
		warehouses.GET("", ListWarehousesHandler(store))
		warehouses.GET("/:id", GetWarehouseHandler(store))
	}
}
