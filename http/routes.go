package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/shanebailey05/future_backend_homework/data"
)

func Routes(s *Service) *mux.Router {
	router := mux.NewRouter()
	router.Schemes("https")

	router.HandleFunc("/future/available/appointments/{trainer_id}/{starts_at}/{ends_at}", s.AvailableAppointments).Methods("GET")
	router.HandleFunc("/future/scheduled/appointments/{trainer_id}", s.ScheduledAppointments).Methods("GET")
	router.HandleFunc("/future/save/appointment", s.SaveAppointment).Methods("POST")

	return router
}

func (s *Service) AvailableAppointments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	trainerID, err := strconv.Atoi(vars["trainer_id"])
	if err != nil {
		s.errorResponse(http.StatusBadRequest, fmt.Sprintf("invalid trainer_id: %s\n", err.Error()), w, r)
		return
	}

	startsAt, err := time.Parse("2006-01-02T15:04:05-07:00", vars["starts_at"])
	if err != nil {
		s.errorResponse(http.StatusBadRequest, fmt.Sprintf("invalid starts_at: %s\n", err.Error()), w, r)
		return
	}

	endsAt, err := time.Parse("2006-01-02T15:04:05-07:00", vars["ends_at"])
	if err != nil {
		s.errorResponse(http.StatusBadRequest, fmt.Sprintf("invalid ends_at: %s\n", err.Error()), w, r)
		return
	}

	if startsAt.After(endsAt) || startsAt.Equal(endsAt) {
		s.errorResponse(http.StatusBadRequest, "invalid dates provided", w, r)
		return
	}

	appts, err := s.app.AvailableAppointments(trainerID, startsAt, endsAt)
	if err != nil {
		s.errorResponse(http.StatusInternalServerError, err.Error(), w, r)
		return
	}

	b, err := json.Marshal(appts)
	if err != nil {
		s.errorResponse(http.StatusInternalServerError, err.Error(), w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(b); err != nil {
		s.errorResponse(http.StatusInternalServerError, err.Error(), w, r)
	}
}

func (s *Service) ScheduledAppointments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	trainerID, err := strconv.Atoi(vars["trainer_id"])
	if err != nil {
		s.errorResponse(http.StatusBadRequest, fmt.Sprintf("invalid trainer_id: %s\n", err.Error()), w, r)
		return
	}

	appts, err := s.app.ScheduledAppointments(trainerID)
	if err != nil {
		s.errorResponse(http.StatusInternalServerError, err.Error(), w, r)
		return
	}

	b, err := json.Marshal(appts)
	if err != nil {
		s.errorResponse(http.StatusInternalServerError, err.Error(), w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(b); err != nil {
		s.errorResponse(http.StatusInternalServerError, err.Error(), w, r)
	}
}

func (s *Service) SaveAppointment(w http.ResponseWriter, r *http.Request) {
	appt := new(data.Appointment)
	if err := json.NewDecoder(r.Body).Decode(&appt); err != nil {
		s.errorResponse(http.StatusBadRequest, "Error decoding request data", w, r)
		return
	}

	if err := s.app.SaveAppointment(appt); err != nil {
		code := http.StatusInternalServerError
		if strings.Contains(err.Error(), "Appointment is invalid") {
			code = http.StatusBadRequest
		}
		s.errorResponse(code, err.Error(), w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
