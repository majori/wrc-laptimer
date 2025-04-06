package route

import (
	"math"

	"github.com/majori/wrc-laptimer/pkg/telemetry"
	"github.com/majori/wrc-laptimer/internal/database"
)

type RouteUpdateTelemetrySample struct {
	StageCurrentDistance      int
	VehiclePositionX          float32
	VehiclePositionY          float32
	VehiclePositionZ          float32
	VehicleVelocityX          float32
	VehicleVelocityY          float32
	VehicleVelocityZ          float32
	VehicleAccelerationX      float32
	VehicleAccelerationY      float32
	VehicleAccelerationZ      float32
	VehicleThrottle           float32
	VehicleEngineRpmCurrent   float32
	VehicleEngineRpmIdle      float32
	VehicleEngineRpmMax       float32
	VehicleBrake              float32
	VehicleHandbrake          float32
	VehicleGearIndex          uint8
	VehicleGearIndexNeutral   uint8
	VehicleGearIndexReverse   uint8
	VehicleGearMaximum        uint8
}

const (
	ROUTE_PATH_POINT_INTERVAL = 1 // TODO: read from an env var
)

var (
	previousTelemetry *telemetry.TelemetrySessionUpdate
	samples = make([]RouteUpdateTelemetrySample, 0)
)

func getVectorLen(a float32, b float32, c float32) float32 {
	return float32(math.Sqrt(float64(a*a + b*b + c*c)))
}

// addRoutePoint adds a new route point to the telemetry samples.
// It interpolates the values between two telemetry samples to
// create a new sample at the specified stage distance.
// The stageDistance has to be between the two samples.
func addRoutePoint(
	stageDistance int,
	sample1 *telemetry.TelemetrySessionUpdate,
	sample2 *telemetry.TelemetrySessionUpdate,
) {
	diffTo1 := float64(stageDistance) - sample1.StageCurrentDistance
	diff := sample2.StageCurrentDistance - sample1.StageCurrentDistance

	sample2Weight := float32(1)
	if diff > 0 {
		sample2Weight = float32(diffTo1 / diff)
	}

	samples = append(samples, RouteUpdateTelemetrySample{
		StageCurrentDistance:    stageDistance,
		VehiclePositionX:        (1-sample2Weight) * sample1.VehiclePositionX + sample2Weight * sample2.VehiclePositionX,
		VehiclePositionY:        (1-sample2Weight) * sample1.VehiclePositionY + sample2Weight * sample2.VehiclePositionY,
		VehiclePositionZ:        (1-sample2Weight) * sample1.VehiclePositionZ + sample2Weight * sample2.VehiclePositionZ,
		VehicleVelocityX:        (1-sample2Weight) * sample1.VehicleVelocityX + sample2Weight * sample2.VehicleVelocityX,
		VehicleVelocityY:        (1-sample2Weight) * sample1.VehicleVelocityY + sample2Weight * sample2.VehicleVelocityY,
		VehicleVelocityZ:        (1-sample2Weight) * sample1.VehicleVelocityZ + sample2Weight * sample2.VehicleVelocityZ,
		VehicleAccelerationX:    (1-sample2Weight) * sample1.VehicleAccelerationX + sample2Weight * sample2.VehicleAccelerationX,
		VehicleAccelerationY:    (1-sample2Weight) * sample1.VehicleAccelerationY + sample2Weight * sample2.VehicleAccelerationY,
		VehicleAccelerationZ:    (1-sample2Weight) * sample1.VehicleAccelerationZ + sample2Weight * sample2.VehicleAccelerationZ,
		VehicleThrottle:         sample2.VehicleThrottle,
		VehicleEngineRpmCurrent: sample2.VehicleEngineRpmCurrent,
		VehicleEngineRpmIdle:    sample2.VehicleEngineRpmIdle,
		VehicleEngineRpmMax:     sample2.VehicleEngineRpmMax,
		VehicleBrake:            sample2.VehicleBrake,
		VehicleHandbrake:        sample2.VehicleHandbrake,
		VehicleGearIndex:        sample2.VehicleGearIndex,
		VehicleGearIndexNeutral: sample2.VehicleGearIndexNeutral,
		VehicleGearIndexReverse: sample2.VehicleGearIndexReverse,
		VehicleGearMaximum:      sample2.VehicleGearMaximum,
	})
}

func clearSamples() {
	samples = samples[:0]
}


