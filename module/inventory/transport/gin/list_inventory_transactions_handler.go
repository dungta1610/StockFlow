package gin

import (
	"net/http"
	"strconv"

	ginpkg "github.com/gin-gonic/gin"

	inventorybiz "stockflow/module/inventory/biz"
	"stockflow/module/inventory/model"
	"stockflow/module/inventory/storage"
)

func ListInventoryTransactionsHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	return func(c *ginpkg.Context) {
		filter := &model.TransactionFilter{
			InventoryID:   c.Query("inventory_id"),
			ProductID:     c.Query("product_id"),
			WarehouseID:   c.Query("warehouse_id"),
			OrderID:       c.Query("order_id"),
			ReservationID: c.Query("reservation_id"),
			TxnType:       c.Query("txn_type"),
		}

		paging := &model.Paging{
			Page:  parseIntOrDefault(c.Query("page"), 1),
			Limit: parseIntOrDefault(c.Query("limit"), 10),
		}

		biz := inventorybiz.NewListInventoryTransactionsBiz(store)

		items, err := biz.ListInventoryTransactions(c.Request.Context(), filter, paging)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ginpkg.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, ginpkg.H{
			"data": items,
			"paging": ginpkg.H{
				"page":  paging.Page,
				"limit": paging.Limit,
			},
		})
	}
}

func parseIntOrDefault(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}

	return v
}
