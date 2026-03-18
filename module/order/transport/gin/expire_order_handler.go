package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	"stockflow/module/order/biz"
	"stockflow/module/order/model"
	"stockflow/module/order/storage"
)

type expireOrderRequest struct {
	Reason string `json:"reason"`
	By     string `json:"by"`
}

func ExpireOrderHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	expireOrderBiz := biz.NewExpireOrderBiz(store)

	return func(c *ginpkg.Context) {
		orderID := c.Param("id")

		var req expireOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{"error": err.Error()})
			return
		}

		data := &model.OrderExpire{
			OrderID: orderID,
			Reason:  req.Reason,
			By:      req.By,
		}

		updatedOrder, err := expireOrderBiz.ExpireOrder(c.Request.Context(), data)
		if err != nil {
			statusCode := http.StatusBadRequest
			if err.Error() == "order not found" {
				statusCode = http.StatusNotFound
			}

			c.JSON(statusCode, ginpkg.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ginpkg.H{"data": updatedOrder})
	}
}
