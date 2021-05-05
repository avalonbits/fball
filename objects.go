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

type RoundResponse struct {
	CommonResponse
	Rounds []string `json:"response"`
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
	League  League  `json:"league"`
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

type Ranking struct {
	Rank        int        `json:"rank"`
	Team        TeamData   `json:"team"`
	Points      int        `json:"points"`
	Goalsdiff   int        `json:"goalsDiff"`
	Group       string     `json:"group"`
	Form        string     `json:"form"`
	Status      string     `json:"status"`
	Description string     `json:"description"`
	All         RankTotals `json:"all"`
	Home        RankTotals `json:"home"`
	Away        RankTotals `json:"away"`
	Update      string     `json:"update"`
}

type League struct {
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Country  string      `json:"country"`
	Logo     string      `json:"logo"`
	Flag     string      `json:"flag"`
	Season   int         `json:"season"`
	Round    string      `json:"round"`
	Rankings [][]Ranking `json:"standings"`
}

type TeamInfoResponse struct {
	CommonResponse
	TeamInfo []TeamInfo `json:"response"`
}

type VenueResponse struct {
	CommonResponse
	Venue []Venue `json:"response"`
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
	Update   string
}

type Venue struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	City     string `json:"city"`
	Capacity int    `json:"capacity"`
	Surface  string `json:"surface"`
	Country  string `json:"country"`
	Image    string `json:"image"`
}

type StandingsResponse struct {
	CommonResponse
	Standings []struct {
		League League `json:"league"`
	} `json:"response"`
}

type RankTotals struct {
	Played int `json:"played"`
	Win    int `json:"win"`
	Draw   int `json:"draw"`
	Lose   int `json:"lose"`
	Goals  struct {
		For     int `json:"for"`
		Against int `json:"against"`
	} `json:"goals"`
}

type TeamStatsResponse struct {
	CommonResponse

	TeamStats struct {
		League   League   `json:"league"`
		Team     TeamData `json:"team"`
		Form     string   `json:"form"`
		Fixtures struct {
			Played Totals `json:"played"`
			Wins   Totals `json:"wins"`
			Draws  Totals `json:"draws"`
			Loses  Totals `json:"loses"`
		} `json:"fixtures"`
		Goals struct {
			For struct {
				Total   Totals   `json:"total"`
				Average TotalStr `json:"average"`
				Minute  GameTime `json:"minute"`
			} `json:"for"`
			Against struct {
				Total   Totals   `json:"total"`
				Average TotalStr `json:"average"`
				Minute  GameTime `json:"minute"`
			} `json:"against"`
		} `json:"goals"`
		Biggest struct {
			Streak struct {
				Wins  int `json:"wins"`
				Draws int `json:"draws"`
				Loses int `json:"loses"`
			} `json:"streak"`
			Wins  TotalStr `json:"wins"`
			Losev TotalStr `json:"loses"`
			Goals struct {
				For     Totals `json:"for"`
				Against Totals `json:"against"`
			} `json:"goals"`
		} `json:"biggest"`
		CleanSheet    Totals `json:"clean_sheet"`
		FailedToScore Totals `json:"failed_to_score"`
		Penalty       struct {
			Scored TotalPercent `json:"scored"`
			Missed TotalPercent `json:"missed"`
			Total  int          `json:"total"`
		} `json:"penalty"`
		Lineups []struct {
			Formation string `json:"formation"`
			Played    int    `json:"played"`
		} `json:"lineups"`
		Cards struct {
			Yellow GameTime `json:"yellow"`
			Red    GameTime `json:"red"`
		} `json:"cards"`
	} `json:"response"`
}

type Totals struct {
	Home  int `json:"home"`
	Away  int `json:"away"`
	Total int `json:"total"`
}

type TotalStr struct {
	Home  string `json:"home"`
	Away  string `json:"away"`
	Total string `json:"total"`
}

type GameTime struct {
	P15  TotalPercent `json:"0-15"`
	P30  TotalPercent `json:"16-30"`
	P45  TotalPercent `json:"31-45"`
	P60  TotalPercent `json:"46-60"`
	P75  TotalPercent `json:"61-75"`
	P90  TotalPercent `json:"76-90"`
	P105 TotalPercent `json:"91-105"`
	P120 TotalPercent `json:"105-120"`
}

type TotalPercent struct {
	Total      int    `json:"total"`
	Percentage string `json:"percentage"`
}

type FixtureInfoResponse struct {
	CommonResponse

	FixtureInfo []Head2Head `json:"response"`
}

