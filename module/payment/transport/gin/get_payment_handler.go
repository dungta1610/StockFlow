package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	"stockflow/module/payment/biz"
	"stockflow/module/payment/model"
	"stockflow/module/payment/storage"
)

func GetPaymentHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	getPaymentBiz := biz.NewGetPaymentBiz(store)

	return func(c *ginpkg.Context) {
		id := c.Param("id")

		payment, err := getPaymentBiz.GetPayment(c.Request.Context(), id)
		if err != nil {
			statusCode := http.StatusBadRequest
			if err == model.ErrPaymentNotFound {
				statusCode = http.StatusNotFound
			}

			c.JSON(statusCode, ginpkg.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ginpkg.H{"data": payment})
	}
}
