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
)

type noParams struct{}

func (np noParams) URLQueryString() string {
	return ""
}

type urlQueryStringer interface {
	URLQueryString() string
}

func transact(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) (dberr error) {
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
