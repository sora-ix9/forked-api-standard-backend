package utils

import (
	"time"
)

const (
	// BangkokTimezone is the timezone for Bangkok (UTC+07:00)
	BangkokTimezone = "Asia/Bangkok"
)

var (
	// bangkokLocation is the cached location for Bangkok timezone
	bangkokLocation *time.Location
)

// init initializes the Bangkok timezone location
func init() {
	loc, err := time.LoadLocation(BangkokTimezone)
	if err != nil {
		// Fallback to UTC+07:00 if Asia/Bangkok is not available
		loc = time.FixedZone("UTC+7", 7*60*60)
	}
	bangkokLocation = loc
}

// NowBangkok returns the current time in Bangkok timezone (UTC+07:00)
func NowBangkok() time.Time {
	return time.Now().In(bangkokLocation)
}

// ParseBangkok parses a time string and converts it to Bangkok timezone
func ParseBangkok(layout, value string) (time.Time, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(bangkokLocation), nil
}

// ToBangkok converts a time to Bangkok timezone
func ToBangkok(t time.Time) time.Time {
	return t.In(bangkokLocation)
}

// GetBangkokLocation returns the Bangkok timezone location
func GetBangkokLocation() *time.Location {
	return bangkokLocation
}

// FormatBangkokRFC3339 formats time in Bangkok timezone to RFC3339 format
func FormatBangkokRFC3339(t time.Time) string {
	return t.In(bangkokLocation).Format(time.RFC3339)
}

// SetGlobalTimezone sets the default timezone for the entire application
func SetGlobalTimezone() {
	time.Local = bangkokLocation
}
