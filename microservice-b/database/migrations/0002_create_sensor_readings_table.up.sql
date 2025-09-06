CREATE TABLE sensor_readings (
                                 id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                                 id1 VARCHAR(16) NOT NULL,
                                 id2 INT NOT NULL,
                                 sensor_type VARCHAR(32) NOT NULL,
                                 value DOUBLE NOT NULL,
                                 ts DATETIME(6) NOT NULL,
                                 created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
                                 updated_at DATETIME(6) NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(6),
                                 archived_at DATETIME(6) NULL DEFAULT NULL,
                                 INDEX IX_id_combo (id1, id2),
                                 INDEX IX_type_ts (sensor_type, ts),
                                 INDEX IX_combo_ts (id1, id2, ts),
                                 INDEX IX_ts (ts)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
