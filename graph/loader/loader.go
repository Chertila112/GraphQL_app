package loader

import (
	"context"
	"net/http"

	"github.com/vikstrous/dataloadgen"

	"github.com/Chertila112/server/internal/db/sqlc"
)

type ctxKey struct{}

type Loaders struct {
	UserByID               *dataloadgen.Loader[string, *sqlc.User]
	ExerciseByID           *dataloadgen.Loader[string, *sqlc.Exercise]
	DaysByScheduleID       *dataloadgen.Loader[string, []*sqlc.ScheduleDay]
	ExercisesByDayID       *dataloadgen.Loader[string, []*sqlc.ScheduleExercise]
	ActiveScheduleByUserID *dataloadgen.Loader[string, *sqlc.Schedule]
}

func New(q *sqlc.Queries) *Loaders {
	return &Loaders{
		UserByID:               newUserLoader(q),
		ExerciseByID:           newExerciseLoader(q),
		DaysByScheduleID:       newDaysLoader(q),
		ExercisesByDayID:       newExercisesByDayLoader(q),
		ActiveScheduleByUserID: newActiveScheduleLoader(q),
	}
}

func Middleware(q *sqlc.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ctxKey{}, New(q))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func For(ctx context.Context) *Loaders {
	l, ok := ctx.Value(ctxKey{}).(*Loaders)
	if !ok {
		panic("loader: not in context — middleware missing on /query handler")
	}
	return l
}

func newUserLoader(q *sqlc.Queries) *dataloadgen.Loader[string, *sqlc.User] {
	return dataloadgen.NewMappedLoader[string, *sqlc.User](
		func(ctx context.Context, ids []string) (map[string]*sqlc.User, error) {
			rows, err := q.GetUsersByIDs(ctx, ids)
			if err != nil {
				return nil, err
			}

			out := make(map[string]*sqlc.User, len(rows))
			for i := range rows {
				out[rows[i].ID] = &rows[i]
			}
			return out, nil
		},
	)
}

func newExerciseLoader(q *sqlc.Queries) *dataloadgen.Loader[string, *sqlc.Exercise] {
	return dataloadgen.NewMappedLoader[string, *sqlc.Exercise](
		func(ctx context.Context, ids []string) (map[string]*sqlc.Exercise, error) {
			rows, err := q.GetExercisesByIDs(ctx, ids)
			if err != nil {
				return nil, err
			}

			out := make(map[string]*sqlc.Exercise, len(rows))
			for i := range rows {
				out[rows[i].ID] = &rows[i]
			}
			return out, nil
		},
	)
}

func newDaysLoader(q *sqlc.Queries) *dataloadgen.Loader[string, []*sqlc.ScheduleDay] {
	return dataloadgen.NewMappedLoader[string, []*sqlc.ScheduleDay](
		func(ctx context.Context, scheduleIDs []string) (map[string][]*sqlc.ScheduleDay, error) {
			rows, err := q.GetScheduleDaysByScheduleIDs(ctx, scheduleIDs)
			if err != nil {
				return nil, err
			}

			out := make(map[string][]*sqlc.ScheduleDay, len(scheduleIDs))
			for _, id := range scheduleIDs {
				out[id] = []*sqlc.ScheduleDay{}
			}

			for i := range rows {
				sid := rows[i].ScheduleID
				out[sid] = append(out[sid], &rows[i])
			}
			return out, nil
		},
	)
}

func newExercisesByDayLoader(q *sqlc.Queries) *dataloadgen.Loader[string, []*sqlc.ScheduleExercise] {
	return dataloadgen.NewMappedLoader[string, []*sqlc.ScheduleExercise](
		func(ctx context.Context, dayIDs []string) (map[string][]*sqlc.ScheduleExercise, error) {
			rows, err := q.GetScheduleExercisesByDayIDs(ctx, dayIDs)
			if err != nil {
				return nil, err
			}

			out := make(map[string][]*sqlc.ScheduleExercise, len(dayIDs))
			for _, id := range dayIDs {
				out[id] = []*sqlc.ScheduleExercise{}
			}

			for i := range rows {
				did := rows[i].ScheduleDayID
				out[did] = append(out[did], &rows[i])
			}
			return out, nil
		},
	)
}

func newActiveScheduleLoader(q *sqlc.Queries) *dataloadgen.Loader[string, *sqlc.Schedule] {
	return dataloadgen.NewMappedLoader[string, *sqlc.Schedule](
		func(ctx context.Context, userIDs []string) (map[string]*sqlc.Schedule, error) {
			rows, err := q.GetSchedulesByUserIDs(ctx, userIDs)
			if err != nil {
				return nil, err
			}

			out := make(map[string]*sqlc.Schedule, len(rows))
			for i := range rows {
				out[rows[i].UserID] = &rows[i]
			}
			return out, nil
		},
	)
}
