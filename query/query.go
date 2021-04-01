package query

import (
	"context"
	"database/sql"
	"time"

	"git.cana.pw/avalonbits/fball"
)

type Handler struct {
	db *sql.DB
}

func New(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

type Range struct {
	Latest   time.Time
	Earliest time.Time
}

var noParamQuery = `
SELECT Response, Timestamp from RequestCache
	WHERE
		Endpoint = ?
	ORDER BY
		Timestamp DESC
	LIMIT ?;
`

var rangeNoParamQuery = `
SELECT Response, Timestamp from RequestCache
	WHERE
		Endpoint = ?
		AND
			Timestamp <= ?
		AND
			Timestamp >= ?
	ORDER BY
		Timestamp DESC
	LIMIT ?
`

func (h *Handler) Timezone(ctx context.Context, max int, r *Range) ([]fball.TimezoneResponse, error) {
	if max < 1 {
		max = 1
	}

	tr := []fball.TimezoneResponse{}
	err := h.transact(ctx, func(tx *sql.Tx) error {
		if r == nil || (r.Latest.IsZero() && r.Earliest.IsZero()) {
			stmt, err := tx.PrepareContext(ctx, noParamQuery)
			if err != nil {
				return err
			}
			res, err := stmt.Exec(ctx, max)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return tr, nil
}

func (h *Handler) transact(ctx context.Context, fn func(*sql.Tx) error) (dberr error) {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if dberr != nil {
			tx.Rollback()
		} else {
			dberr = tx.Commit()
		}
	}()

	dberr = fn(tx)
	return
}
