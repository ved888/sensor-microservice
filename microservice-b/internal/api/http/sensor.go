package http

import (
	"microservice-b/internal/repository"
	"microservice-b/model"
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

// GetSensors godoc
// @Summary Retrieve sensor readings with filters
// @Description This endpoint retrieves sensor readings from the database.You can filter results by `id1`, `id2`, or by a time range (`from`, `to`).You can also combine filters (e.g., ID1 + time range).Pagination is supported via `page` and `limit` query parameters.- `page`: Page number starting from 1- `limit`: Number of records per page (default: 10) Time parameters must be in RFC3339 format (UTC). Example: `2025-09-06T15:04:05Z`
// @Tags MicroserviceB
// @Accept json
// @Produce json
// @Param id1 query string false "Filter by ID1 (string identifier)" example("A")
// @Param id2 query int false "Filter by ID2 (integer identifier)" example(1)
// @Param from query string false "Filter from timestamp (RFC3339 format)" example("2025-09-06T10:00:00Z")
// @Param to query string false "Filter to timestamp (RFC3339 format)" example("2025-09-06T12:00:00Z")
// @Param page query int false "Page number (starting from 1)" default(1) example(1)
// @Param limit query int false "Page size (number of records per page)" default(10) example(10)
// @Success 200 {object} map[string]interface{} "Paginated sensor readings with metadata"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/sensors [get]
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

// DeleteSensors godoc
// @Summary Delete sensor readings with filters
// @Description This endpoint deletes sensor readings from the database.You can filter records by `id1`, `id2`, or by a time range (`from`, `to`).You can also combine filters (e.g., ID1 + time range).Time parameters must be in RFC3339 format (UTC).Example: `2025-09-06T15:04:05Z`If no filters are provided, **no rows will be deleted**.
// @Tags MicroserviceB
// @Accept json
// @Produce json
// @Param id1 query string false "Filter by ID1 (string identifier)" example("A")
// @Param id2 query int false "Filter by ID2 (integer identifier)" example(1)
// @Param from query string false "Filter from timestamp (RFC3339 format)" example("2025-09-06T10:00:00Z")
// @Param to query string false "Filter to timestamp (RFC3339 format)" example("2025-09-06T12:00:00Z")
// @Success 200 {object} map[string]interface{} "Number of deleted rows, e.g. {\"deleted\": 3}"
// @Failure 500 {object} map[string]string "Internal server error, e.g. {\"error\": \"database failure\"}"
// @Security BearerAuth
// @Router /api/sensors [delete]
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

// EditSensors godoc
// @Summary Update sensor readings values with filters
// @Description This endpoint allows updating sensor values based on optional filters such as `id1`, `id2`, and a time range (`from`, `to`).If no filters are provided, no rows will be updated.Time parameters should be provided in RFC3339 format (UTC).Example time format: `2025-09-06T15:04:05Z`
// @Tags MicroserviceB
// @Accept json
// @Produce json
// @Param id1 query string false "Filter by ID1 (string identifier)" example("A")
// @Param id2 query string false "Filter by ID2 (integer identifier)" example(1)
// @Param from query string false "Start timestamp in RFC3339 format (e.g., 2025-09-06T10:00:00Z)"
// @Param to query string false "End timestamp in RFC3339 format (e.g., 2025-09-06T12:00:00Z)"
// @Param payload body model.EditSensorsRequest true "Sensor update request payload"
// @Success 200 {object} map[string]interface{} "Number of updated rows, e.g. {\"updated\": 5}"
// @Failure 400 {object} map[string]string "Invalid request body, e.g. {\"error\": \"invalid body\"}"
// @Failure 500 {object} map[string]string "Internal server error, e.g. {\"error\": \"database failure\"}"
// @Security BearerAuth
// @Router /api/sensors [patch]
func (h *SensorHandler) EditSensors(c echo.Context) error {
	req := new(model.EditSensorsRequest)
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
