package gin

import (
	ginpkg "github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"stockflow/module/user/storage"
)

func RegisterRoutes(r *ginpkg.Engine, db *pgxpool.Pool) {
	store := storage.NewSQLStore(db)

	users := r.Group("/users")
	{
		users.POST("", CreateUserHandler(store))
		users.GET("", ListUsersHandler(store))
		users.GET("/:id", GetUserHandler(store))
		users.PUT("/:id", UpdateUserHandler(store))
	}
}
