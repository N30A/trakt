package osmand

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

type osmandPosition struct {
	DeviceUniqueID string
	Latitude       float64
	Longitude      float64
	Timestamp      time.Time // seconds or milliseconds since epoch, ISO 8601 format, or "yyyy-MM-dd HH:mm:ss"
	Speed          *float64  // in knots
	Course         *float64  // in degrees
	Altitude       *float64  // in meters
	Accuracy       *float64  // in meters
}

func parsePosition(r *http.Request) (osmandPosition, error) {
	deviceID, ok := getDeviceIDParam(r)
	if !ok {
		return osmandPosition{}, errors.New("missing required id or deviceid param")
	}

	latitude, err := strconv.ParseFloat(r.FormValue("lat"), 64)
	if err != nil {
		return osmandPosition{}, errors.New("missing required lat param")
	}

	longitude, err := strconv.ParseFloat(r.FormValue("lon"), 64)
	if err != nil {
		return osmandPosition{}, errors.New("missing required lon param")
	}

	timestamp, ok := getTimestampParam(r)
	if !ok {
		return osmandPosition{}, errors.New("missing required timestamp param")
	}

	speed := getOptionalParam(r, "speed", func(key string) (float64, error) {
		return strconv.ParseFloat(key, 64)
	})

	course := getOptionalCourseParam(r)

	altitude := getOptionalParam(r, "altitude", func(key string) (float64, error) {
		return strconv.ParseFloat(key, 64)
	})

	accuracy := getOptionalParam(r, "accuracy", func(key string) (float64, error) {
		return strconv.ParseFloat(key, 64)
	})

	return osmandPosition{
		DeviceUniqueID: deviceID,
		Latitude:       latitude,
		Longitude:      longitude,
		Timestamp:      timestamp,
		Speed:          speed,
		Course:         course,
		Altitude:       altitude,
		Accuracy:       accuracy,
	}, nil
}

func getDeviceIDParam(r *http.Request) (string, bool) {
	for _, key := range []string{"id", "deviceid"} {
		if value := r.FormValue(key); value != "" {
			return value, true
		}
	}
	return "", false
}

func getTimestampParam(r *http.Request) (time.Time, bool) {
	value := r.FormValue("timestamp")
	if value == "" {
		return time.Time{}, false
	}

	// Unix timestamp (seconds eller milliseconds)
	if number, err := strconv.ParseInt(value, 10, 64); err == nil {
		if number > 1_000_000_000_000 {
			return time.UnixMilli(number).UTC(), true
		}

		return time.Unix(number, 0).UTC(), true
	}

	// ISO 8601 / RFC3339
	if time, err := time.Parse(time.RFC3339, value); err == nil {
		return time.UTC(), true
	}

	// "yyyy-MM-dd HH:mm:ss"
	if time, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.UTC); err == nil {
		return time, true
	}

	return time.Time{}, false
}

func getOptionalCourseParam(r *http.Request) *float64 {
	for _, key := range []string{"bearing", "heading"} {
		if value := r.FormValue(key); value != "" {
			if course, err := strconv.ParseFloat(value, 64); err == nil {
				return &course
			}
		}
	}
	return nil
}

func getOptionalParam[T any](r *http.Request, key string, parse func(string) (T, error)) *T {
	if value, err := parse(r.FormValue(key)); err == nil {
		return &value
	}
	return nil
}
