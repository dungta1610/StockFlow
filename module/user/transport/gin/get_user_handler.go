package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	"stockflow/module/user/biz"
	"stockflow/module/user/model"
	"stockflow/module/user/storage"
)

func GetUserHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	getUserBiz := biz.NewGetUserBiz(store)

	return func(c *ginpkg.Context) {
		id := c.Param("id")

		user, err := getUserBiz.GetUser(c.Request.Context(), id)
		if err != nil {
			statusCode := http.StatusBadRequest
			if err == model.ErrUserNotFound {
				statusCode = http.StatusNotFound
			}

			c.JSON(statusCode, ginpkg.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ginpkg.H{"data": user})
	}
}
