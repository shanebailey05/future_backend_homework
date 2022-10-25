package app

import "github.com/shanebailey05/future_backend_homework/data"

type Service struct {
	appointments *data.Appointments
}

func New() (*Service, error) {
	ret := new(Service)
	appts, err := data.AllAppointments()
	if err != nil {
		return nil, err
	}
	ret.appointments = appts
	return ret, nil
}
