package position

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PositionRepo struct {
	pool *pgxpool.Pool
}

func NewPositionRepo(pool *pgxpool.Pool) *PositionRepo {
	return &PositionRepo{pool}
}

func (r *PositionRepo) AddPosition(ctx context.Context, position Position) error {
	query := `
		INSERT INTO positions(
		    device_id, latitude, longitude, fix_time, server_time,
		    protocol, altitude, speed, course, accuracy
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`

	if _, err := r.pool.Exec(ctx, query,
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
	); err != nil {
		return err
	}

	return nil
}

func (r *PositionRepo) DevicePositions(ctx context.Context, deviceIDs []int, from, to time.Time) ([]Position, error) {
	query := `
		SELECT
			id, device_id, latitude, longitude, fix_time,
			server_time, protocol, altitude, speed, course, accuracy
		FROM positions
		WHERE device_id = ANY($1)
		  AND fix_time >= $2
		  AND fix_time < $3
		ORDER BY fix_time ASC
	`

	rows, err := r.pool.Query(ctx, query, deviceIDs, from, to)
	if err != nil {
		return nil, err
	}

	positions, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Position])
	if err != nil {
		return nil, err
	}

	return positions, nil
}

func (r *PositionRepo) LatestPositions(ctx context.Context) ([]Position, error) {
	query := `
		SELECT DISTINCT ON (device_id)
	    id, device_id, latitude, longitude, fix_time,
	    server_time, protocol, altitude, speed, course, accuracy
		FROM positions
		ORDER BY device_id, fix_time DESC;
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	positions, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Position])
	if err != nil {
		return nil, err
	}

	return positions, nil
}
