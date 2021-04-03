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
	"time"

	"git.cana.pw/avalonbits/fball"
	"git.cana.pw/avalonbits/fball/client"
	"git.cana.pw/avalonbits/fball/db"
)

type refreshPolicy time.Duration

func (rp refreshPolicy) Valid(now time.Time, tsnano int64) bool {
	return now.UTC().Sub(time.Unix(0, tsnano)) < time.Duration(rp)
}

const (
	rp_OneDay   = refreshPolicy(86400 * time.Second)
	rp_Infinite = refreshPolicy(1<<63 - 1)
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
	// Query the timezone from the database.
	tr, err := c.query.Timezone(ctx, 1, db.Range{})
	if err == nil && len(tr) != 0 {
		if len(tr) != 0 && rp_Infinite.Valid(time.Now(), tr[0].When()) {
			return tr, nil
		}
	}

	// No timezone found, let's retrieve it from the api server.
	trQ, err := c.fballc.Timezone(ctx)
	if err != nil {
		// We tolerate stale data if the api call fails. We still log it.
		if len(tr) != 0 {
			c.logger.Printf("WARNING - unable to query timezone: %v.", err)
			c.logger.Printf("WARNING - returning stale data for countries.")
			return tr, nil
		} else {
			return nil, err
		}
	}

	// Ok, so now we can store it back in the database. Note that if storing fails, we still want to
	// return the data since the db is just a cache.
	if err := c.insert.Timezone(ctx, trQ[0]); err != nil {
		c.logger.Printf("ERROR - unable to write timezone to cache: %v", err)
	}

	return trQ, nil
}

func (c Corpus) Country(ctx context.Context, cp client.CountryParams) ([]fball.CountryResponse, error) {
	// Query the countries from the database.
	cr, err := c.query.Country(ctx, cp, 1, db.Range{})
	if err == nil && len(cr) != 0 {
		// Country data is valid for one day.
		if len(cr) != 0 && rp_OneDay.Valid(time.Now(), cr[0].When()) {
			return cr, nil
		}
	}

	// Either the data is not available or it has expired.
	crQ, err := c.fballc.Country(ctx, cp)
	if err != nil {
		// We tolerate stale data if the api call fails. We still log it.
		if len(cr) != 0 {
			c.logger.Printf("WARNING - unable to query countries: %v.", err)
			c.logger.Printf("WARNING - returning stale data for countries.")
			return cr, nil
		} else {
			return nil, err
		}
	}

	if err := c.insert.Country(ctx, crQ[0], cp); err != nil {
		c.logger.Printf("ERROR - unable to write country to cache: %v", err)
	}

	return crQ, nil
}
