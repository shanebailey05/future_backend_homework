package app

import (
	"time"

	"github.com/shanebailey05/future_backend_homework/data"
)

func (s *Service) AvailableAppointments(trainerID int, startsAt, endsAt time.Time) ([]*data.Appointment, error) {
	if s.appointments == nil {
		s.appointments = new(data.Appointments)
	}

	startsAt = startsAt.Round(time.Minute * 30)
	endsAt = endsAt.Round(time.Minute * 30)

	appts, err := s.ScheduledAppointments(trainerID)
	if err != nil {
		return nil, err
	}

	ret := []*data.Appointment{}
	for {
		if startsAt.Equal(endsAt) {
			break
		}

		a := &data.Appointment{
			TrainerID: trainerID,
			StartedAt: startsAt,
			EndedAt:   startsAt.Add(time.Minute * 30),
		}

		for _, sa := range *appts {
			if sa.StartedAt.Equal(a.StartedAt) {
				a.StartedAt = sa.EndedAt.Add(time.Minute * 30)
				a.EndedAt = a.StartedAt.Add(time.Minute * 30)
				startsAt = a.EndedAt
			}
		}

		ret = append(ret, a)

		startsAt = startsAt.Add(time.Minute * 30)
	}

	return ret, nil
}

func (s *Service) ScheduledAppointments(trainerID int) (*data.Appointments, error) {
	if s.appointments == nil {
		s.appointments = new(data.Appointments)
	}

	ret := new(data.Appointments)
	for _, a := range *s.appointments {
		if a.TrainerID == trainerID {
			ret.Append(a)
		}
	}

	return ret, nil
}

func (s *Service) SaveAppointment(appt *data.Appointment) error {
	err := s.appointments.Save(appt)
	if err != nil {
		return err
	}

	s.appointments, err = data.AllAppointments()
	if err != nil {
		return err
	}

	return nil
}
