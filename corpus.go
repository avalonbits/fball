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
	"log"
	"time"
)

type Corpus struct {
	logger   *log.Logger
	fballc   *Client
	cache    *cache
	useStale bool
}

func NewCorpus(fballc *Client, logger *log.Logger, dbs *sql.DB) *Corpus {
	return &Corpus{
		logger: logger,
		fballc: fballc,
		cache:  &cache{DB: dbs},
	}
}

func (c *Corpus) WithStale(stale bool) *Corpus {
	c.useStale = stale
	return c
}

func (c *Corpus) Timezone(ctx context.Context) (TimezoneResponse, error) {
	tr := TimezoneResponse{}
	err := c.get(ctx, EP_Timezone, RP_Infinite, NoParams{}, &tr)
	return tr, err
}

type CountryParams struct {
	Name   string
	Code   string
	Search string
}

func (c *Corpus) Country(ctx context.Context, cp CountryParams) (CountryResponse, error) {
	cr := CountryResponse{}
	err := c.get(ctx, EP_Countries, RP_OneDay, toURLQueryString{cp}, &cr)
	return cr, err
}

func (c *Corpus) Season(ctx context.Context) (SeasonResponse, error) {
	sr := SeasonResponse{}
	err := c.get(ctx, EP_Season, RP_OneDay, NoParams{}, &sr)
	return sr, err

}

type LeagueInfoParams struct {
	ID      string
	Name    string
	Country string
	Code    string
	Season  string
	Team    string
	Type    string
	Current string
	Search  string
	Last    string
}

func (c *Corpus) LeagueInfo(ctx context.Context, params LeagueInfoParams) (LeagueInfoResponse, error) {
	lir := LeagueInfoResponse{}
	err := c.get(ctx, EP_LeagueInfo, RP_OneHour, toURLQueryString{params}, &lir)
	return lir, err
}

type TeamInfoParams struct {
	ID      string
	Name    string
	League  string
	Season  string
	Country string
	Search  string
}

func (c *Corpus) TeamInfo(ctx context.Context, params TeamInfoParams) (TeamInfoResponse, error) {
	tir := TeamInfoResponse{}
	err := c.get(ctx, EP_TeamInfo, RP_OneDay, toURLQueryString{params}, &tir)
	return tir, err
}

type TeamStatsParams struct {
	League string
	Season string
	Team   string
	Date   string
}

func (c *Corpus) TeamStats(ctx context.Context, params TeamStatsParams) (TeamStatsResponse, error) {
	tsr := TeamStatsResponse{}
	err := c.get(ctx, EP_TeamStats, RP_OneDay, toURLQueryString{params}, &tsr)
	return tsr, err
}

type VenueParams struct {
	ID      string
	Name    string
	City    string
	Country string
	Search  string
}

func (c *Corpus) Venue(ctx context.Context, params VenueParams) (VenueResponse, error) {
	vr := VenueResponse{}
	err := c.get(ctx, EP_Venue, RP_OneDay, toURLQueryString{params}, &vr)
	return vr, err
}

type StandingsParams struct {
	League string
	Season string
	Team   string
}

func (c *Corpus) Standings(ctx context.Context, params StandingsParams) (StandingsResponse, error) {
	sr := StandingsResponse{}
	err := c.get(ctx, EP_Standings, RP_OneHour, toURLQueryString{params}, &sr)
	return sr, err
}

type RoundParams struct {
	League  string
	Season  string
	Current string
}

func (c *Corpus) Round(ctx context.Context, params RoundParams) (RoundResponse, error) {
	rr := RoundResponse{}
	err := c.get(ctx, EP_Round, RP_OneDay, toURLQueryString{params}, &rr)
	return rr, err
}

type FixtureInfoParams struct {
	ID       string
	Live     string
	Date     string
	League   string
	Season   string
	Team     string
	Last     string
	Next     string
	From     string
	To       string
	Round    string
	Status   string
	Timezone string
}

func (c *Corpus) FixtureInfo(ctx context.Context, params FixtureInfoParams) (FixtureInfoResponse, error) {
	fir := FixtureInfoResponse{}
	err := c.get(ctx, EP_FixtureInfo, RP_OneMinute, toURLQueryString{params}, &fir)
	return fir, err
}

type Head2HeadParams struct {
	H2H      string
	Date     string
	League   string
	Season   string
	Last     string
	Next     string
	From     string
	To       string
	Status   string
	Timezone string
}

func (c *Corpus) Head2Head(ctx context.Context, params Head2HeadParams) (Head2HeadResponse, error) {
	h2hr := Head2HeadResponse{}
	err := c.get(ctx, EP_Head2Head, RP_OneMinute, toURLQueryString{params}, &h2hr)
	return h2hr, err
}

type FixtureStatsParams struct {
	Fixture string
	Team    string
	Type    string
}

func (c *Corpus) FixtureStats(ctx context.Context, params FixtureStatsParams) (FixtureStatsResponse, error) {
	fsr := FixtureStatsResponse{}
	err := c.get(ctx, EP_FixtureStats, RP_OneMinute, toURLQueryString{params}, &fsr)
	return fsr, err
}

