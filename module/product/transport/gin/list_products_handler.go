package gin

import (
	"net/http"
	"strconv"

	ginpkg "github.com/gin-gonic/gin"

	productbiz "stockflow/module/product/biz"
	"stockflow/module/product/model"
	"stockflow/module/product/storage"
)

func ListProductsHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	return func(c *ginpkg.Context) {
		var filter model.Filter
		filter.SKU = c.Query("sku")
		filter.Name = c.Query("name")

		if activeStr := c.Query("is_active"); activeStr != "" {
			val, err := strconv.ParseBool(activeStr)

			if err != nil {
				c.JSON(http.StatusBadRequest, ginpkg.H{
					"error": "is_active must be true or false",
				})

				return
			}

			filter.IsActive = &val
		}

		paging := &model.Paging{
			Page:  parseIntOrDefault(c.Query("page"), 1),
			Limit: parseIntOrDefault(c.Query("limit"), 10),
		}

		biz := productbiz.NewListProductsBiz(store)
		products, err := biz.ListProducts(c.Request.Context(), &filter, paging)

		if err != nil {
			c.JSON(http.StatusInternalServerError, ginpkg.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, ginpkg.H{
			"data": products,
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
