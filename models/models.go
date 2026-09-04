package models

import "time"

type Protocol string

type Device struct {
	ID       int
	UniqueID string
	Name     string
}

type Position struct {
	ID        int
	DeviceID  int
	Latitude  float64
	Longitude float64

	//DeviceTime time.Time
	FixTime    time.Time // UTC
	ServerTime time.Time // UTC

	Protocol Protocol

	Altitude *float64
	Speed    *float64
	Course   *float64
	Accuracy *float64
}
