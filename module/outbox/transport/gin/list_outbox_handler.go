package gin

import (
	"net/http"
	"strconv"
	"strings"

	ginpkg "github.com/gin-gonic/gin"

	"stockflow/module/outbox/biz"
	"stockflow/module/outbox/model"
	"stockflow/module/outbox/storage"
)

func ListOutboxHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	listPendingEventsBiz := biz.NewListPendingEventsBiz(store)

	return func(c *ginpkg.Context) {
		filter := &model.Filter{
			AggregateType: strings.TrimSpace(c.Query("aggregate_type")),
			AggregateID:   strings.TrimSpace(c.Query("aggregate_id")),
			EventType:     strings.TrimSpace(c.Query("event_type")),
			Status:        strings.TrimSpace(c.Query("status")),
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

		events, err := listPendingEventsBiz.ListPendingEvents(c.Request.Context(), filter, paging)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ginpkg.H{
			"data": events,
			"paging": ginpkg.H{
				"page":  paging.Page,
				"limit": paging.Limit,
			},
		})
	}
}
