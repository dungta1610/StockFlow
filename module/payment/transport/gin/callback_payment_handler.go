package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	"stockflow/module/payment/biz"
	"stockflow/module/payment/model"
	"stockflow/module/payment/storage"
)

func CallbackPaymentHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	callbackPaymentBiz := biz.NewCallbackPaymentBiz(store)

	return func(c *ginpkg.Context) {
		var data model.PaymentCallback

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{"error": err.Error()})
			return
		}

		updatedPayment, err := callbackPaymentBiz.CallbackPayment(c.Request.Context(), &data)
		if err != nil {
			statusCode := http.StatusBadRequest
			if err == model.ErrPaymentNotFound {
				statusCode = http.StatusNotFound
			}

			c.JSON(statusCode, ginpkg.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ginpkg.H{"data": updatedPayment})
	}
}
