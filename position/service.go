package position

import (
	"context"
	"errors"
	"time"

	"github.com/N30A/trakt/device"
	"github.com/N30A/trakt/protocol"
)

type PositionInput struct {
	DeviceUniqueID string
	Latitude       float64
	Longitude      float64
	FixTime        time.Time // UTC
	ServerTime     time.Time // UTC
	Speed          *float64
	Course         *float64
	Altitude       *float64
	Accuracy       *float64
	Protocol       protocol.Protocol
}

type PositionService struct {
	deviceRepo   *device.DeviceRepo
	positionRepo *PositionRepo
}

func NewPositionService(deviceRepo *device.DeviceRepo, positionRepo *PositionRepo) *PositionService {
	return &PositionService{
		deviceRepo:   deviceRepo,
		positionRepo: positionRepo,
	}
}

func (s *PositionService) SavePosition(ctx context.Context, input PositionInput) error {
	device, err := s.deviceRepo.GetDeviceByUniqueID(ctx, input.DeviceUniqueID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}

		return err
	}

	position := Position{
		DeviceID:   device.ID,
		Latitude:   input.Latitude,
		Longitude:  input.Longitude,
		FixTime:    input.FixTime,
		ServerTime: input.ServerTime,
		Protocol:   input.Protocol,
		Altitude:   input.Altitude,
		Speed:      input.Speed,
		Course:     input.Course,
		Accuracy:   input.Accuracy,
	}

	return s.positionRepo.AddPosition(ctx, position)
}
