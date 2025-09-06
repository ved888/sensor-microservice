package http

import (
	"net/http"
	"strconv"
	"time"

	"microservice-a/internal/api/grpcclient"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	generator *grpcclient.Generator
}

func NewHandler(gen *grpcclient.Generator) *Handler {
	return &Handler{generator: gen}
}

func (h *Handler) UpdateFrequency(c echo.Context) error {
	freqStr := c.QueryParam("freq")
	if freqStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "freq is required"})
	}

	var duration time.Duration
	// parsing as integer (ms)
	if ms, err := strconv.Atoi(freqStr); err == nil && ms > 0 {
		duration = time.Duration(ms) * time.Millisecond
	} else {
		// parsing as time duration string
		d, err := time.ParseDuration(freqStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid freq"})
		}
		duration = d
	}

	h.generator.UpdateFrequency(duration)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "Frequency updated",
		"frequency": duration.String(),
	})
}