type EventParams struct {
	Fixture string
	Team    string
	Player  string
	Type    string
}

func (c *Corpus) Event(ctx context.Context, params EventParams) (EventResponse, error) {
	er := EventResponse{}
	err := c.get(ctx, EP_FixtureEvent, RP_OneMinute, toURLQueryString{params}, &er)
	return er, err
}

type LineupParams struct {
	Fixture string
	Team    string
	Player  string
	Type    string
}

func (c *Corpus) Lineup(ctx context.Context, params LineupParams) (LineupResponse, error) {
	lr := LineupResponse{}
	err := c.get(ctx, EP_Lineup, RP_15Minutes, toURLQueryString{params}, &lr)
	return lr, err
}

type PlayerStatsParams struct {
	Fixture string
	Team    string
}

func (c *Corpus) PlayerStats(ctx context.Context, params PlayerStatsParams) (PlayerStatsResponse, error) {
	psr := PlayerStatsResponse{}
	err := c.get(ctx, EP_PlayerStats, RP_OneMinute, toURLQueryString{params}, &psr)
	return psr, err
}

type refreshPolicy time.Duration

func (rp refreshPolicy) Valid(now time.Time, tsnano int64) bool {
	return now.UTC().Sub(time.Unix(0, tsnano)) < time.Duration(rp)
}

func (rp refreshPolicy) Range(now time.Time) tRange {
	now = now.UTC()
	return tRange{
		Latest:   now,
		Earliest: now.Add(-time.Duration(rp)),
	}
}

const (
	RP_OneMinute = refreshPolicy(time.Minute)
	RP_15Minutes = refreshPolicy(15 * time.Minute)
	RP_OneHour   = refreshPolicy(time.Hour)
	RP_OneDay    = refreshPolicy(86400 * time.Second)
	RP_Infinite  = refreshPolicy(1<<63 - 1)
)

func Get[T Response](
	ctx context.Context, c *Corpus, endpoint string, policy refreshPolicy, params urlQueryStringer, data T) error {
	q1 := time.Now()
	found := false
	pRange := tRange{}
	if !c.useStale {
		pRange = policy.Range(q1)
	}
	err := c.cache.Query(ctx, endpoint, params, 1, pRange, func(bs []byte) error {
		if err := json.Unmarshal(bs, data); err != nil {
			return err
		}
		found = true
		return nil
	})
	q2 := time.Now()
	c.logger.Printf("INFO - %q query time: %dms", endpoint, q2.Sub(q1)/time.Millisecond)

	if err == nil && found && !c.useStale {
		return nil
	} else if err != nil {
		c.logger.Printf("WARNING - query error: %v", err)
	}

	// Either the data is not available or it has expired.
	s1 := time.Now()
	err = c.fballc.Get(ctx, endpoint, data, params)
	s2 := time.Now()
	c.logger.Printf("INFO - %q api call time: %dms", endpoint, s2.Sub(s1)/time.Millisecond)

	if err != nil {
		// There was an error in the api call. If we can, return stale data.
		if !c.useStale {
			return err
		} else if found {
			c.logger.Printf("WARNING - returning stale data because of api call error: %v", err)
		}
	} else {
		i1 := time.Now()
		err = c.cache.Insert(ctx, endpoint, data, params)
		i2 := time.Now()
		c.logger.Printf("INFO - %q insert time: %dms", endpoint, i2.Sub(i1)/time.Millisecond)

		if err != nil {
			c.logger.Printf("ERROR - unable to write country to cache: %v", err)
		}
	}
	return nil
}

func (c *Corpus) get(
	ctx context.Context, endpoint string, policy refreshPolicy, params urlQueryStringer, data Response) error {
	q1 := time.Now()
	found := false
	pRange := tRange{}
	if !c.useStale {
		pRange = policy.Range(q1)
	}
	err := c.cache.Query(ctx, endpoint, params, 1, pRange, func(bs []byte) error {
		if err := json.Unmarshal(bs, data); err != nil {
			return err
		}
		found = true
		return nil
	})
	q2 := time.Now()
	c.logger.Printf("INFO - %q query time: %dms", endpoint, q2.Sub(q1)/time.Millisecond)

	if err == nil && found && !c.useStale {
		return nil
	} else if err != nil {
		c.logger.Printf("WARNING - query error: %v", err)
	}

	// Either the data is not available or it has expired.
	s1 := time.Now()
	err = c.fballc.Get(ctx, endpoint, data, params)
	s2 := time.Now()
	c.logger.Printf("INFO - %q api call time: %dms", endpoint, s2.Sub(s1)/time.Millisecond)

	if err != nil {
		// There was an error in the api call. If we can, return stale data.
		if !c.useStale {
			return err
		} else if found {
			c.logger.Printf("WARNING - returning stale data because of api call error: %v", err)
		}
	} else {
		i1 := time.Now()
		err = c.cache.Insert(ctx, endpoint, data, params)
		i2 := time.Now()
		c.logger.Printf("INFO - %q insert time: %dms", endpoint, i2.Sub(i1)/time.Millisecond)

		if err != nil {
			c.logger.Printf("ERROR - unable to write country to cache: %v", err)
		}
	}
	return nil
}
