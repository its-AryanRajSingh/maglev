package restapi

import (
	"log/slog"
	"net/http"

	"maglev.onebusaway.org/internal/logging"
	"maglev.onebusaway.org/internal/models"
	"maglev.onebusaway.org/internal/utils"
)

func (api *RestAPI) reportProblemWithTripHandler(w http.ResponseWriter, r *http.Request) {

	tripID := utils.ExtractIDFromParams(r)

	// Validate tripID is not empty
	if tripID == "" {
		fieldErrors := map[string][]string{
			"id": {"trip ID is required"},
		}
		api.validationErrorResponse(w, r, fieldErrors)
		return
	}

	// Validate tripID format (must be in agency_tripCode format)
	_, _, err := utils.ExtractAgencyIDAndCodeID(tripID)
	if err != nil {
		fieldErrors := map[string][]string{
			"id": {"invalid trip ID format: " + err.Error()},
		}
		api.validationErrorResponse(w, r, fieldErrors)
		return
	}

	query := r.URL.Query()

	serviceDate := query.Get("serviceDate")
	vehicleID := query.Get("vehicleId")
	stopID := query.Get("stopId")
	code := query.Get("code")
	userComment := query.Get("userComment")
	userOnVehicle := query.Get("userOnVehicle")
	userVehicleNumber := query.Get("userVehicleNumber")
	userLat := query.Get("userLat")
	userLon := query.Get("userLon")
	userLocationAccuracy := query.Get("userLocationAccuracy")

	// Validate required code parameter
	if code == "" {
		fieldErrors := map[string][]string{
			"code": {"problem code is required"},
		}
		api.validationErrorResponse(w, r, fieldErrors)
		return
	}

	// Log the problem report
	logger := logging.FromContext(r.Context()).With(slog.String("component", "problem_reporting"))
	logging.LogOperation(logger, "problem_report_received_for_trip",
		slog.String("trip_id", tripID),
		slog.String("code", code),
		slog.String("service_date", serviceDate),
		slog.String("vehicle_id", vehicleID),
		slog.String("stop_id", stopID),
		slog.String("user_comment", userComment),
		slog.String("user_on_vehicle", userOnVehicle),
		slog.String("user_vehicle_number", userVehicleNumber),
		slog.String("user_lat", userLat),
		slog.String("user_lon", userLon),
		slog.String("user_location_accuracy", userLocationAccuracy))

	api.sendResponse(w, r, models.NewOKResponse(struct{}{}, api.Clock))
}
