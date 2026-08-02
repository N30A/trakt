package repos

import (
	"context"
	"time"

	"github.com/N30A/trakt/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PositionRepo struct {
	pool *pgxpool.Pool
}

func NewPositionRepo(pool *pgxpool.Pool) *PositionRepo {
	return &PositionRepo{pool}
}

func (r *PositionRepo) AddPosition(ctx context.Context, position models.Position) error {
	query := `
		INSERT INTO positions(
		    device_id, latitude, longitude, fix_time, server_time,
		    protocol, altitude, speed, course, accuracy
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`

	_, err := r.pool.Exec(ctx, query,
		position.DeviceID,
		position.Latitude,
		position.Longitude,
		position.FixTime,
		position.ServerTime,
		position.Protocol,
		position.Altitude,
		position.Speed,
		position.Course,
		position.Accuracy,
	)
	if err != nil {
		return ErrInternal
	}

	return nil
}

func (r *PositionRepo) GetPositionsByDevice(ctx context.Context, deviceID int, from, to time.Time) ([]models.Position, error) {
	query := `
		SELECT
		    id, device_id, latitude, longitude, fix_time,
		    server_time, protocol, altitude, speed, course, accuracy
		FROM positions
		WHERE device_id = $1 AND (fix_time >= $2 AND fix_time <= $3)
		ORDER BY fix_time ASC;
	`

	rows, err := r.pool.Query(ctx, query, deviceID, from, to)
	if err != nil {
		return nil, ErrInternal
	}
	defer rows.Close()

	var positions []models.Position

	for rows.Next() {
		var position models.Position

		if err := rows.Scan(
			&position.ID, &position.DeviceID, &position.Latitude, &position.Longitude,
			&position.FixTime, &position.ServerTime, &position.Protocol, &position.Altitude,
			&position.Speed, &position.Course, &position.Accuracy,
		); err != nil {
			return nil, ErrInternal
		}

		positions = append(positions, position)
	}

	return positions, nil
}

func (r *PositionRepo) GetPositions(ctx context.Context, from, to time.Time) ([]models.Position, error) {
	query := `
		SELECT
		    id, device_id, latitude, longitude, fix_time,
		    server_time, protocol, altitude, speed, course, accuracy
		FROM positions
		WHERE fix_time >= $1 AND fix_time <= $2
		ORDER BY fix_time ASC;
	`

	rows, err := r.pool.Query(ctx, query, from, to)
	if err != nil {
		return nil, ErrInternal
	}
	defer rows.Close()

	var positions []models.Position

	for rows.Next() {
		var position models.Position

		if err := rows.Scan(
			&position.ID, &position.DeviceID, &position.Latitude, &position.Longitude,
			&position.FixTime, &position.ServerTime, &position.Protocol, &position.Altitude,
			&position.Speed, &position.Course, &position.Accuracy,
		); err != nil {
			return nil, ErrInternal
		}

		positions = append(positions, position)
	}

	return positions, nil
}
