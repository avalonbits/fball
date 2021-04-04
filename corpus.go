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
	"log"
	"time"
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
	fballc *Client
	handle *Handle
}

func New(fballc *Client, logger *log.Logger, dbs *sql.DB) Corpus {
	return Corpus{
		logger: logger,
		fballc: fballc,
		handle: &Handle{DB: dbs},
	}
}

func (c Corpus) Timezone(ctx context.Context) ([]TimezoneResponse, error) {
	return c.getTimezoneResponse(ctx, EP_Timezone, 1, Range{}, rp_Infinite, NoParams{})
}

type CountryParams struct {
	Name   string
	Code   string
	Search string
}

func (cp CountryParams) URLQueryString() string {
	return StructToURLQueryString(cp)
}

func (c Corpus) Country(ctx context.Context, cp CountryParams) ([]CountryResponse, error) {
	return c.getCountryResponse(ctx, EP_Countries, 1, Range{}, rp_OneDay, cp)
}