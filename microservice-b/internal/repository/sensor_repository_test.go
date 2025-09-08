package repository

import (
	pb "microservice-b/pb/shared-proto"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSensorRepository_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := NewSensorRepository(sqlxDB)

	// Test data
	sensorData := &pb.SensorData{
		Value:      25.5,
		SensorType: "Temperature",
		Id1:        "A",
		Id2:        "1",
		Timestamp:  timestamppb.Now(),
	}

	// Mock expectations
	mock.ExpectExec("INSERT INTO sensor_readings").
		WithArgs(25.5, "Temperature", "A", "1", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute
	err = repo.Save(sensorData)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSensorRepository_GetSensors(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := NewSensorRepository(sqlxDB)

	// Prepare expected rows
	rows := sqlmock.NewRows([]string{"id1", "id2", "ts", "value"}).
		AddRow("A", "1", time.Date(2025, 9, 8, 10, 0, 0, 0, time.UTC), 42.5).
		AddRow("B", "2", time.Date(2025, 9, 8, 10, 1, 0, 0, time.UTC), 50.0)

	// Filters
	from := time.Date(2025, 9, 8, 10, 0, 0, 0, time.UTC)
	to := time.Date(2025, 9, 8, 11, 0, 0, 0, time.UTC)
	filters := map[string]interface{}{
		"id1":  "A",
		"from": from,
		"to":   to,
	}

	limit := 10
	offset := 0

	// Expected SQL query pattern
	expectedQuery := "SELECT \\* FROM sensor_readings WHERE 1=1.*ORDER BY ts DESC LIMIT \\? OFFSET \\?"
	// Expect the query
	mock.ExpectQuery(expectedQuery).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), limit, offset).
		WillReturnRows(rows)

	// Call the method
	sensors, err := repo.GetSensors(filters, limit, offset)
	require.NoError(t, err)
	require.Len(t, sensors, 2)
	require.Equal(t, "A", sensors[0].ID1)
	require.Equal(t, 42.5, sensors[0].Value)
	require.Equal(t, "B", sensors[1].ID1)
	require.Equal(t, 50.0, sensors[1].Value)

	// Ensure all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSensorRepository_CountSensors(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := NewSensorRepository(sqlxDB)

	// Test data
	filters := map[string]interface{}{
		"id1": "A",
	}

	// Mock expectations
	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
	mock.ExpectQuery("SELECT COUNT.*FROM sensor_readings").
		WithArgs("A").
		WillReturnRows(rows)

	// Execute
	count, err := repo.CountSensors(filters)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSensorRepository_DeleteSensors(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := NewSensorRepository(sqlxDB)

	// Test data
	filters := map[string]interface{}{
		"id1": "A",
	}

	// Mock expectations
	mock.ExpectExec("DELETE FROM sensor_readings").
		WithArgs("A").
		WillReturnResult(sqlmock.NewResult(0, 3))

	// Execute
	deleted, err := repo.DeleteSensors(filters)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, int64(3), deleted)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSensorRepository_EditSensors(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := NewSensorRepository(sqlxDB)

	// Test data
	filters := map[string]interface{}{
		"id1": "A",
	}
	newValue := 30.0

	// Mock expectations
	mock.ExpectExec("UPDATE sensor_readings SET value").
		WithArgs(30.0, "A").
		WillReturnResult(sqlmock.NewResult(0, 2))

	// Execute
	updated, err := repo.EditSensors(filters, newValue)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, int64(2), updated)
	assert.NoError(t, mock.ExpectationsWereMet())
}
