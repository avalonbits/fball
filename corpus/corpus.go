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
	"encoding/json"
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
	handle *db.Handle
}

func New(fballc *client.Client, logger *log.Logger, dbs *sql.DB) Corpus {
	return Corpus{
		logger: logger,
		fballc: fballc,
		handle: &db.Handle{DB: dbs},
	}
}

func (c Corpus) Timezone(ctx context.Context) ([]fball.TimezoneResponse, error) {
	// Query the timezone from the database.
	tzResp := []fball.TimezoneResponse{}
	err := c.handle.Query(ctx, fball.EP_Timezone, db.NoParams{}, 1, db.Range{}, func(data []byte) error {
		tr := fball.TimezoneResponse{}
		if err := json.Unmarshal(data, &tr); err != nil {
			return err
		}
		tzResp = append(tzResp, tr)
		return nil
	})

	if err == nil && len(tzResp) != 0 && rp_Infinite.Valid(time.Now(), tzResp[0].When()) {
		return tzResp, nil
	} else if err != nil {
		c.logger.Printf("WARNING - query error for timezone: %v", err)
	}

	// No timezone found, let's retrieve it from the api server.
	trQ, err := c.fballc.Timezone(ctx)
	if err != nil {
		// We tolerate stale data if the api call fails. We still log it.
		if len(tzResp) != 0 {
			c.logger.Printf("WARNING - unable to query timezone: %v.", err)
			c.logger.Printf("WARNING - returning stale data for countries.")
			return tzResp, nil
		} else {
			return nil, err
		}
	}

	// Ok, so now we can store it back in the database. Note that if storing fails, we still want to
	// return the data since the db is just a cache.
	if err := c.handle.Insert(ctx, fball.EP_Timezone, trQ[0], db.NoParams{}); err != nil {
		c.logger.Printf("ERROR - unable to write timezone to cache: %v", err)
	}

	return trQ, nil
}

func (c Corpus) Country(ctx context.Context, cp client.CountryParams) ([]fball.CountryResponse, error) {
	// Query the countries from the database.
	crResp := []fball.CountryResponse{}
	err := c.handle.Query(ctx, fball.EP_Countries, cp, 1, db.Range{}, func(data []byte) error {
		cr := fball.CountryResponse{}
		if err := json.Unmarshal(data, &cr); err != nil {
			return err
		}
		crResp = append(crResp, cr)
		return nil
	})

	if err == nil && len(crResp) != 0 && rp_OneDay.Valid(time.Now(), crResp[0].When()) {
		return crResp, nil
	} else if err != nil {
		c.logger.Printf("WARNING - query error for countries: %v", err)
	}

	// Either the data is not available or it has expired.
	crQ, err := c.fballc.Country(ctx, cp)
	if err != nil {
		// We tolerate stale data if the api call fails. We still log it.
		if len(crResp) != 0 {
			c.logger.Printf("WARNING - unable to query countries: %v.", err)
			c.logger.Printf("WARNING - returning stale data for countries.")
			return crResp, nil
		} else {
			return nil, err
		}
	}

	if err := c.handle.Insert(ctx, fball.EP_Countries, crQ[0], cp); err != nil {
		c.logger.Printf("ERROR - unable to write country to cache: %v", err)
	}

	return crQ, nil
}
