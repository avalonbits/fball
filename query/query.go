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

func (r Range) UnixNano() (top, bottom int64) {
	if r.Latest.IsZero() {
		top = time.Now().UTC().UnixNano()
	} else {
		top = r.Latest.UTC().UnixNano()
	}
	if !r.Earliest.IsZero() {
		bottom = r.Earliest.UTC().UnixNano()
	}
	return
}

func (r Range) IsZero() bool {
	return r.Latest.IsZero() && r.Earliest.IsZero()
}

var noParamQuery = `
SELECT Response from RequestCache
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

const (
	timezoneEP = "/timezone"
)

func (h *Handler) Timezone(ctx context.Context, max int, r Range) ([]fball.TimezoneResponse, error) {
	if max < 1 {
		max = 1
	}

	tzResp := []fball.TimezoneResponse{}
	err := h.transact(ctx, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, noParamQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		top, bottom := r.UnixNano()
		rows, err := stmt.QueryContext(ctx, fball.EP_Timezone, top, bottom, max)
		if err != nil {
			return err
		}

		for rows.Next() {
			tr := fball.TimezoneResponse{}
			if err := rows.Scan(&tr); err != nil {
				return err
			}
			tzResp = append(tzResp, tr)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return tzResp, nil
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
