package position

import (
	"time"

	"github.com/N30A/trakt/protocol"
)

type Position struct {
	ID        int
	DeviceID  int
	Latitude  float64
	Longitude float64

	//DeviceTime time.Time
	FixTime    time.Time // UTC
	ServerTime time.Time // UTC

	Protocol protocol.Protocol

	Altitude *float64
	Speed    *float64
	Course   *float64
	Accuracy *float64
}
