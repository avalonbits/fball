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
	"encoding/json"

	"git.cana.pw/avalonbits/fball"
	"git.cana.pw/avalonbits/fball/client"
)

type Querier struct {
	DB *sql.DB
}

func (q *Querier) Country(
	ctx context.Context, params client.CountryParams, max int, r Range) ([]fball.CountryResponse, error) {
	resp := []fball.CountryResponse{}
	if err := q.query(ctx, fball.EP_Countries, params, max, r, func(data []byte) error {
		cr := fball.CountryResponse{}
		if err := json.Unmarshal(data, &cr); err != nil {
			return err
		}
		resp = append(resp, cr)
		return nil
	}); err != nil {
		return nil, err
	}
	return resp, nil
}

type QueryCB func([]byte) error

func (q *Querier) query(
	ctx context.Context, endpoint string, params urlQueryStringer, max int, r Range, cb QueryCB) error {
	if q == nil || q.DB == nil {
		return nil
	}

	if max < 1 {
		max = 1
	}

	return transact(ctx, q.DB, func(tx *sql.Tx) error {
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