type Head2Head struct {
	Fixture Fixture `json:"fixture"`
	League  League  `json:"league"`
	Teams   struct {
		Home struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Logo   string `json:"logo"`
			Winner bool   `json:"winner"`
		} `json:"home"`
		Away struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Logo   string `json:"logo"`
			Winner bool   `json:"winner"`
		} `json:"away"`
	} `json:"teams"`
	Goals Totals `json:"goals"`
	Score struct {
		Halftime  Totals `json:"halftime"`
		Fulltime  Totals `json:"fulltime"`
		Extratime Totals `json:"extratime"`
		Penalty   Totals `json:"penalty"`
	} `json:"score"`
}

type Fixture struct {
	ID        int    `json:"id"`
	Referee   string `json:"referee"`
	Timezone  string `json:"timezone"`
	Date      string `json:"date"`
	Timestamp int64  `json:"timestamp"`
	Periods   struct {
		First  int `json:"first"`
		Second int `json:"second"`
	} `json:"periods"`
	Venue  Venue `json:"venue"`
	Status struct {
		Long    string `json:"long"`
		Short   string `json:"short"`
		Elapsed int    `json:"elapsed"`
	} `json:"status"`
}

type Head2HeadResponse struct {
	CommonResponse

	Head2Head []Head2Head `json:"response"`
}

type StatisticsResponse struct {
	CommonResponse

	Statistics []Statistics `json:"response"`
}

type Statistics struct {
	Team TeamData `json:"team"`
	Info []struct {
		Type  string `json:"type"`
		Value int    `json:"value"`
	}
}

type EventResponse struct {
	CommonResponse

	Event []Event `json:"response"`
}

type Event struct {
	Time struct {
		Elapsed int `json:"elapsed"`
		Extra   int `json:"extra"`
	} `json:"time"`
	Team     TeamData `json:"team"`
	Player   Player   `json:"player"`
	Assist   Player   `json:"assist"`
	Type     string   `json:"type"`
	Detail   string   `json:"detail"`
	Comments string   `json:"comments"`
}

type Player struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Number int    `json:"number"`
	Pos    string `json:"pos"`
}

type LineupResponse struct {
	CommonResponse

	Lineup []Lineup `json:"response"`
}

type Lineup struct {
	Team  TeamData `json:"team"`
	Coach struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Photo string `json:"photo"`
	} `json:"coach"`
	Formation string `json:"formation"`
	StartXI   []struct {
		Player Player `json:"player"`
	} `json:"startXI"`
	Substitutes []struct {
		Player Player `json:"player"`
	} `json:"substitutes"`
}

type PlayerStatsResponse struct {
	CommonResponse

	PlayerStats []PlayerStats `json:"response"`
}

type PlayerStats struct {
	Team    TeamData `json:"team"`
	Players []struct {
		Player     Player `json:"player"`
		Statistics []struct {
			Games struct {
				Minutes    int    `json:"minutes"`
				Number     int    `json:"number"`
				Position   string `json:"position"`
				Rating     string `json:"rating"`
				Captain    bool   `json:"captain"`
				Substitute bool   `json:"substitute"`
			} `json:"games"`
			Offsides interface{} `json:"offsides"`
			Shots    struct {
				Total int `json:"total"`
				On    int `json:"on"`
			} `json:"shots"`
			Goals struct {
				Total    int `json:"total"`
				Conceded int `json:"conceded"`
				Assists  int `json:"assists"`
				Saves    int `json:"saves"`
			} `json:"goals"`
			Passes struct {
				Total    int    `json:"total"`
				Key      int    `json:"key"`
				Accuracy string `json:"accuracy"`
			} `json:"passes"`
			Tackles struct {
				Total         int `json:"total"`
				Blocks        int `json:"blocks"`
				Interceptions int `json:"interceptions"`
			} `json:"tackles"`
			Duels struct {
				Total int `json:"total"`
				Won   int `json:"won"`
			} `json:"duels"`
			Dribbles struct {
				Attempts int `json:"attempts"`
				Success  int `json:"success"`
				Past     int `json:"past"`
			} `json:"dribbles"`
			Fouls struct {
				Drawn     int `json:"drawn"`
				Committed int `json:"committed"`
			} `json:"fouls"`
			Cards struct {
				Yellow int `json:"yellow"`
				Red    int `json:"red"`
			} `json:"cards"`
			Penalty struct {
				Won      int `json:"won"`
				Commited int `json:"commited"`
				Scored   int `json:"scored"`
				Missed   int `json:"missed"`
				Saved    int `json:"saved"`
			} `json:"penalty"`
		} `json:"statistics"`
	} `json:"players"`
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
