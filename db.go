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

package fball

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Handle struct {
	DB *sql.DB
}

var insertSQL = `
INSERT INTO RequestCache(Endpoint, Params, Timestamp, Response)
				  VALUES(?, ?, ?, ?);`

type response interface {
	When() int64
}

func (h *Handle) Insert(ctx context.Context, endpoint string, data response, params URLQueryStringer) error {
	return transact(ctx, h.DB, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, insertSQL)
		if err != nil {
			return err
		}
		defer stmt.Close()

		blob, err := json.Marshal(data)
		if err != nil {
			return err
		}

		res, err := stmt.ExecContext(ctx, endpoint, params.URLQueryString(), data.When(), blob)
		if err != nil {
			return err
		}
		_, err = res.RowsAffected()
		return err

	})
}

var querySQL = `
SELECT Response from RequestCache
	WHERE
		Endpoint = ?
		AND
			Params = ?
		AND
			Timestamp <= ?
		AND
			Timestamp >= ?
	ORDER BY
		Timestamp DESC
	LIMIT ?
`

type QueryCB func([]byte) error

func (h *Handle) Query(
	ctx context.Context, endpoint string, params URLQueryStringer, max int, r Range, cb QueryCB) error {
	if h == nil || h.DB == nil {
		return nil
	}

	if max < 1 {
		max = 1
	}

	return transact(ctx, h.DB, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, querySQL)
		if err != nil {
			return err
		}
		defer stmt.Close()

		top, bottom := r.UnixNano()
		rows, err := stmt.QueryContext(ctx, endpoint, params.URLQueryString(), top, bottom, max)
		if err != nil {
			return err
		}

		for rows.Next() {
			bytes := []byte{}
			if err := rows.Scan(&bytes); err != nil {
				return err
			}
			if err := cb(bytes); err != nil {
				return err
			}
		}
		return nil
	})
}

type NoParams struct{}

func (np NoParams) URLQueryString() string {
	return ""
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

type URLQueryStringer interface {
	URLQueryString() string
}

func transact(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) (dberr error) {
	if db == nil {
		return fmt.Errorf("no valid database")
	}

	tx, err := db.BeginTx(ctx, nil)
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