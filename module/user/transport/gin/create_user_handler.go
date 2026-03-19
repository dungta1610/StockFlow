package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	"stockflow/module/user/biz"
	"stockflow/module/user/model"
	"stockflow/module/user/storage"
)

func CreateUserHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	createUserBiz := biz.NewCreateUserBiz(store)

	return func(c *ginpkg.Context) {
		var data model.UserCreate

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{"error": err.Error()})
			return
		}

		createdUser, err := createUserBiz.CreateUser(c.Request.Context(), &data)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, ginpkg.H{"data": createdUser})
	}
}
