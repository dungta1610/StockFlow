package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	warehousebiz "stockflow/module/warehouse/biz"
	"stockflow/module/warehouse/model"
	"stockflow/module/warehouse/storage"
)

func GetWarehouseHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	return func(c *ginpkg.Context) {
		id := c.Param("id")
		biz := warehousebiz.NewGetWarehouseBiz(store)
		warehouse, err := biz.GetWarehouse(c.Request.Context(), id)

		if err != nil {
			if err == model.ErrWarehouseIDIsBlank {
				c.JSON(http.StatusBadRequest, ginpkg.H{
					"error": err.Error(),
				})

				return
			}

			if err == model.ErrWarehouseNotFound {
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
			"data": warehouse,
		})
	}
}
