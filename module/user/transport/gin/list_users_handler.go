package gin

import (
	"net/http"
	"strconv"
	"strings"

	ginpkg "github.com/gin-gonic/gin"

	"stockflow/module/user/biz"
	"stockflow/module/user/model"
	"stockflow/module/user/storage"
)

func ListUsersHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	listUsersBiz := biz.NewListUsersBiz(store)

	return func(c *ginpkg.Context) {
		filter := &model.Filter{
			Email:    strings.TrimSpace(c.Query("email")),
			FullName: strings.TrimSpace(c.Query("full_name")),
			Role:     strings.TrimSpace(c.Query("role")),
		}

		if isActiveStr := strings.TrimSpace(c.Query("is_active")); isActiveStr != "" {
			isActive, err := strconv.ParseBool(isActiveStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, ginpkg.H{"error": "invalid is_active"})
				return
			}
			filter.IsActive = &isActive
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

		users, err := listUsersBiz.ListUsers(c.Request.Context(), filter, paging)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ginpkg.H{
			"data": users,
			"paging": ginpkg.H{
				"page":  paging.Page,
				"limit": paging.Limit,
			},
		})
	}
}
