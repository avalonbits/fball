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
)

type cache struct {
	DB *sql.DB
}

var insertSQL = `
INSERT INTO RequestCache(Endpoint, Params, Timestamp, Response)
				  VALUES(?, ?, ?, ?);`

func (c *cache) Insert(ctx context.Context, endpoint string, data Response, params urlQueryStringer) error {
	return transact(ctx, c.DB, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, insertSQL)
		if err != nil {
			return err
		}
		defer stmt.Close()

		blob, err := json.Marshal(data)
		if err != nil {
			return err
		}

		res, err := stmt.ExecContext(ctx, endpoint, params.urlQueryString(), data.When(), blob)
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

type queryCB func([]byte) error

func (c *cache) Query(
	ctx context.Context, endpoint string, params urlQueryStringer, max int, r tRange, cb queryCB) error {
	if c == nil || c.DB == nil {
		return nil
	}

	if max < 1 {
		max = 1
	}

	return transact(ctx, c.DB, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, querySQL)
		if err != nil {
			return err
		}
		defer stmt.Close()

		top, bottom := r.UnixNano()
		rows, err := stmt.QueryContext(ctx, endpoint, params.urlQueryString(), top, bottom, max)
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
