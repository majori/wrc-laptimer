package http

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
      "vehicle_class_id": 2
  }

  Example Response:
  {
	  "series_id": 1

*/

type CreateSeriesRequest struct {
	Name           string  `json:"name"`
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
		var vehicleClassID sql.NullInt16
		if req.VehicleClassID != nil {
			vehicleClassID = sql.NullInt16{Int16: int16(*req.VehicleClassID), Valid: true}
		}

		// Call the CreateSeries function
		seriesID, err := db.CreateSeries(req.Name, vehicleClassID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create series: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with the created series ID
		response := CreateSeriesResponse{SeriesID: seriesID}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

/*
Example Request:

	{
	    "name": "Event 1",
	    "race_series_id": 1,
	    "location_id": 2,
		"route_id": 3,
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
	RaceSeriesID   *uint16 `json:"race_series_id"`
	LocationID     *uint16 `json:"location_id"`
	RouteID        *uint16 `json:"route_id"`
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

		var seriesID sql.NullInt32
		if req.RaceSeriesID != nil {
			seriesID = sql.NullInt32{Int32: int32(*req.RaceSeriesID), Valid: true}
		}

		// Convert nullable fields to sql.NullInt16
		var locationID sql.NullInt16
		if req.LocationID != nil {
			locationID = sql.NullInt16{Int16: int16(*req.LocationID), Valid: true}
		}

		var routeID sql.NullInt16
		if req.RouteID != nil {
			routeID = sql.NullInt16{Int16: int16(*req.RouteID), Valid: true}
		}

		var vehicleClassID sql.NullInt16
		if req.VehicleClassID != nil {
			vehicleClassID = sql.NullInt16{Int16: int16(*req.VehicleClassID), Valid: true}
		}

		// Call the CreateEvent function
		eventID, err := db.CreateEvent(req.Name, seriesID, locationID, routeID, vehicleClassID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create event: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with the created event ID
		response := CreateEventResponse{EventID: eventID}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
			return
		}
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
		eventID, err := parseIDFromPath(r, "/api/admin/events/", "/start")
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
		if _, err := w.Write([]byte(`{"status": "event started successfully"}`)); err != nil {
			http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
			return
		}
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
		eventID, err := parseIDFromPath(r, "/api/admin/events/", "/end")
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
		if _, err := w.Write([]byte(`{"status": "event ended successfully"}`)); err != nil {
			http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
			return
		}
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
		if _, err := w.Write([]byte(`{"status": "series started successfully"}`)); err != nil {
			http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
			return
		}
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
		if _, err := w.Write([]byte(`{"status": "series ended successfully"}`)); err != nil {
			http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
