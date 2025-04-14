package route

import (
    "testing"

	"github.com/majori/wrc-laptimer/pkg/telemetry"
	"github.com/majori/wrc-laptimer/internal/database"
)

type SpyRoutePointDb struct {
	routePoints []*database.RoutePoint
}

func (s *SpyRoutePointDb) SetRoutePoints(routeID int, routePoints []*database.RoutePoint) error {
	s.routePoints = routePoints
	return nil
}

func (s *SpyRoutePointDb) GetRoutePoints(routeID int) ([]*database.RoutePoint, error) {
	return s.routePoints, nil
}

func generateTelemetry(size int, stageLength int) []*telemetry.TelemetrySessionUpdate {
	telemetrySamples := make([]*telemetry.TelemetrySessionUpdate, size)
	for i := 0; i < size; i++ {
		telemetrySamples[i] = &telemetry.TelemetrySessionUpdate{
			StageCurrentDistance: float64(stageLength) / float64(size) * float64(i),
			VehiclePositionX:     float32(i),
			VehiclePositionY:     float32(i),
			VehiclePositionZ:     float32(i),
			VehicleVelocityX:     float32(i),
			VehicleVelocityY:     float32(i),
			VehicleVelocityZ:     float32(i),
			VehicleAccelerationX: float32(i),
			VehicleAccelerationY: float32(i),
			VehicleAccelerationZ: float32(i),
			VehicleThrottle:      float32(i),
			VehicleEngineRpmCurrent:   float32(((i * 100) % 6500) + 2000),
			VehicleEngineRpmIdle:      2000,
			VehicleEngineRpmMax:       8500,
			VehicleBrake:              float32(i),
			VehicleHandbrake:          float32(i),
			VehicleGearIndex:          uint8(i % 6),
			VehicleGearIndexNeutral:   0,
			VehicleGearIndexReverse:   10,
			VehicleGearMaximum:        6,
		}
	}
	return telemetrySamples
}

func TestRouteDataUpdate(t *testing.T) {
	// Create a new SpyRoutePointDb instance
	spyDb := &SpyRoutePointDb{}

	telemetrySamples := generateTelemetry(100, 20)

	for i := 0; i < len(telemetrySamples); i++ {
		ProcessTelemetryForRouteData(telemetrySamples[i])
	}

	UpdateRoutePoints(spyDb, 1, 1)

	// Check if the route points were set correctly
	if len(spyDb.routePoints) != 20 {
		t.Errorf("Expected 20 route points, got %d", len(spyDb.routePoints))
	}
}
