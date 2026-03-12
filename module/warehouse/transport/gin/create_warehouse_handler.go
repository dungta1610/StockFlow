package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"

	warehousebiz "stockflow/module/warehouse/biz"
	"stockflow/module/warehouse/model"
	"stockflow/module/warehouse/storage"
)

func CreateWarehouseHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	return func(c *ginpkg.Context) {
		var data model.WarehouseCreate

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{
				"error": err.Error(),
			})
			return
		}

		biz := warehousebiz.NewCreateWarehouseBiz(store)
		warehouse, err := biz.CreateWarehouse(c.Request.Context(), &data)

		if err != nil {
			if err == model.ErrWarehouseDataIsNil ||
				err == model.ErrWarehouseCodeIsBlank ||
				err == model.ErrWarehouseNameIsBlank {
				c.JSON(http.StatusBadRequest, ginpkg.H{
					"error": err.Error(),
				})

				return
			}

			if err == model.ErrWarehouseCodeAlreadyExists {
				c.JSON(http.StatusConflict, ginpkg.H{
					"error": err.Error(),
				})

				return
			}

			if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, ginpkg.H{
					"error": model.ErrWarehouseCodeAlreadyExists.Error(),
				})

				return
			}

			c.JSON(http.StatusInternalServerError, ginpkg.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusCreated, ginpkg.H{
			"data": warehouse,
		})
	}
}
