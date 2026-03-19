package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	"stockflow/module/payment/biz"
	"stockflow/module/payment/model"
	"stockflow/module/payment/storage"
)

func CheckoutPaymentHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	checkoutPaymentBiz := biz.NewCheckoutPaymentBiz(store)

	return func(c *ginpkg.Context) {
		var data model.PaymentCheckout

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{"error": err.Error()})
			return
		}

		createdPayment, err := checkoutPaymentBiz.CheckoutPayment(c.Request.Context(), &data)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, ginpkg.H{"data": createdPayment})
	}
}
