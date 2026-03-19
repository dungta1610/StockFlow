package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	"stockflow/module/user/biz"
	"stockflow/module/user/model"
	"stockflow/module/user/storage"
)

func UpdateUserHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	updateUserBiz := biz.NewUpdateUserBiz(store)

	return func(c *ginpkg.Context) {
		id := c.Param("id")

		var data model.UserUpdate
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{"error": err.Error()})
			return
		}

		updatedUser, err := updateUserBiz.UpdateUser(c.Request.Context(), id, &data)
		if err != nil {
			statusCode := http.StatusBadRequest
			if err == model.ErrUserNotFound {
				statusCode = http.StatusNotFound
			}

			c.JSON(statusCode, ginpkg.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ginpkg.H{"data": updatedUser})
	}
}
