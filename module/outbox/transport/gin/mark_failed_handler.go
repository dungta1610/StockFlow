package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	"stockflow/module/outbox/biz"
	"stockflow/module/outbox/model"
	"stockflow/module/outbox/storage"
)

func MarkFailedHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	markFailedBiz := biz.NewMarkFailedBiz(store)

	return func(c *ginpkg.Context) {
		var data model.OutboxEventMarkFailed

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{"error": err.Error()})
			return
		}

		data.EventID = c.Param("id")

		event, err := markFailedBiz.MarkFailed(c.Request.Context(), &data)
		if err != nil {
			statusCode := http.StatusBadRequest
			if err == model.ErrOutboxEventNotFound {
				statusCode = http.StatusNotFound
			}

			c.JSON(statusCode, ginpkg.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ginpkg.H{"data": event})
	}
}
