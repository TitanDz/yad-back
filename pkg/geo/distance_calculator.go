package geo

import (
	"math"
)

const EarthRadiusKm = 6371.0

// CalculateDistance returns the distance in kilometers between two geographic points
// using the Haversine formula. This implementation matches the Dart frontend calculation
// for consistent behavior across client and server.
//
// Parameters:
// - lat1, lon1: User's latitude and longitude
// - lat2, lon2: Target location's latitude and longitude
//
// Returns: Distance in kilometers
//
// Example:
//
//	distance := geo.CalculateDistance(
//		40.7128, -74.0060,  // New York
//		34.0522, -118.2437  // Los Angeles
//	)
//	// Returns approximately 3944 km
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert latitude difference to radians
	dLat := degreesToRadians(lat2 - lat1)
	// Convert longitude difference to radians
	dLon := degreesToRadians(lon2 - lon1)

	// Haversine formula components
	sinHalfDLat := math.Sin(dLat / 2)
	sinHalfDLon := math.Sin(dLon / 2)

	// Calculate 'a' component
	a := sinHalfDLat*sinHalfDLat +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			sinHalfDLon*sinHalfDLon

	// Calculate 'c' component
	c := 2 * math.Asin(math.Sqrt(a))

	// Return distance in kilometers
	return EarthRadiusKm * c
}

// CalculateDistanceInMiles returns the distance in miles between two geographic points.
// This is a convenience method that converts kilometers to miles.
func CalculateDistanceInMiles(lat1, lon1, lat2, lon2 float64) float64 {
	distanceKm := CalculateDistance(lat1, lon1, lat2, lon2)
	return distanceKm / 1.60934
}

// degreesToRadians converts degrees to radians.
// Used internally by the Haversine formula.
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// IsWithinRadius checks if a target point is within a specified radius from a user location.
// This is a convenience method combining distance calculation and comparison.
func IsWithinRadius(userLat, userLon, targetLat, targetLon, radiusKm float64) bool {
	distance := CalculateDistance(userLat, userLon, targetLat, targetLon)
	return distance <= radiusKm
}
