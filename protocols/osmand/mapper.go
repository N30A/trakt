package osmand

import (
	"time"

	"github.com/N30A/trakt/service"
)

func toPositionInput(position osmandPosition, serverTime time.Time) service.PositionInput {
	return service.PositionInput{
		DeviceUniqueID: position.DeviceUniqueID,
		Latitude:       position.Latitude,
		Longitude:      position.Longitude,
		FixTime:        position.Timestamp,
		ServerTime:     serverTime,
		Speed:          position.Speed,
		Course:         position.Course,
		Altitude:       position.Altitude,
		Accuracy:       position.Accuracy,
		Protocol:       ProtocolOsmAnd,
	}
}
