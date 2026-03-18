package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	"stockflow/module/order/biz"
	"stockflow/module/order/storage"
)

func GetOrderHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	getOrderBiz := biz.NewGetOrderBiz(store)

	return func(c *ginpkg.Context) {
		id := c.Param("id")

		order, err := getOrderBiz.GetOrder(c.Request.Context(), id)
		if err != nil {
			statusCode := http.StatusBadRequest
			if err.Error() == "order not found" {
				statusCode = http.StatusNotFound
			}

			c.JSON(statusCode, ginpkg.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ginpkg.H{"data": order})
	}
}
