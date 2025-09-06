package model

import "time"

type SensorReading struct {
	ID         uint64     `db:"id" json:"id"`
	ID1        string     `db:"id1" json:"id1"`
	ID2        int        `db:"id2" json:"id2"`
	SensorType string     `db:"sensor_type" json:"sensor_type"`
	Value      float64    `db:"value" json:"value"`
	TS         time.Time  `db:"ts" json:"ts"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	ArchivedAt *time.Time `db:"archived_at" json:"archived_at,omitempty"`
}
