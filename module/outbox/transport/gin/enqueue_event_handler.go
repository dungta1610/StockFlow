package gin

import (
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	"stockflow/module/outbox/biz"
	"stockflow/module/outbox/model"
	"stockflow/module/outbox/storage"
)

func EnqueueEventHandler(store *storage.SQLStore) ginpkg.HandlerFunc {
	enqueueEventBiz := biz.NewEnqueueEventBiz(store)

	return func(c *ginpkg.Context) {
		var data model.OutboxEventCreate

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{"error": err.Error()})
			return
		}

		event, err := enqueueEventBiz.EnqueueEvent(c.Request.Context(), &data)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, ginpkg.H{"data": event})
	}
}
