package data

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Appointment struct {
	ID        int       `json:"id,omitempty"`
	TrainerID int       `json:"trainer_id,omitempty"`
	UserID    int       `json:"user_id,omitempty"`
	StartedAt time.Time `json:"started_at,omitempty"`
	EndedAt   time.Time `json:"ended_at,omitempty"`
}

type Appointments []*Appointment

const file string = "../data/appointments.json" // TODO: change this path to data/appointments.json

func (appts *Appointments) Append(a *Appointment) {
	*appts = append(*appts, a)
}

func (appts *Appointments) Remove(a *Appointment) {
	found, index := appts.findIndexByID(a.ID)

	if !found {
		return
	}

	newAppts := new(Appointments)
	for i, a := range *appts {
		if i == index {
			continue
		}
		newAppts.Append(a)
	}

	*appts = *newAppts
}

func (appts *Appointments) findIndexByID(id int) (bool, int) {
	for i, a := range *appts {
		if a.ID == id {
			return true, i
		}
	}
	return false, 0
}

func (appts *Appointments) Sort() {
	// TODO: sort the slice
}

func (appts *Appointments) Validate(a *Appointment) error {
	var errs []string

	if a.TrainerID == 0 {
		errs = append(errs, "trainer was not provided")
	}

	if a.UserID == 0 {
		errs = append(errs, "user was not provided")
	}

	if a.StartedAt.IsZero() {
		errs = append(errs, "starting date was not provided")
	}

	if a.EndedAt.IsZero() {
		errs = append(errs, "ending date was not provided")
	}

	if (a.EndedAt.Sub(a.StartedAt)/time.Minute)%30 != 0 {
		errs = append(errs, "appointment is not in 30 minute increments")
	}

	if a.StartedAt.After(a.EndedAt) {
		errs = append(errs, "starting time is after the ending time")
	}

	if a.EndedAt.Before(a.StartedAt) {
		errs = append(errs, "ending time is before the starting time")
	}

	if a.StartedAt.Minute() != 0 || a.StartedAt.Minute() != 30 {
		errs = append(errs, "starting time must have either :00 or :30 for its minutes")
	}

	if a.EndedAt.Minute() != 0 || a.EndedAt.Minute() != 30 {
		errs = append(errs, "ending time must have :00 or :30 for its minutes")
	}

	if a.StartedAt.Hour() < 8 {
		errs = append(errs, "starting time is before business hours")
	}

	if a.EndedAt.Hour() > 17 {
		errs = append(errs, "ending time is after business hours")
	}

	if a.StartedAt.Day() == int(time.Saturday) {
		errs = append(errs, "appointment cannot be on Saturday")
	}

	if a.StartedAt.Day() == int(time.Sunday) {
		errs = append(errs, "appointment cannot be on Sunday")
	}

	if !appts.isAppointmentAvailable(a) {
		errs = append(errs, "trainer is already booked within the requested timeframe")
	}

	if len(errs) > 0 {
		return fmt.Errorf("Appointment is invalid for the following reasons: %s", strings.Join(errs, ", "))
	}

	return nil
}

func (appts *Appointments) isAppointmentAvailable(a *Appointment) bool {
	for _, appt := range *appts {
		if a.TrainerID == appt.TrainerID && a.StartedAt == appt.StartedAt && a.EndedAt == appt.EndedAt {
			return false
		}
	}

	return true
}

func AllAppointments() (*Appointments, error) {
	file, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	ret := new(Appointments)
	if err = json.Unmarshal(file, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (appts *Appointments) Save(a *Appointment) error {
	if err := appts.Validate(a); err != nil {
		return err
	}

	if a.ID == 0 {
		return appts.insert(a)
	}

	return appts.update(a)
}

func (appts *Appointments) insert(a *Appointment) error {
	a.ID = appts.maxID() + 1
	appts.Append(a)
	appts.Sort()

	/* TODO: uncomment this
	var str string
	for _, a := range *appts {
		str += a.String()
		str += ","
	}

	str = strings.TrimSuffix(str, ",")
	if err := os.WriteFile(file, []byte(str), 0644); err != nil {
		return fmt.Errorf("error inserting appointment: %s", err.Error())
	}
	*/

	return nil
}

func (appts *Appointments) update(a *Appointment) error {
	appts.Remove(a)
	appts.Append(a)
	appts.Sort()

	/* TODO: uncomment
	var str string
	for _, a := range *appts {
		str += a.String()
		str += ","
	}

	str = strings.TrimSuffix(str, ",")
	if err := os.WriteFile(file, []byte(str), 0644); err != nil {
		return fmt.Errorf("error inserting appointment: %s", err.Error())
	}
	*/

	return nil
}

func (appts *Appointments) maxID() int {
	var ret int
	for _, a := range *appts {
		if a.ID > ret {
			ret = a.ID
		}
	}

	return ret
}

func (a *Appointment) String() string {
	return fmt.Sprintf(
		"{\n\t\"id\": %d\n\t\"trainer_id\": %d\n\t\"user_id\": %d\n\t\"started_at\": \"%s\"\n\t\"ended_at\": \"%s\"\n}",
		a.ID,
		a.TrainerID,
		a.UserID,
		a.StartedAt.String(),
		a.EndedAt.String(),
	)
}
