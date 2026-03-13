package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	inventorybiz "stockflow/module/inventory/biz"
	"stockflow/module/inventory/model"
	"stockflow/module/inventory/storage"
)

func GetInventoryHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	return func(c *ginpkg.Context) {
		id := c.Query("id")
		productID := c.Query("product_id")
		warehouseID := c.Query("warehouse_id")

		biz := inventorybiz.NewGetInventoryBiz(store)

		inventory, err := biz.GetInventory(c.Request.Context(), id, productID, warehouseID)
		if err != nil {
			switch err {
			case model.ErrInventoryProductIDIsBlank,
				model.ErrInventoryWarehouseIDIsBlank:
				c.JSON(http.StatusBadRequest, ginpkg.H{
					"error": err.Error(),
				})
				return

			case model.ErrInventoryNotFound:
				c.JSON(http.StatusNotFound, ginpkg.H{
					"error": err.Error(),
				})
				return
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
