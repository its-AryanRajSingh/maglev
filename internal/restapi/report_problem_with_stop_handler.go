package restapi

import (
	"log/slog"
	"net/http"

	"maglev.onebusaway.org/internal/logging"
	"maglev.onebusaway.org/internal/models"
	"maglev.onebusaway.org/internal/utils"
)

func (api *RestAPI) reportProblemWithStopHandler(w http.ResponseWriter, r *http.Request) {
	stopID := utils.ExtractIDFromParams(r)

	// Validate stopID is not empty
	if stopID == "" {
		fieldErrors := map[string][]string{
			"id": {"stop ID is required"},
		}
		api.validationErrorResponse(w, r, fieldErrors)
		return
	}

	// Validate stopID format (must be in agency_stopCode format)
	_, _, err := utils.ExtractAgencyIDAndCodeID(stopID)
	if err != nil {
		fieldErrors := map[string][]string{
			"id": {"invalid stop ID format: " + err.Error()},
		}
		api.validationErrorResponse(w, r, fieldErrors)
		return
	}

	query := r.URL.Query()

	code := query.Get("code")
	userComment := query.Get("userComment")
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
	logging.LogOperation(logger, "problem_report_received_for_stop",
		slog.String("stop_id", stopID),
		slog.String("code", code),
		slog.String("user_comment", userComment),
		slog.String("user_lat", userLat),
		slog.String("user_lon", userLon),
		slog.String("user_location_accuracy", userLocationAccuracy))

	api.sendResponse(w, r, models.NewOKResponse(struct{}{}, api.Clock))
}
