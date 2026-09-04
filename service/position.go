package service

import (
	"context"
	"errors"
	"time"

	"github.com/N30A/trakt/models"
	"github.com/N30A/trakt/repos"
)

var ErrDeviceNotFound = errors.New("device not found")

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
	Protocol       models.Protocol
}

type PositionService struct {
	deviceRepo   *repos.DeviceRepo
	positionRepo *repos.PositionRepo
}

func NewPositionService(deviceRepo *repos.DeviceRepo, positionRepo *repos.PositionRepo) *PositionService {
	return &PositionService{
		deviceRepo:   deviceRepo,
		positionRepo: positionRepo,
	}
}

func (s *PositionService) SavePosition(ctx context.Context, input PositionInput) error {
	device, err := s.deviceRepo.GetDeviceByUniqueID(ctx, input.DeviceUniqueID)
	if err != nil {
		if errors.Is(err, repos.ErrNotFound) {
			return ErrDeviceNotFound
		}

		return err
	}

	position := models.Position{
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
