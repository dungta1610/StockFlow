package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	"stockflow/module/order/biz"
	"stockflow/module/order/model"
	"stockflow/module/order/storage"
)

type cancelOrderRequest struct {
	Reason string `json:"reason"`
	By     string `json:"by"`
}

func CancelOrderHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	cancelOrderBiz := biz.NewCancelOrderBiz(store)

	return func(c *ginpkg.Context) {
		orderID := c.Param("id")

		var req cancelOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{"error": err.Error()})
			return
		}

		data := &model.OrderCancel{
			OrderID: orderID,
			Reason:  req.Reason,
			By:      req.By,
		}

		updatedOrder, err := cancelOrderBiz.CancelOrder(c.Request.Context(), data)
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
