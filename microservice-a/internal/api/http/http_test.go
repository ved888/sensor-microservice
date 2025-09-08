package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"microservice-a/internal/api/grpcclient"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_UpdateFrequency_Milliseconds(t *testing.T) {
	// Setup
	e := echo.New()
	gen := grpcclient.NewGenerator("localhost:50051", 1*time.Second)
	handler := NewHandler(gen)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/frequency?freq=1000", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.UpdateFrequency(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Frequency updated", response["message"])
	assert.Equal(t, "1s", response["frequency"])
}

func TestHandler_UpdateFrequency_DurationString(t *testing.T) {
	// Setup
	e := echo.New()
	gen := grpcclient.NewGenerator("localhost:50051", 1*time.Second)
	handler := NewHandler(gen)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/frequency?freq=500ms", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.UpdateFrequency(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Frequency updated", response["message"])
	assert.Equal(t, "500ms", response["frequency"])
}

func TestHandler_UpdateFrequency_Seconds(t *testing.T) {
	// Setup
	e := echo.New()
	gen := grpcclient.NewGenerator("localhost:50051", 1*time.Second)
	handler := NewHandler(gen)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/frequency?freq=2s", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.UpdateFrequency(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Frequency updated", response["message"])
	assert.Equal(t, "2s", response["frequency"])
}

func TestHandler_UpdateFrequency_MissingFreq(t *testing.T) {
	// Setup
	e := echo.New()
	gen := grpcclient.NewGenerator("localhost:50051", 1*time.Second)
	handler := NewHandler(gen)

	// Create request without freq parameter
	req := httptest.NewRequest(http.MethodPost, "/frequency", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.UpdateFrequency(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "freq is required", response["error"])
}

func TestHandler_UpdateFrequency_InvalidFreq(t *testing.T) {
	// Setup
	e := echo.New()
	gen := grpcclient.NewGenerator("localhost:50051", 1*time.Second)
	handler := NewHandler(gen)

	// Create request with invalid freq parameter
	req := httptest.NewRequest(http.MethodPost, "/frequency?freq=invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.UpdateFrequency(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid freq", response["error"])
}

func TestHandler_UpdateFrequency_ZeroFreq(t *testing.T) {
	// Setup
	e := echo.New()
	gen := grpcclient.NewGenerator("localhost:50051", 1*time.Second)
	handler := NewHandler(gen)

	// Create request with zero freq parameter
	req := httptest.NewRequest(http.MethodPost, "/frequency?freq=0", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.UpdateFrequency(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Frequency updated", response["message"])
	assert.Equal(t, "0s", response["frequency"])
}

func TestHandler_UpdateFrequency_NegativeFreq(t *testing.T) {
	// Setup
	e := echo.New()
	gen := grpcclient.NewGenerator("localhost:50051", 1*time.Second)
	handler := NewHandler(gen)

	// Create request with negative freq parameter
	req := httptest.NewRequest(http.MethodPost, "/frequency?freq=-1000", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.UpdateFrequency(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid freq", response["error"])
}
