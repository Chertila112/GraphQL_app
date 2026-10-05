package graph

import (
	"time"

	"github.com/Chertila112/server/graph/model"
	"github.com/Chertila112/server/internal/db/sqlc"
)

// DB → GraphQL
func toGraphQLUser(u *sqlc.User) *model.User {
	if u == nil {
		return nil
	}

	// Check if pgtype.Date is
	var dob *time.Time
	if u.DateOfBirth.Valid {
		t := u.DateOfBirth.Time
		dob = &t
	}

	return &model.User{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		DateOfBirth: dob,
		CreatedAt:   u.CreatedAt,
	}
}

func toGraphQLExercise(e *sqlc.Exercise) *model.Exercise {
	if e == nil {
		return nil
	}
	var equipment *model.Equipment
	if e.Equipment != nil {
		v := model.Equipment(*e.Equipment)
		equipment = &v
	}
	return &model.Exercise{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		MuscleGroup: model.MuscleGroup(e.MuscleGroup),
		Equipment:   equipment,
		CreatedAt:   e.CreatedAt,
	}
}

func toGraphQLSchedule(s *sqlc.Schedule) *model.Schedule {
	if s == nil {
		return nil
	}
	return &model.Schedule{
		ID:        s.ID,
		UserID:    s.UserID,
		Name:      s.Name,
		IsActive:  s.IsActive,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func toGraphQLScheduleDay(d *sqlc.ScheduleDay) *model.ScheduleDay {
	if d == nil {
		return nil
	}
	return &model.ScheduleDay{
		ID:         d.ID,
		ScheduleID: d.ScheduleID,
		DayOfWeek:  d.DayOfWeek,
		IsRestDay:  d.IsRestDay,
		Notes:      d.Notes,
	}
}

func toGraphQLScheduleExercise(e *sqlc.ScheduleExercise) *model.ScheduleExercise {
	if e == nil {
		return nil
	}
	return &model.ScheduleExercise{
		ID:            e.ID,
		ScheduleDayID: e.ScheduleDayID,
		ExerciseID:    e.ExerciseID,
		Position:      e.Position,
		Sets:          e.Sets,
		Reps:          e.Reps,
		Weight:        e.Weight,
		Rest:          e.Rest,
		Duration:      e.Duration,
		Notes:         e.Notes,
	}
}

func toGraphQLUsers(users []sqlc.User) []*model.User {
	result := make([]*model.User, 0, len(users))
	for i := range users {
		result = append(result, toGraphQLUser(&users[i]))
	}
	return result
}

func toGraphQLExercises(exercises []sqlc.Exercise) []*model.Exercise {
	result := make([]*model.Exercise, 0, len(exercises))
	for i := range exercises {
		result = append(result, toGraphQLExercise(&exercises[i]))
	}
	return result
}

func toGraphQLSchedules(schedules []sqlc.Schedule) []*model.Schedule {
	result := make([]*model.Schedule, 0, len(schedules))
	for i := range schedules {
		result = append(result, toGraphQLSchedule(&schedules[i]))
	}
	return result
}

func toGraphQLScheduleDays(days []sqlc.ScheduleDay) []*model.ScheduleDay {
	result := make([]*model.ScheduleDay, 0, len(days))
	for i := range days {
		result = append(result, toGraphQLScheduleDay(&days[i]))
	}
	return result
}

func toGraphQLScheduleExercises(items []sqlc.ScheduleExercise) []*model.ScheduleExercise {
	result := make([]*model.ScheduleExercise, 0, len(items))
	for i := range items {
		result = append(result, toGraphQLScheduleExercise(&items[i]))
	}
	return result
}
