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
	"fmt"
	"strings"

	"github.com/kr/pretty"
)

type TimezoneResponse struct {
	CommonResponse
	Timezone []string `json:"response"`
}

type CountryResponse struct {
	CommonResponse
	Country []Country `json:"response"`
}

type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Flag string `json:"flag"`
}

type SeasonResponse struct {
	CommonResponse
	Season []int `json:"response"`
}

type LeagueInfoResponse struct {
	CommonResponse
	LeagueInfo []LeagueInfo `json:"response"`
}

type LeagueInfo struct {
	League struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		Logo string `json:"logo"`
	} `json:"league"`
	Country Country `json:"country"`
	Seasons []struct {
		Year     int    `json:"year"`
		Start    string `json:"start"`
		End      string `json:"end"`
		Current  bool   `json:"current"`
		Coverage struct {
			Fixtures struct {
				Events             bool `json:"events"`
				Lineups            bool `json:"lineups"`
				StatisticsFixtures bool `json:"statistics_fixtures"`
				StatisticsPlayers  bool `json:"statistics_players"`
			} `json:"fixtures"`
			Standings   bool `json:"standings"`
			Players     bool `json:"players"`
			TopScorers  bool `json:"top_scorers"`
			TopAssists  bool `json:"top_assists"`
			TopCards    bool `json:"top_cards"`
			Predictions bool `json:"predictions"`
			Odds        bool `json:"odds"`
		} `json:"coverage"`
	} `json:"seasons"`
}

type TeamInfoResponse struct {
	CommonResponse
	TeamInfo []TeamInfo `json:"response"`
}

type TeamInfo struct {
	Team  TeamData `json:"team"`
	Venue Venue    `json:"venue"`
}

type TeamData struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Country  string `json:"country"`
	Founded  int    `json:"founded"`
	National bool   `json:"national"`
	Logo     string `json:"logo"`
}

type Venue struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	City     string `json:"city"`
	Capacity int    `json:"capacity"`
	Surface  string `json:"surface"`
	Image    string `json:"image"`
}

type CommonResponse struct {
	Get        string      `json:"get"`
	Parameters interface{} `json:"parameters"`
	Errors     interface{} `json:"errors"`
	Results    int         `json:"results"`
	Paging     PagingToken `json:"paging"`
	Timestamp  int64
}

func (cr *CommonResponse) SetWhen(timestamp int64) {
	cr.Timestamp = timestamp
}

func (cr CommonResponse) When() int64 {
	return cr.Timestamp
}

func (cr CommonResponse) Err() error {
	if cr.Errors == nil {
		return nil
	}

	// If there are no errros, Errors gets parsed as a []interface. So let's check that first.
	nErrs, ok := cr.Errors.([]interface{})
	if ok {
		if len(nErrs) != 0 {
			return fmt.Errorf("%s", pretty.Sprint(nErrs))
		} else {
			return nil
		}
	}

	// Now check if actual errors were returned.
	errs, ok := cr.Errors.(map[string]interface{})
	if !ok {
		return fmt.Errorf("%s", pretty.Sprint(errs))
	}
	if len(errs) == 0 {
		return nil
	}

	errList := make([]string, 0, len(errs))
	for k, v := range errs {
		errList = append(errList, fmt.Sprintf("\t%q: %q", k, v))
	}
	return fmt.Errorf("\nerrors:\n%s\n", strings.Join(errList, "\n"))
}

type PagingToken struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}
