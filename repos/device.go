package repos

import (
	"context"
	"errors"

	"github.com/N30A/trakt/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeviceRepo struct {
	pool *pgxpool.Pool
}

func NewDeviceRepo(pool *pgxpool.Pool) *DeviceRepo {
	return &DeviceRepo{pool}
}

func (r *DeviceRepo) GetDeviceByID(ctx context.Context, deviceID int) (models.Device, error) {
	query := `
		SELECT id, unique_id, name
		FROM devices
		WHERE id = $1;
	`

	var device models.Device

	err := r.pool.QueryRow(ctx, query, deviceID).Scan(&device.ID, &device.UniqueID, &device.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Device{}, errors.Join(ErrNotFound, err)
		}

		return models.Device{}, errors.Join(ErrInternal, err)
	}

	return device, nil
}

func (r *DeviceRepo) GetDeviceByUniqueID(ctx context.Context, uniqueID string) (models.Device, error) {
	query := `
		SELECT id, unique_id, name
		FROM devices
		WHERE unique_id = $1;
	`

	var device models.Device

	err := r.pool.QueryRow(ctx, query, uniqueID).Scan(&device.ID, &device.UniqueID, &device.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Device{}, errors.Join(ErrNotFound, err)
		}

		return models.Device{}, errors.Join(ErrInternal, err)
	}

	return device, nil
}

func (r *DeviceRepo) GetDevices(ctx context.Context) ([]models.Device, error) {
	query := `
		SELECT id, unique_id, name
		FROM devices
		ORDER BY id ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errors.Join(ErrInternal, err)
	}
	defer rows.Close()

	var devices []models.Device

	for rows.Next() {
		var device models.Device

		if err := rows.Scan(&device.ID, &device.UniqueID, &device.Name); err != nil {
			return nil, errors.Join(ErrInternal, err)
		}

		devices = append(devices, device)
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
		return errors.Join(ErrInternal, err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
