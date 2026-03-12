package gin

import (
	ginpkg "github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"stockflow/module/product/storage"
)

func RegisterRoutes(r *ginpkg.Engine, db *pgxpool.Pool) {
	store := storage.NewSQLStore(db)

	products := r.Group("/products")
	{
		products.POST("", CreateProductHandler(store))
		products.GET("", ListProductsHandler(store))
		products.GET("/:id", GetProductHandler(store))
	}
}
