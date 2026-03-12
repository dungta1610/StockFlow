package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	productbiz "stockflow/module/product/biz"
	"stockflow/module/product/model"
	"stockflow/module/product/storage"
)

func GetProductHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	return func(c *ginpkg.Context) {
		id := c.Param("id")
		biz := productbiz.NewGetProductBiz(store)
		product, err := biz.GetProduct(c.Request.Context(), id)

		if err != nil {
			if err == model.ErrProductIDIsBlank {
				c.JSON(http.StatusBadRequest, ginpkg.H{
					"error": err.Error(),
				})

				return
			}

			if err == model.ErrProductNotFound {
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
			"data": product,
		})
	}
}
