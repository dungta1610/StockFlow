package gin

import (
	ginpkg "github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"stockflow/module/outbox/storage"
)

func RegisterRoutes(r *ginpkg.Engine, db *pgxpool.Pool) {
	store := storage.NewSQLStore(db)

	outbox := r.Group("/outbox")
	{
		outbox.POST("/events", EnqueueEventHandler(store))
		outbox.GET("/events", ListPendingEventsHandler(store))
		outbox.POST("/events/:id/processed", MarkProcessedHandler(store))
		outbox.POST("/events/:id/failed", MarkFailedHandler(store))
	}
}
