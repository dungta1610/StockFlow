package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	"stockflow/module/outbox/biz"
	"stockflow/module/outbox/model"
	"stockflow/module/outbox/storage"
)

func MarkProcessedHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	markProcessedBiz := biz.NewMarkProcessedBiz(store)

	return func(c *ginpkg.Context) {
		data := &model.OutboxEventMarkProcessed{
			EventID: c.Param("id"),
		}

		event, err := markProcessedBiz.MarkProcessed(c.Request.Context(), data)
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
