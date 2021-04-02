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

type Inserter struct {
	DB *sql.DB
}

var insertSQL = `
INSERT INTO RequestCache(Endpoint, Params, Timestamp, Response)
				  VALUES(?, ?, ?, ?);`

func (i *Inserter) Timezone(ctx context.Context, tr fball.TimezoneResponse) error {
	return i.insert(ctx, tr, nil)
}

func (i *Inserter) Country(ctx context.Context, cr fball.CountryResponse, cp client.CountryParams) error {
	return i.insert(ctx, cr, cp)
}

type response interface {
	When() int64
}

type urlQueryStringer interface {
	URLQueryString() string
}

func (i *Inserter) insert(ctx context.Context, data response, params urlQueryStringer) error {
	return transact(ctx, i.DB, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, insertSQL)
		if err != nil {
			return err
		}
		defer stmt.Close()

		blob, err := json.Marshal(data)
		if err != nil {
			return err
		}

		var urlqp string
		if params != nil {
			urlqp = params.URLQueryString()
		}
		res, err := stmt.ExecContext(ctx, fball.EP_Countries, urlqp, data.When(), blob)
		if err != nil {
			return err
		}
		_, err = res.RowsAffected()
		return err

	})
}
