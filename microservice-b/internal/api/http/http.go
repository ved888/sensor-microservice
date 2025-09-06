package http

import (
	"microservice-b/internal/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type SensorHandler struct {
	repo *repository.SensorRepository
}

func NewSensorHandler(repo *repository.SensorRepository) *SensorHandler {
	return &SensorHandler{repo: repo}
}

// GET /sensors
func (h *SensorHandler) GetSensors(c echo.Context) error {
	filters := make(map[string]interface{})
	if id1 := c.QueryParam("id1"); id1 != "" {
		filters["id1"] = id1
	}
	if id2 := c.QueryParam("id2"); id2 != "" {
		filters["id2"] = id2
	}
	loc, _ := time.LoadLocation("Asia/Kolkata")

	if from := c.QueryParam("from"); from != "" {
		t, _ := time.Parse(time.RFC3339, from) // parsed in UTC
		filters["from"] = t.In(loc)            // convert to IST
	}
	if to := c.QueryParam("to"); to != "" {
		t, _ := time.Parse(time.RFC3339, to)
		filters["to"] = t.In(loc)
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 10
	}
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	data, err := h.repo.GetSensors(filters, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Get total count for pagination metadata
	total, err := h.repo.CountSensors(filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := map[string]interface{}{
		"data":        data,
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	}
	return c.JSON(http.StatusOK, response)
}

// DELETE /sensors
func (h *SensorHandler) DeleteSensors(c echo.Context) error {
	filters := make(map[string]interface{})
	if id1 := c.QueryParam("id1"); id1 != "" {
		filters["id1"] = id1
	}
	if id2 := c.QueryParam("id2"); id2 != "" {
		filters["id2"] = id2
	}
	if from := c.QueryParam("from"); from != "" {
		t, _ := time.Parse(time.RFC3339, from)
		filters["from"] = t
	}
	if to := c.QueryParam("to"); to != "" {
		t, _ := time.Parse(time.RFC3339, to)
		filters["to"] = t
	}

	rows, err := h.repo.DeleteSensors(filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"deleted": rows})
}

// PATCH /sensors
func (h *SensorHandler) EditSensors(c echo.Context) error {
	type request struct {
		Value float64 `json:"value"`
	}
	req := new(request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	filters := make(map[string]interface{})
	if id1 := c.QueryParam("id1"); id1 != "" {
		filters["id1"] = id1
	}
	if id2 := c.QueryParam("id2"); id2 != "" {
		filters["id2"] = id2
	}
	loc, _ := time.LoadLocation("Asia/Kolkata")

	if from := c.QueryParam("from"); from != "" {
		t, _ := time.Parse(time.RFC3339, from) // parsed in UTC
		filters["from"] = t.In(loc)            // convert to IST
	}
	if to := c.QueryParam("to"); to != "" {
		t, _ := time.Parse(time.RFC3339, to)
		filters["to"] = t.In(loc)
	}

	rows, err := h.repo.EditSensors(filters, req.Value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"updated": rows})
}
