package graph

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"

	"github.com/Chertila112/server/auth"
	"github.com/Chertila112/server/graph/model"
	"github.com/Chertila112/server/internal/db/sqlc"
)

// MUTATION: USER
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u, err := r.Queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:          uuid.NewString(),
		Username:    input.Username,
		Email:       input.Email,
		Password:    string(hash),
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		DateOfBirth: toPgDate(input.DateOfBirth),
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return toGraphQLUser(&u), nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.UpdateUser) (*model.User, error) {
	params := sqlc.UpdateUserParams{
		ID:        id,
		Username:  input.Username,
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	if input.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("hash password: %w", err)
		}
		s := string(hash)
		params.Password = &s
	}
	u, err := r.Queries.UpdateUser(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("update user: %w", err)
	}
	return toGraphQLUser(&u), nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	if err := r.Queries.DeleteUser(ctx, id); err != nil {
		return false, fmt.Errorf("delete user: %w", err)
	}
	return true, nil
}

// MUTATION: EXERCISE
func (r *mutationResolver) CreateExercise(ctx context.Context, input model.NewExercise) (*model.Exercise, error) {
	var eq *sqlc.Equipment
	if input.Equipment != nil {
		v := sqlc.Equipment(*input.Equipment)
		eq = &v
	}

	e, err := r.Queries.CreateExercise(ctx, sqlc.CreateExerciseParams{
		ID:          uuid.NewString(),
		Name:        input.Name,
		Description: input.Description,
		MuscleGroup: sqlc.MuscleGroup(input.MuscleGroup),
		Equipment:   eq,
	})
	if err != nil {
		return nil, fmt.Errorf("create exercise: %w", err)
	}
	return toGraphQLExercise(&e), nil
}

func (r *mutationResolver) UpdateExercise(ctx context.Context, id string, input model.UpdateExercise) (*model.Exercise, error) {
	params := sqlc.UpdateExerciseParams{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
	}

	if input.MuscleGroup != nil {
		mg := sqlc.MuscleGroup(*input.MuscleGroup)
		params.MuscleGroup = &mg
	}
	if input.Equipment != nil {
		eq := sqlc.Equipment(*input.Equipment)
		params.Equipment = &eq
	}

	e, err := r.Queries.UpdateExercise(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("exercise not found")
		}
		return nil, fmt.Errorf("update exercise: %w", err)
	}
	return toGraphQLExercise(&e), nil
}

func (r *mutationResolver) DeleteExercise(ctx context.Context, id string) (bool, error) {
	if err := r.Queries.SoftDeleteExercise(ctx, id); err != nil {
		return false, fmt.Errorf("delete exercise: %w", err)
	}
	return true, nil
}

// MUTATION: SCHEDULE
func (r *mutationResolver) CreateSchedule(ctx context.Context, input model.NewSchedule) (*model.Schedule, error) {
	if _, err := r.Queries.GetUser(ctx, input.UserID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	s, err := r.Queries.CreateSchedule(ctx, sqlc.CreateScheduleParams{
		ID:     uuid.NewString(),
		UserID: input.UserID,
		Name:   input.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("create schedule: %w", err)
	}
	return toGraphQLSchedule(&s), nil
}

func (r *mutationResolver) UpdateSchedule(ctx context.Context, id string, input model.UpdateSchedule) (*model.Schedule, error) {
	if input.IsActive != nil && *input.IsActive {
		current, err := r.Queries.GetSchedule(ctx, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("schedule not found")
			}
			return nil, fmt.Errorf("get schedule: %w", err)
		}
		if err := r.Queries.DeactivateOtherSchedules(ctx, sqlc.DeactivateOtherSchedulesParams{
			UserID: current.UserID,
			ID:     id,
		}); err != nil {
			return nil, fmt.Errorf("deactivate other schedules: %w", err)
		}
	}

	s, err := r.Queries.UpdateSchedule(ctx, sqlc.UpdateScheduleParams{
		ID:       id,
		Name:     input.Name,
		IsActive: input.IsActive,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("schedule not found")
		}
		return nil, fmt.Errorf("update schedule: %w", err)
	}
	return toGraphQLSchedule(&s), nil
}

func (r *mutationResolver) DeleteSchedule(ctx context.Context, id string) (bool, error) {
	if err := r.Queries.DeleteSchedule(ctx, id); err != nil {
		return false, fmt.Errorf("delete schedule: %w", err)
	}
	return true, nil
}

// MUTATION: SCHEDULE DAY
func (r *mutationResolver) UpdateScheduleDay(ctx context.Context, id string, input model.UpdateScheduleDay) (*model.ScheduleDay, error) {
	params := sqlc.UpdateScheduleDayParams{
		ID:    id,
		Notes: input.Notes,
	}

	if input.DayOfWeek != nil {
		v := int32(*input.DayOfWeek)
		params.DayOfWeek = &v
	}
	if input.IsRestDay != nil {
		params.IsRestDay = input.IsRestDay
	}

	d, err := r.Queries.UpdateScheduleDay(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("schedule day not found")
		}
		return nil, fmt.Errorf("update schedule day: %w", err)
	}
	return toGraphQLScheduleDay(&d), nil
}