func ProcessTelemetryForRouteData(sample *telemetry.TelemetrySessionUpdate) {
	if previousTelemetry == nil || len(samples) == 0 {
		previousTelemetry = sample
		addRoutePoint(0, sample, sample)
	}

	nextStageDistance := samples[len(samples)-1].StageCurrentDistance + ROUTE_PATH_POINT_INTERVAL

	if sample.StageCurrentDistance < float64(nextStageDistance) {
		return
	}

	addRoutePoint(nextStageDistance, previousTelemetry, sample)
}

func UpdateRoutePoints(db *database.Database, routeID int, placement int) error {
	if previousTelemetry == nil {
		return nil
	}

	previousTelemetry = nil

	defer clearSamples()

	oldRoutePoints, err := db.GetRoutePoints(routeID)
	if err != nil {
		return err
	}

	newSampleMult := 1.0 / (float32(placement) * 2.0)

	oldSampleMult := float64(len(oldRoutePoints)) / float64(len(samples))
	newRoutePoints := make([]*database.RoutePoint, len(samples))
	for i := 0; i < len(samples); i++ {
		prevSample := oldRoutePoints[int(math.Floor(float64(i) * oldSampleMult))]
		nextSample := oldRoutePoints[int(math.Ceil(float64(i) * oldSampleMult))]

		newRoutePoints[i] = &database.RoutePoint{
			RouteID:         routeID,
			StageDistance:   i * ROUTE_PATH_POINT_INTERVAL,
			X:               (prevSample.X + nextSample.X) / 2.0,
			Y:               (prevSample.Y + nextSample.Y) / 2.0,
			Z:               (prevSample.Z + nextSample.Z) / 2.0,
			AvgVelocity:     (prevSample.AvgVelocity + nextSample.AvgVelocity) / 2.0,
			AvgAcceleration: (prevSample.AvgAcceleration + nextSample.AvgAcceleration) / 2.0,
			AvgThrottle:     (prevSample.AvgThrottle + nextSample.AvgThrottle) / 2.0,
			AvgRpm:          (prevSample.AvgRpm + nextSample.AvgRpm) / 2.0,
			AvgBreaking:     (prevSample.AvgBreaking + nextSample.AvgBreaking) / 2.0,
			AvgHandbreak:    (prevSample.AvgHandbreak + nextSample.AvgHandbreak) / 2.0,
			AvgGear:         (prevSample.AvgGear + nextSample.AvgGear) / 2.0,
		}

		newRoutePoints[i].X = (1-newSampleMult) * newRoutePoints[i].X + newSampleMult * samples[i].VehiclePositionX
		newRoutePoints[i].Y = (1-newSampleMult) * newRoutePoints[i].Y + newSampleMult * samples[i].VehiclePositionY
		newRoutePoints[i].Z = (1-newSampleMult) * newRoutePoints[i].Z + newSampleMult * samples[i].VehiclePositionZ
		newRoutePoints[i].AvgVelocity =
			(1-newSampleMult) * newRoutePoints[i].AvgVelocity +
			newSampleMult * getVectorLen(
				samples[i].VehicleVelocityX,
				samples[i].VehicleVelocityY,
				samples[i].VehicleVelocityZ,
			)
		newRoutePoints[i].AvgAcceleration =
			(1-newSampleMult) * newRoutePoints[i].AvgAcceleration +
			newSampleMult * getVectorLen(
				samples[i].VehicleAccelerationX,
				samples[i].VehicleAccelerationY,
				samples[i].VehicleAccelerationZ,
			)
		newRoutePoints[i].AvgThrottle = (1-newSampleMult) * newRoutePoints[i].AvgThrottle + newSampleMult * samples[i].VehicleThrottle
		newRoutePoints[i].AvgRpm =
			(1-newSampleMult) * newRoutePoints[i].AvgRpm +
			newSampleMult * (samples[i].VehicleEngineRpmCurrent / samples[i].VehicleEngineRpmMax)
		newRoutePoints[i].AvgBreaking = (1-newSampleMult) * newRoutePoints[i].AvgBreaking + newSampleMult * samples[i].VehicleBrake
		newRoutePoints[i].AvgHandbreak = (1-newSampleMult) * newRoutePoints[i].AvgHandbreak + newSampleMult * samples[i].VehicleHandbrake

		gear := samples[i].VehicleGearIndex
		if gear == samples[i].VehicleGearIndexNeutral || gear == samples[i].VehicleGearIndexReverse {
			gear = 0
		}
		newRoutePoints[i].AvgGear =
			(1-newSampleMult) * newRoutePoints[i].AvgGear +
			newSampleMult * (float32(gear) / float32(samples[i].VehicleGearMaximum))
	}


	err = db.SetRoutePoints(routeID, newRoutePoints)
	if err != nil {
		return err
	}

	return nil
}