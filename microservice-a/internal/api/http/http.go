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

// UpdateFrequency godoc
// @Summary Update sensor data generation frequency
// @Description Change how often sensor data is generated.The `freq` parameter supports two formats:-- 1. **Milliseconds as integer** (e.g., `1000` = 1 second), 2. **Go duration string** (e.g., `1s`, `500ms`, `2m`) Example usages:- `POST /frequency?freq=1000` → Updates frequency to 1 second, - `POST /frequency?freq=500ms` → Updates frequency to 500 milliseconds
// @Tags MicroserviceA
// @Accept json
// @Produce json
// @Param freq query string true "New frequency for sensor data generation (milliseconds or duration string)" example:"1000"
// @Success 200 {object} map[string]interface{} "Frequency successfully updated"
// @Failure 400 {object} map[string]string "Invalid or missing frequency"
// @Router /frequency [post]
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
