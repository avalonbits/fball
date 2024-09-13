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
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/avalonbits/fball/object"
)

const (
	ep_Timezone     = "/timezone"
	ep_Countries    = "/countries"
	ep_Season       = "/leagues/seasons"
	ep_LeagueInfo   = "/leagues"
	ep_TeamInfo     = "/teams"
	ep_TeamStats    = "/teams/statistics"
	ep_Venue        = "/venues"
	ep_Standings    = "/standings"
	ep_Round        = "/fixtures/rounds"
	ep_FixtureInfo  = "/fixtures"
	ep_Head2Head    = "/fixtures/headtohead"
	ep_FixtureStats = "/fixtures/statistics"
	ep_FixtureEvent = "/fixtures/events"
	ep_Lineup       = "/fixtures/lineups"
	ep_PlayerStats  = "/fixtures/players"
)

// Doer is an interface for perfomring http requests.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// Client is an api-football.com client.
type Client struct {
	key  string
	doer Doer
}

// NewClient creates an api-football.com client. The key is the one provided by the
// service when you register it. limit will rate limit any request and if no logger
// is provided, a default logger is used.
func NewClient(key string, doer Doer) *Client {
	return &Client{
		doer: doer,
		key:  key,
	}
}

// Response is an interface for api-football.com responses.
type Response interface {
	// Err returns the error from the response, if any.
	Err() error

	// When returns the timestamp for the response.
	When() int64

	// SetWhen sets the timestamp for the response.
	SetWhen(int64)
}

const base = "https://v3.football.api-sports.io"

func (c *Client) Timezone(ctx context.Context) (object.TimezoneResponse, error) {
	tr := object.TimezoneResponse{}
	err := c.get(ctx, "/timezone", struct{}{}, &tr)
	return tr, err
}

type CountryParams struct {
	Name   string
	Code   string
	Search string
}

func (c *Client) Country(ctx context.Context, params CountryParams) (object.CountryResponse, error) {
	cr := object.CountryResponse{}
	err := c.get(ctx, ep_Countries, params, &cr)
	return cr, err
}

func (c *Client) Season(ctx context.Context) (object.SeasonResponse, error) {
	sr := object.SeasonResponse{}
	err := c.get(ctx, ep_Season, struct{}{}, &sr)
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

func (c *Client) LeagueInfo(ctx context.Context, params LeagueInfoParams) (object.LeagueInfoResponse, error) {
	lir := object.LeagueInfoResponse{}
	err := c.get(ctx, ep_LeagueInfo, params, &lir)
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

func (c *Client) TeamInfo(ctx context.Context, params TeamInfoParams) (object.TeamInfoResponse, error) {
	tir := object.TeamInfoResponse{}
	err := c.get(ctx, ep_TeamInfo, params, &tir)
	return tir, err
}

type TeamStatsParams struct {
	League string
	Season string
	Team   string
	Date   string
}

func (c *Client) TeamStats(ctx context.Context, params TeamStatsParams) (object.TeamStatsResponse, error) {
	tsr := object.TeamStatsResponse{}
	err := c.get(ctx, ep_TeamStats, params, &tsr)
	return tsr, err
}

type VenueParams struct {
	ID      string
	Name    string
	City    string
	Country string
	Search  string
}

func (c *Client) Venue(ctx context.Context, params VenueParams) (object.VenueResponse, error) {
	vr := object.VenueResponse{}
	err := c.get(ctx, ep_Venue, params, &vr)
	return vr, err
}

type StandingsParams struct {
	League string
	Season string
	Team   string
}

func (c *Client) Standings(ctx context.Context, params StandingsParams) (object.StandingsResponse, error) {
	sr := object.StandingsResponse{}
	err := c.get(ctx, ep_Standings, params, &sr)
	return sr, err
}

type RoundParams struct {
	League  string
	Season  string
	Current string
}

func (c *Client) Round(ctx context.Context, params RoundParams) (object.RoundResponse, error) {
	rr := object.RoundResponse{}
	err := c.get(ctx, ep_Round, params, &rr)
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

func (c *Client) FixtureInfo(ctx context.Context, params FixtureInfoParams) (object.FixtureInfoResponse, error) {
	fir := object.FixtureInfoResponse{}
	err := c.get(ctx, ep_FixtureInfo, params, &fir)
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

func (c *Client) Head2Head(ctx context.Context, params Head2HeadParams) (object.Head2HeadResponse, error) {
	h2hr := object.Head2HeadResponse{}
	err := c.get(ctx, ep_Head2Head, params, &h2hr)
	return h2hr, err
}

type FixtureStatsParams struct {
	Fixture string
	Team    string
	Type    string
}

func (c *Client) FixtureStats(ctx context.Context, params FixtureStatsParams) (object.FixtureStatsResponse, error) {
	fsr := object.FixtureStatsResponse{}
	err := c.get(ctx, ep_FixtureStats, params, &fsr)
	return fsr, err
}

type EventParams struct {
	Fixture string
	Team    string
	Player  string
	Type    string
}

func (c *Client) Event(ctx context.Context, params EventParams) (object.EventResponse, error) {
	er := object.EventResponse{}
	err := c.get(ctx, ep_FixtureEvent, params, &er)
	return er, err
}

type LineupParams struct {
	Fixture string
	Team    string
	Player  string
	Type    string
}

func (c *Client) Lineup(ctx context.Context, params LineupParams) (object.LineupResponse, error) {
	lr := object.LineupResponse{}
	err := c.get(ctx, ep_Lineup, params, &lr)
	return lr, err
}

type PlayerStatsParams struct {
	Fixture string
	Team    string
}

func (c *Client) PlayerStats(ctx context.Context, params PlayerStatsParams) (object.PlayerStatsResponse, error) {
	psr := object.PlayerStatsResponse{}
	err := c.get(ctx, ep_PlayerStats, params, &psr)
	return psr, err
}

// Get will perform a GET request against the api-football service.
// The response is returned in the data out param.
func (c *Client) get(ctx context.Context, endpoint string, params any, data Response) error {
	if data == nil {
		return fmt.Errorf("inalid data: must be non-nil")
	}
	if len(endpoint) == 0 {
		return fmt.Errorf("invalid endpoint: empty string")
	}

	queryStr := toURLQueryString(params)
	url := base + endpoint
	if queryStr != "" {
		url += "?"
		url += queryStr
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	now := time.Now().UTC().UnixNano()
	req.Header.Set("X-RapidAPI-Key", c.key)
	resp, err := c.doer.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(data); err != nil {
		return err
	}
	data.SetWhen(now)

	return data.Err()
}

func toURLQueryString(data any) string {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("expected a struct, got %v", v.Kind()))
	}
	if v.NumField() == 0 {
		return ""
	}

	t := reflect.TypeOf(data)
	strs := []string{}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() != reflect.String {
			continue
		}
		val := f.Interface().(string)
		if val == "" {
			continue
		}

		key := strings.ToLower(t.Field(i).Name)
		strs = append(strs, template.URLQueryEscaper(key)+"="+template.URLQueryEscaper(val))
	}
	if len(strs) == 0 {
		return ""
	}

	sort.Strings(strs)
	return strings.Join(strs, "&")
}
