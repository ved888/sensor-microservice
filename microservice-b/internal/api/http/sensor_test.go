package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"microservice-b/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSensorHandler_GetSensors(t *testing.T) {
	// Setup
	e := echo.New()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := repository.NewSensorRepository(sqlxDB)
	handler := NewSensorHandler(repo)

	// Test data
	rows := sqlmock.NewRows([]string{"id", "id1", "id2", "sensor_type", "value", "ts", "created_at"}).
		AddRow(1, "A", 1, "Temperature", 25.5, time.Now(), time.Now()).
		AddRow(2, "A", 1, "Temperature", 26.0, time.Now(), time.Now())

	mock.ExpectQuery("SELECT.*FROM sensor_readings").
		WillReturnRows(rows)

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery("SELECT COUNT.*FROM sensor_readings").
		WillReturnRows(countRows)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/sensors?limit=10&page=1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err = handler.GetSensors(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "data")
	assert.Contains(t, response, "page")
	assert.Contains(t, response, "limit")
	assert.Contains(t, response, "total")

	data := response["data"].([]interface{})
	assert.Len(t, data, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSensorHandler_EditSensors_InvalidBody(t *testing.T) {
	// Setup
	e := echo.New()
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := repository.NewSensorRepository(sqlxDB)
	handler := NewSensorHandler(repo)

	// Create request with invalid JSON
	req := httptest.NewRequest(http.MethodPatch, "/api/sensors?id1=A", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err = handler.EditSensors(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid body", response["error"])
}