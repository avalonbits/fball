package db

import (
	"context"
	"database/sql"
	"time"

	"git.cana.pw/avalonbits/fball"
)

type Querier struct {
	DB *sql.DB
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

func (q *Querier) Timezone(ctx context.Context, max int, r Range) ([]fball.TimezoneResponse, error) {
	if max < 1 {
		max = 1
	}

	tzResp := []fball.TimezoneResponse{}
	err := transact(ctx, q.DB, func(tx *sql.Tx) error {
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
