package live

import "math"

const earthRadiusKM = 6371

func HaversineKM(lat1, lng1, lat2, lng2 float64) float64 {
	toRad := func(v float64) float64 { return v * math.Pi / 180 }
	dLat := toRad(lat2 - lat1)
	dLng := toRad(lng2 - lng1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*math.Sin(dLng/2)*math.Sin(dLng/2)
	return 2 * earthRadiusKM * math.Asin(math.Sqrt(a))
}