// MUTATION: SCHEDULE EXERCISE
func (r *mutationResolver) AddExerciseToDay(ctx context.Context, input model.NewScheduleExercise) (*model.ScheduleExercise, error) {
	if _, err := r.Queries.GetScheduleDay(ctx, input.ScheduleDayID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("schedule day not found")
		}
		return nil, fmt.Errorf("get schedule day: %w", err)
	}
	if _, err := r.Queries.GetExercise(ctx, input.ExerciseID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("exercise not found")
		}
		return nil, fmt.Errorf("get exercise: %w", err)
	}

	se, err := r.Queries.AddExerciseToDay(ctx, sqlc.AddExerciseToDayParams{
		ID:            uuid.NewString(),
		ScheduleDayID: input.ScheduleDayID,
		ExerciseID:    input.ExerciseID,
		Position:      int32(input.Position),
		Sets:          input.Sets,
		Reps:          input.Reps,
		Weight:        input.Weight,
		Rest:          input.Rest,
		Duration:      input.Duration,
		Notes:         input.Notes,
	})
	if err != nil {
		return nil, fmt.Errorf("add exercise to day: %w", err)
	}
	return toGraphQLScheduleExercise(&se), nil
}

func (r *mutationResolver) UpdateScheduleExercise(ctx context.Context, id string, input model.UpdateScheduleExercise) (*model.ScheduleExercise, error) {
	params := sqlc.UpdateScheduleExerciseParams{
		ID:       id,
		Sets:     input.Sets,
		Reps:     input.Reps,
		Rest:     input.Rest,
		Duration: input.Duration,
		Notes:    input.Notes,
	}

	if input.Position != nil {
		v := int32(*input.Position)
		params.Position = &v
	}
	if input.Weight != nil {
		params.Weight = input.Weight
	}

	se, err := r.Queries.UpdateScheduleExercise(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("schedule exercise not found")
		}
		return nil, fmt.Errorf("update schedule exercise: %w", err)
	}
	return toGraphQLScheduleExercise(&se), nil
}

func (r *mutationResolver) RemoveExerciseFromDay(ctx context.Context, id string) (bool, error) {
	if err := r.Queries.RemoveExerciseFromDay(ctx, id); err != nil {
		return false, fmt.Errorf("remove exercise from day: %w", err)
	}
	return true, nil
}

// QUERY: USER
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	u, err := r.Queries.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	return toGraphQLUser(&u), nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	u, err := r.Queries.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	return toGraphQLUser(&u), nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	rows, err := r.Queries.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return toGraphQLUsers(rows), nil
}

// QUERY: EXERCISE
func (r *queryResolver) Exercise(ctx context.Context, id string) (*model.Exercise, error) {
	e, err := r.Queries.GetExercise(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get exercise: %w", err)
	}
	return toGraphQLExercise(&e), nil
}

func (r *queryResolver) Exercises(ctx context.Context, muscleGroup *model.MuscleGroup, equipment *model.Equipment) ([]*model.Exercise, error) {
	params := sqlc.ListExercisesParams{}
	if muscleGroup != nil {
		mg := sqlc.MuscleGroup(*muscleGroup)
		params.MuscleGroup = &mg
	}
	if equipment != nil {
		eq := sqlc.Equipment(*equipment)
		params.Equipment = &eq
	}

	rows, err := r.Queries.ListExercises(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("list exercises: %w", err)
	}
	return toGraphQLExercises(rows), nil
}

// QUERY: SCHEDULE
func (r *queryResolver) Schedule(ctx context.Context, id string) (*model.Schedule, error) {
	s, err := r.Queries.GetSchedule(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get schedule: %w", err)
	}
	return toGraphQLSchedule(&s), nil
}

func (r *queryResolver) ScheduleByUser(ctx context.Context, userID string) (*model.Schedule, error) {
	s, err := r.Queries.GetActiveScheduleByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get schedule by user: %w", err)
	}
	return toGraphQLSchedule(&s), nil
}

func (r *queryResolver) Schedules(ctx context.Context) ([]*model.Schedule, error) {
	rows, err := r.Queries.ListSchedules(ctx)
	if err != nil {
		return nil, fmt.Errorf("list schedules: %w", err)
	}
	return toGraphQLSchedules(rows), nil
}

// QUERY: SCHEDULE DAY
func (r *queryResolver) ScheduleDay(ctx context.Context, id string) (*model.ScheduleDay, error) {
	d, err := r.Queries.GetScheduleDay(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get schedule day: %w", err)
	}
	return toGraphQLScheduleDay(&d), nil
}

func (r *queryResolver) ScheduleDays(ctx context.Context, scheduleID string) ([]*model.ScheduleDay, error) {
	rows, err := r.Queries.GetScheduleDaysByScheduleIDs(ctx, []string{scheduleID})
	if err != nil {
		return nil, fmt.Errorf("list schedule days: %w", err)
	}
	return toGraphQLScheduleDays(rows), nil
}

// SUBSCRIPTION
func (r *subscriptionResolver) ScheduleUpdated(ctx context.Context, userID string) (<-chan *model.Schedule, error) {
	ch := make(chan *model.Schedule)
	go func() {
		<-ctx.Done()
		close(ch)
	}()
	return ch, nil
}

func toPgDate(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: *t, Valid: true}
}

// ============================================================
// HELPERS
// ============================================================

func toIntPtr(v *int32) *int {
	if v == nil {
		return nil
	}
	x := int(*v)
	return &x
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type (
	mutationResolver     struct{ *Resolver }
	queryResolver        struct{ *Resolver }
	subscriptionResolver struct{ *Resolver }
)
