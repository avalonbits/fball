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

type Corpus struct {
	logger *log.Logger
	fballc *Client
	cache  *cache
}

func NewCorpus(fballc *Client, logger *log.Logger, dbs *sql.DB) *Corpus {
	return &Corpus{
		logger: logger,
		fballc: fballc,
		cache:  &cache{DB: dbs},
	}
}

func (c *Corpus) Timezone(ctx context.Context) ([]TimezoneResponse, error) {
	return c.getTimezoneResponse(ctx, EP_Timezone, 1, tRange{}, rp_Infinite, NoParams{})
}

type CountryParams struct {
	Name   string
	Code   string
	Search string
}

func (c *Corpus) Country(ctx context.Context, cp CountryParams) ([]CountryResponse, error) {
	return c.getCountryResponse(ctx, EP_Countries, 1, tRange{}, rp_OneDay, toURLQueryString{cp})
}

func (c *Corpus) Season(ctx context.Context) ([]SeasonResponse, error) {
	return c.getSeasonResponse(ctx, EP_Season, 1, tRange{}, rp_OneDay, NoParams{})
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

func (c *Corpus) LeagueInfo(ctx context.Context, params LeagueInfoParams) ([]LeagueInfoResponse, error) {
	return c.getLeagueInfoResponse(ctx, EP_LeagueInfo, 1, tRange{}, rp_OneHour, toURLQueryString{params})
}

type TeamInfoParams struct {
	ID      string
	Name    string
	League  string
	Season  string
	Country string
	Search  string
}

func (c *Corpus) TeamInfo(ctx context.Context, params TeamInfoParams) ([]TeamInfoResponse, error) {
	return c.getTeamInfoResponse(ctx, EP_TeamInfo, 1, tRange{}, rp_OneDay, toURLQueryString{params})
}

type TeamStatsParams struct {
	League string
	Season string
	Team   string
	Date   string
}

func (c *Corpus) TeamStats(ctx context.Context, params TeamStatsParams) ([]TeamStatsResponse, error) {
	return c.getTeamStatsResponse(ctx, EP_TeamStats, 1, tRange{}, rp_OneDay, toURLQueryString{params})
}

type VenueParams struct {
	ID      string
	Name    string
	City    string
	Country string
	Search  string
}

func (c *Corpus) Venue(ctx context.Context, params VenueParams) ([]VenueResponse, error) {
	return c.getVenueResponse(ctx, EP_Venue, 1, tRange{}, rp_OneDay, toURLQueryString{params})
}

type StandingsParams struct {
	League string
	Season string
	Team   string
}

func (c *Corpus) Standings(ctx context.Context, params StandingsParams) ([]StandingsResponse, error) {
	return c.getStandingsResponse(ctx, EP_Standings, 1, tRange{}, rp_OneHour, toURLQueryString{params})
}

type RoundParams struct {
	League  string
	Season  string
	Current string
}

func (c *Corpus) Round(ctx context.Context, params RoundParams) ([]RoundResponse, error) {
	return c.getRoundResponse(ctx, EP_Round, 1, tRange{}, rp_OneDay, toURLQueryString{params})
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

func (c *Corpus) FixtureInfo(ctx context.Context, params FixtureInfoParams) ([]FixtureInfoResponse, error) {
	return c.getFixtureInfoResponse(ctx, EP_FixtureInfo, 1, tRange{}, rp_OneMinute, toURLQueryString{params})
}

type refreshPolicy time.Duration

func (rp refreshPolicy) Valid(now time.Time, tsnano int64) bool {
	return now.UTC().Sub(time.Unix(0, tsnano)) < time.Duration(rp)
}

const (
	rp_OneMinute = refreshPolicy(time.Minute)
	rp_OneHour   = refreshPolicy(time.Hour)
	rp_OneDay    = refreshPolicy(86400 * time.Second)
	rp_Infinite  = refreshPolicy(1<<63 - 1)
)
