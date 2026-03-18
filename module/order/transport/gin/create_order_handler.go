package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	"stockflow/module/order/biz"
	"stockflow/module/order/model"
	"stockflow/module/order/storage"
)

func CreateOrderHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	createOrderBiz := biz.NewCreateOrderBiz(store)

	return func(c *ginpkg.Context) {
		var data model.OrderCreate

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{"error": err.Error()})
			return
		}

		createdOrder, err := createOrderBiz.CreateOrder(c.Request.Context(), &data)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, ginpkg.H{"data": createdOrder})
	}
}
