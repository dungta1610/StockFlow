package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"

	productbiz "stockflow/module/product/biz"
	"stockflow/module/product/model"
	"stockflow/module/product/storage"
)

func CreateProductHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	return func(c *ginpkg.Context) {
		var data model.ProductCreate

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{
				"error": err.Error(),
			})

			return
		}

		biz := productbiz.NewCreateProductBiz(store)
		product, err := biz.CreateProduct(c.Request.Context(), &data)

		if err != nil {
			if err == model.ErrProductDataIsNil ||
				err == model.ErrProductSKUIsBlank ||
				err == model.ErrProductNameIsBlank ||
				err == model.ErrProductPriceInvalid {
				c.JSON(http.StatusBadRequest, ginpkg.H{
					"error": err.Error(),
				})

				return
			}

			if err == model.ErrProductSKUAlreadyExists {
				c.JSON(http.StatusConflict, ginpkg.H{
					"error": err.Error(),
				})

				return
			}

			if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, ginpkg.H{
					"error": model.ErrProductSKUAlreadyExists.Error(),
				})

				return
			}

			c.JSON(http.StatusInternalServerError, ginpkg.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusCreated, ginpkg.H{
			"data": product,
		})
	}
}
