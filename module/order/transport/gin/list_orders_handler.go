package gin

import (
	"net/http"
	"strconv"
	"strings"

	ginpkg "github.com/gin-gonic/gin"

	"stockflow/module/order/biz"
	"stockflow/module/order/model"
	"stockflow/module/order/storage"
)

func ListOrdersHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	listOrdersBiz := biz.NewListOrdersBiz(store)

	return func(c *ginpkg.Context) {
		orderCode := strings.TrimSpace(c.Query("order_code"))
		if orderCode == "" {
			orderCode = strings.TrimSpace(c.Query("code"))
		}

		filter := &model.Filter{
			OrderCode:   orderCode,
			UserID:      strings.TrimSpace(c.Query("user_id")),
			WarehouseID: strings.TrimSpace(c.Query("warehouse_id")),
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

		orders, err := listOrdersBiz.ListOrders(c.Request.Context(), filter, paging)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ginpkg.H{
			"data": orders,
			"paging": ginpkg.H{
				"page":  paging.Page,
				"limit": paging.Limit,
			},
		})
	}
}
