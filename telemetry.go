package main

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Header struct {
	Packet4CC [4]byte
	PacketUid uint64
}

type TelemetrySessionStart struct {
	VehicleID             uint8
	VehicleClassID        uint8
	VehicleManufacturerID uint8
	GameMode              uint8
	LocationID            uint8
	RouteID               uint8
	StageLength           float64
	StageShakedown        bool
}

type TelemetrySessionUpdate struct {
	StageCurrentDistance      float64
	StageCurrentTime          float32
	StagePreviousSplitTime    float32
	StageProgress             float32
	VehicleAccelerationX      float32
	VehicleAccelerationY      float32
	VehicleAccelerationZ      float32
	VehicleBrake              float32
	VehicleBrakeTemperatureBl float32
	VehicleBrakeTemperatureBr float32
	VehicleBrakeTemperatureFl float32
	VehicleBrakeTemperatureFr float32
	VehicleClutch             float32
	VehicleClusterAbs         bool
	VehicleCpForwardSpeedBl   float32
	VehicleCpForwardSpeedBr   float32
	VehicleCpForwardSpeedFl   float32
	VehicleCpForwardSpeedFr   float32
	VehicleEngineRpmCurrent   float32
	VehicleEngineRpmIdle      float32
	VehicleEngineRpmMax       float32
	VehicleForwardDirectionX  float32
	VehicleForwardDirectionY  float32
	VehicleForwardDirectionZ  float32
	VehicleGearIndex          uint8
	VehicleGearIndexNeutral   uint8
	VehicleGearIndexReverse   uint8
	VehicleGearMaximum        uint8
	VehicleHandbrake          float32
	VehicleHubPositionBl      float32
	VehicleHubPositionBr      float32
	VehicleHubPositionFl      float32
	VehicleHubPositionFr      float32
	VehicleHubVelocityBl      float32
	VehicleHubVelocityBr      float32
	VehicleHubVelocityFl      float32
	VehicleHubVelocityFr      float32
	VehicleLeftDirectionX     float32
	VehicleLeftDirectionY     float32
	VehicleLeftDirectionZ     float32
	VehiclePositionX          float32
	VehiclePositionY          float32
	VehiclePositionZ          float32
	VehicleSpeed              float32
	VehicleSteering           float32
	VehicleThrottle           float32
	VehicleTransmissionSpeed  float32
	VehicleTyreStateBl        uint8
	VehicleTyreStateBr        uint8
	VehicleTyreStateFl        uint8
	VehicleTyreStateFr        uint8
	VehicleUpDirectionX       float32
	VehicleUpDirectionY       float32
	VehicleUpDirectionZ       float32
	VehicleVelocityX          float32
	VehicleVelocityY          float32
	VehicleVelocityZ          float32
}

type TelemetrySessionPause struct {
	StageCurrentTime     float32
	StageCurrentDistance float64
}

type TelemetrySessionResume struct {
	StageCurrentTime     float32
	StageCurrentDistance float64
}

type TelemetrySessionEnd struct {
	Header
	StageResultTime        float32
	StageResultTimePenalty float32
	StageResultStatus      uint8
}

const (
	PacketHeaderSize = 16

	// Packet type constants
	Packet4CCSessionStart  = "sess"
	Packet4CCSessionUpdate = "sesu"
	Packet4CCSessionPause  = "sesp"
	Packet4CCSessionResume = "sesr"
	Packet4CCSessionEnd    = "sese"
)

var (
	ErrInvalidPacket     = errors.New("invalid packet")
	ErrUnknownPacketType = errors.New("unknown packet type")
)

func UnmarshalBinary(data []byte) (any, error) {
	// Check if data is large enough to contain at least a header
	if len(data) < PacketHeaderSize {
		return nil, ErrInvalidPacket
	}

	// Create a reader for binary data
	buf := bytes.NewReader(data)

	var header Header
	if err := binary.Read(buf, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	// Convert 4CC to string for easier comparison
	fourCCStr := string(header.Packet4CC[:])

	var packet any
	switch fourCCStr {
	case Packet4CCSessionStart:
		packet = new(TelemetrySessionStart)

	case Packet4CCSessionUpdate:
		packet = new(TelemetrySessionUpdate)

	case Packet4CCSessionPause:
		packet = new(TelemetrySessionPause)

	case Packet4CCSessionResume:
		packet = new(TelemetrySessionResume)

	case Packet4CCSessionEnd:
		packet = new(TelemetrySessionEnd)

	default:
		return nil, ErrUnknownPacketType
	}

	err := binary.Read(buf, binary.LittleEndian, packet)
	return packet, err
}
