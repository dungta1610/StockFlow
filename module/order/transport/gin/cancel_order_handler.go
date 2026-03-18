package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	"stockflow/module/order/biz"
	"stockflow/module/order/model"
	"stockflow/module/order/storage"
)

func CancelOrderHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	cancelOrderBiz := biz.NewCancelOrderBiz(store)

	return func(c *ginpkg.Context) {
		data := &model.OrderCancel{
			OrderID: c.Param("id"),
		}

		updatedOrder, err := cancelOrderBiz.CancelOrder(c.Request.Context(), data)
		if err != nil {
			statusCode := http.StatusBadRequest
			if err == model.ErrOrderNotFound {
				statusCode = http.StatusNotFound
			}

			c.JSON(statusCode, ginpkg.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ginpkg.H{"data": updatedOrder})
	}
}
