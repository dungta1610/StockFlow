package gin

import (
	ginpkg "github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"stockflow/module/payment/storage"
)

func RegisterRoutes(r *ginpkg.Engine, db *pgxpool.Pool) {
	store := storage.NewSQLStore(db)

	payments := r.Group("/payments")
	{
		payments.POST("/checkout", CheckoutPaymentHandler(store))
		payments.POST("/callback", CallbackPaymentHandler(store))
		payments.GET("", ListPaymentsHandler(store))
		payments.GET("/:id", GetPaymentHandler(store))
	}
}
