package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/majori/wrc-laptimer/internal/database"
)

/*
  Example:
  {
      "name": "WRC 2025",
      "vehicle_id": 1,
      "vehicle_class_id": 2
  }

  Example Response:
  {
	  "series_id": 1

*/

type CreateSeriesRequest struct {
	Name           string  `json:"name"`
	VehicleID      *uint16 `json:"vehicle_id"`
	VehicleClassID *uint16 `json:"vehicle_class_id"`
}

type CreateSeriesResponse struct {
	SeriesID int `json:"series_id"`
}

func CreateSeriesHandler(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the JSON request body
		var req CreateSeriesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		// Convert nullable fields to sql.NullInt16
		var vehicleID sql.NullInt16
		if req.VehicleID != nil {
			vehicleID = sql.NullInt16{Int16: int16(*req.VehicleID), Valid: true}
		}

		var vehicleClassID sql.NullInt16
		if req.VehicleClassID != nil {
			vehicleClassID = sql.NullInt16{Int16: int16(*req.VehicleClassID), Valid: true}
		}

		// Call the CreateSeries function
		seriesID, err := db.CreateSeries(req.Name, vehicleID, vehicleClassID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create series: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with the created series ID
		response := CreateSeriesResponse{SeriesID: seriesID}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

/*
Example Request:

	{
	    "name": "Event 1",
	    "race_series_id": 1,
	    "location_id": 2,
	    "vehicle_id": 3,
	    "vehicle_class_id": 4
	}

Example Response:

	{
	    "event_id": 1
	}
*/

func parseIDFromPath(r *http.Request, prefix, suffix string) (int, error) {
	// Trim the prefix and suffix from the URL path
	path := strings.TrimPrefix(r.URL.Path, prefix)
	path = strings.TrimSuffix(path, suffix)

	// Convert the extracted string to an integer
	id, err := strconv.Atoi(path)
	if err != nil {
		return 0, fmt.Errorf("invalid ID: %w", err)
	}

	return id, nil
}

type CreateEventRequest struct {
	Name           string  `json:"name"`
	RaceSeriesID   int     `json:"race_series_id"`
	LocationID     *uint16 `json:"location_id"`
	VehicleID      *uint16 `json:"vehicle_id"`
	VehicleClassID *uint16 `json:"vehicle_class_id"`
}

type CreateEventResponse struct {
	EventID int `json:"event_id"`
}

func CreateEventHandler(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the JSON request body
		var req CreateEventRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		// Convert nullable fields to sql.NullInt16
		var locationID sql.NullInt16
		if req.LocationID != nil {
			locationID = sql.NullInt16{Int16: int16(*req.LocationID), Valid: true}
		}

		var vehicleID sql.NullInt16
		if req.VehicleID != nil {
			vehicleID = sql.NullInt16{Int16: int16(*req.VehicleID), Valid: true}
		}

		var vehicleClassID sql.NullInt16
		if req.VehicleClassID != nil {
			vehicleClassID = sql.NullInt16{Int16: int16(*req.VehicleClassID), Valid: true}
		}

		// Call the CreateEvent function
		eventID, err := db.CreateEvent(req.Name, req.RaceSeriesID, locationID, vehicleID, vehicleClassID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create event: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with the created event ID
		response := CreateEventResponse{EventID: eventID}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func StartEventHandler(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure the request method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract the event ID from the URL path
		eventID, err := parseIDFromPath(r, "/api/admin/event/", "/start")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Call the StartEvent function
		err = db.StartEvent(eventID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to start event: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "event started successfully"}`))
	}
}

func EndEventHandler(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure the request method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract the event ID from the URL path
		eventID, err := parseIDFromPath(r, "/api/admin/event/", "/end")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Call the EndEvent function
		err = db.EndEvent(eventID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to end event: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "event ended successfully"}`))
	}
}

func StartSeriesHandler(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure the request method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		seriesID, err := parseIDFromPath(r, "/api/admin/series/", "/start")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Call the StartSeries function
		err = db.StartSeries(seriesID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to start series: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "series started successfully"}`))
	}
}

func EndSeriesHandler(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure the request method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		seriesID, err := parseIDFromPath(r, "/api/admin/series/", "/end")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Call the EndSeries function
		err = db.EndSeries(seriesID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to end series: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "series ended successfully"}`))
	}
}
