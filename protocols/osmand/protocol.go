package osmand

// OsmAnd protocol.
//
// This protocol was originally used by the OsmAnd mobile app and was later
// adopted and extended by the official Traccar Client apps.
// The protocol uses HTTP POST requests with parameters provided either
// in the URL query string or as form data (application/x-www-form-urlencoded)
// in the request body.
//
// Not all parameters defined by the protocol are relevant to this server
// and are therefore intentionally ignored.
//
// Parameters used by this server:
// id or deviceid, timestamp, lat, lon, speed, bearing or heading, altitude,
// accuracy
//
// https://www.traccar.org/osmand/

import "github.com/N30A/trakt/models"

const ProtocolOsmAnd models.Protocol = "osmand"
