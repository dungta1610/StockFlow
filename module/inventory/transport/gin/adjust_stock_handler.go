package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"

	inventorybiz "stockflow/module/inventory/biz"
	"stockflow/module/inventory/model"
	"stockflow/module/inventory/storage"
)

func AdjustStockHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	return func(c *ginpkg.Context) {
		var data model.InventoryAdjust

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{
				"error": err.Error(),
			})
			return
		}

		biz := inventorybiz.NewAdjustStockBiz(store)

		inventory, err := biz.AdjustStock(c.Request.Context(), &data)
		if err != nil {
			switch err {
			case model.ErrInventoryAdjustDataRequired,
				model.ErrInventoryProductIDIsBlank,
				model.ErrInventoryWarehouseIDIsBlank,
				model.ErrInventoryAdjustQtyInvalid,
				model.ErrInventoryAvailableQtyInvalid,
				model.ErrInventoryReservedQtyInvalid:
				c.JSON(http.StatusBadRequest, ginpkg.H{
					"error": err.Error(),
				})
				return

			case model.ErrInventoryNotFound:
				c.JSON(http.StatusNotFound, ginpkg.H{
					"error": err.Error(),
				})
				return

			case model.ErrInventoryNotEnoughStock:
				c.JSON(http.StatusConflict, ginpkg.H{
					"error": err.Error(),
				})
				return
			}

			if pgErr, ok := err.(*pgconn.PgError); ok {
				switch pgErr.Code {
				case "23505":
					c.JSON(http.StatusConflict, ginpkg.H{
						"error": model.ErrInventoryAlreadyExists.Error(),
					})
					return
				case "23514":
					c.JSON(http.StatusBadRequest, ginpkg.H{
						"error": "inventory data violates database constraints",
					})
					return
				}
			}

			c.JSON(http.StatusInternalServerError, ginpkg.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, ginpkg.H{
			"data": inventory,
		})
	}
}
