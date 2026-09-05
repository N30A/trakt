package device

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeviceRepo struct {
	pool *pgxpool.Pool
}

func NewDeviceRepo(pool *pgxpool.Pool) *DeviceRepo {
	return &DeviceRepo{pool}
}

func (r *DeviceRepo) GetDeviceByID(ctx context.Context, deviceID int) (Device, error) {
	query := `
		SELECT id, unique_id, name
		FROM devices
		WHERE id = $1;
	`

	var device Device

	err := r.pool.QueryRow(ctx, query, deviceID).Scan(&device.ID, &device.UniqueID, &device.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Device{}, ErrNotFound
		}

		return Device{}, ErrInternal
	}

	return device, nil
}

func (r *DeviceRepo) GetDeviceByUniqueID(ctx context.Context, uniqueID string) (Device, error) {
	query := `
		SELECT id, unique_id, name
		FROM devices
		WHERE unique_id = $1;
	`

	var device Device

	err := r.pool.QueryRow(ctx, query, uniqueID).Scan(&device.ID, &device.UniqueID, &device.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Device{}, ErrNotFound
		}

		return Device{}, ErrInternal
	}

	return device, nil
}

func (r *DeviceRepo) GetDevices(ctx context.Context) ([]Device, error) {
	query := `
		SELECT id, unique_id, name
		FROM devices
		ORDER BY id ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	devices, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Device])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return devices, nil
}

func (r *DeviceRepo) DeleteDeviceByID(ctx context.Context, deviceID int) error {
	query := `
		DELETE FROM devices
		WHERE id = $1;
	`

	tag, err := r.pool.Exec(ctx, query, deviceID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: %v", ErrNotFound, err)
	}
	return nil
}

func (r *DeviceRepo) AddDevice(ctx context.Context, newDevice Device) (Device, error) {
	query := `
		INSERT INTO devices (unique_id, name)
		VALUES ($1, $2)
		RETURNING id, unique_id, name;
	`

	var device Device
	if err := r.pool.QueryRow(ctx, query, newDevice.UniqueID, newDevice.Name).Scan(
		&device.ID, &device.UniqueID, &device.Name,
	); err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" {
			return Device{}, fmt.Errorf("%w: %v", ErrConflict, err)
		}

		return Device{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return device, nil
}

func (r *DeviceRepo) UpdateDevice(ctx context.Context, updatedDevice Device) (Device, error) {
	query := `
		UPDATE devices
		SET unique_id = $1, name = $2
		WHERE id = $3
		RETURNING id, unique_id, name;
	`

	var device Device
	if err := r.pool.QueryRow(ctx, query, updatedDevice.UniqueID, updatedDevice.Name, updatedDevice.ID).Scan(
		&device.ID, &device.UniqueID, &device.Name,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Device{}, fmt.Errorf("%w: %v", ErrNotFound, err)
		}

		return Device{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return device, nil
}
