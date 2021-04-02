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

package corpus

import (
	"context"
	"database/sql"
	"log"

	"git.cana.pw/avalonbits/fball"
	"git.cana.pw/avalonbits/fball/client"
	"git.cana.pw/avalonbits/fball/db"
)

type Corpus struct {
	logger *log.Logger
	fballc *client.Client
	query  *db.Querier
	insert *db.Inserter
}

func New(fballc *client.Client, logger *log.Logger, dbs *sql.DB) Corpus {
	return Corpus{
		logger: logger,
		fballc: fballc,
		query:  &db.Querier{DB: dbs},
		insert: &db.Inserter{DB: dbs},
	}
}

func (c Corpus) Timezone(ctx context.Context) ([]fball.TimezoneResponse, error) {
	// Query the timezone from the dabase.
	tr, err := c.query.Timezone(ctx, 1, db.Range{})
	if err != nil {
		return nil, err
	}

	// If it's there, there is nothing else to do.
	if len(tr) != 0 {
		return tr, nil
	}

	// No timezone found, let's retrieve it from the api server.
	tr, err = c.fballc.Timezone(ctx)
	if err != nil {
		return nil, err
	}

	if err := c.insert.Timezone(ctx, tr[0]); err != nil {
		c.logger.Printf("ERROR - unable to write timezone to cache: %v", err)
	}

	// Ok, so now we can store it back in the database. Note that if storing fails, we still want to
	// return the data since the db is just a cache.

	return tr, nil
}

func (c Corpus) Country(ctx context.Context, cp client.CountryParams) ([]fball.CountryResponse, error) {
	return c.fballc.Country(ctx, cp)
}
