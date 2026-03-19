package gin

import (
	"net/http"
	"strconv"
	"strings"

	ginpkg "github.com/gin-gonic/gin"

	"stockflow/module/payment/biz"
	"stockflow/module/payment/model"
	"stockflow/module/payment/storage"
)

func ListPaymentsHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	listPaymentsBiz := biz.NewListPaymentsBiz(store)

	return func(c *ginpkg.Context) {
		filter := &model.Filter{
			OrderID:     strings.TrimSpace(c.Query("order_id")),
			PaymentCode: strings.TrimSpace(c.Query("payment_code")),
			Method:      strings.TrimSpace(c.Query("method")),
			Status:      strings.TrimSpace(c.Query("status")),
		}

		paging := model.NewPaging()

		if pageStr := strings.TrimSpace(c.Query("page")); pageStr != "" {
			page, err := strconv.Atoi(pageStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, ginpkg.H{"error": "invalid page"})
				return
			}
			paging.Page = page
		}

		if limitStr := strings.TrimSpace(c.Query("limit")); limitStr != "" {
			limit, err := strconv.Atoi(limitStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, ginpkg.H{"error": "invalid limit"})
				return
			}
			paging.Limit = limit
		}

		payments, err := listPaymentsBiz.ListPayments(c.Request.Context(), filter, paging)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ginpkg.H{
			"data": payments,
			"paging": ginpkg.H{
				"page":  paging.Page,
				"limit": paging.Limit,
			},
		})
	}
}
