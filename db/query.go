/*
 * Copyright (C) 2021  Igor Cananea <icc@avalonbits.com>
 * Author: Igor Cananea <icc@avalonbits.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
	if q == nil || q.DB == nil {
		return nil, nil
	}

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