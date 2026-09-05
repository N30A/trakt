package osmand

import (
	"time"

	"github.com/N30A/trakt/position"
)

func toPositionInput(pos osmandPosition, serverTime time.Time) position.PositionInput {
	return position.PositionInput{
		DeviceUniqueID: pos.DeviceUniqueID,
		Latitude:       pos.Latitude,
		Longitude:      pos.Longitude,
		FixTime:        pos.Timestamp,
		ServerTime:     serverTime,
		Speed:          pos.Speed,
		Course:         pos.Course,
		Altitude:       pos.Altitude,
		Accuracy:       pos.Accuracy,
		Protocol:       ProtocolOsmAnd,
	}
}
