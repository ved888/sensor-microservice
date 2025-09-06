package repository

import (
	"fmt"
	"log"
	"microservice-b/model"
	"time"

	pb "microservice-b/pb/shared-proto"

	"github.com/jmoiron/sqlx"
)

type SensorRepository struct {
	DB *sqlx.DB
}

func NewSensorRepository(db *sqlx.DB) *SensorRepository {
	return &SensorRepository{DB: db}
}

func (r *SensorRepository) Save(data *pb.SensorData) error {
	query := `INSERT INTO sensor_readings(
                            value, 
                            sensor_type, 
                            id1, 
                            id2, 
                            ts)
                     VALUES (?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(query, data.Value, data.SensorType, data.Id1, data.Id2, data.Timestamp.AsTime())
	if err != nil {
		log.Printf("Failed to insert sensor data: %v", err)
		return err
	}
	return nil
}

// GetSensors with optional filters and pagination
func (r *SensorRepository) GetSensors(filters map[string]interface{}, limit, offset int) ([]model.SensorReading, error) {
	query := "SELECT * FROM sensor_readings WHERE 1=1"
	args := []interface{}{}

	for k, v := range filters {
		switch k {
		case "id1", "id2":
			query += fmt.Sprintf(" AND %s = ?", k)
			args = append(args, v)
		case "from":
			query += " AND ts >= ?"
			if t, ok := v.(time.Time); ok {
				args = append(args, t)
			}
		case "to":
			query += " AND ts <= ?"
			if t, ok := v.(time.Time); ok {
				args = append(args, t)
			}
		}
	}

	query += " ORDER BY ts DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var sensors []model.SensorReading
	err := r.DB.Select(&sensors, query, args...)
	return sensors, err
}

func (r *SensorRepository) CountSensors(filters map[string]interface{}) (int64, error) {
	query := "SELECT COUNT(*) FROM sensor_readings WHERE 1=1"
	args := []interface{}{}

	for k, v := range filters {
		switch k {
		case "id1", "id2":
			query += fmt.Sprintf(" AND %s = ?", k)
			args = append(args, v)
		case "from":
			query += " AND ts >= ?"
			if t, ok := v.(time.Time); ok {
				args = append(args, t)
			}
		case "to":
			query += " AND ts <= ?"
			if t, ok := v.(time.Time); ok {
				args = append(args, t)
			}
		}
	}

	var total int64
	err := r.DB.Get(&total, query, args...)
	return total, err
}

// DeleteSensors based on filters
func (r *SensorRepository) DeleteSensors(filters map[string]interface{}) (int64, error) {
	query := "DELETE FROM sensor_readings WHERE 1=1"
	args := []interface{}{}

	for k, v := range filters {
		switch k {
		case "id1", "id2":
			query += fmt.Sprintf(" AND %s = ?", k)
			args = append(args, v)
		case "from":
			query += " AND ts >= ?"
			args = append(args, v)
		case "to":
			query += " AND ts <= ?"
			args = append(args, v)
		}
	}

	res, err := r.DB.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// EditSensors update value based on filters
func (r *SensorRepository) EditSensors(filters map[string]interface{}, newValue float64) (int64, error) {
	query := "UPDATE sensor_readings SET value = ? WHERE 1=1"
	args := []interface{}{newValue}

	for k, v := range filters {
		switch k {
		case "id1", "id2":
			query += fmt.Sprintf(" AND %s = ?", k)
			args = append(args, v)
		case "from":
			query += " AND ts >= ?"
			if t, ok := v.(time.Time); ok {
				args = append(args, t)
			}
		case "to":
			query += " AND ts <= ?"
			if t, ok := v.(time.Time); ok {
				args = append(args, t)
			}
		}
	}

	res, err := r.DB.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
