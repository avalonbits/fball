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
	"reflect"
	"sort"
	"strings"
	"text/template"
	"time"
)

const (
	EP_Timezone     = "/timezone"
	EP_Countries    = "/countries"
	EP_Season       = "/leagues/seasons"
	EP_LeagueInfo   = "/leagues"
	EP_TeamInfo     = "/teams"
	EP_TeamStats    = "/teams/statistics"
	EP_Venue        = "/venues"
	EP_Standings    = "/standings"
	EP_Round        = "/fixtures/rounds"
	EP_FixtureInfo  = "/fixtures"
	EP_Head2Head    = "/fixtures/headtohead"
	EP_FixtureStats = "/fixtures/statistics"
	EP_FixtureEvent = "/fixtures/events"
	EP_Lineup       = "/fixtures/lineups"
	EP_PlayerStats  = "/fixtures/players"
)

const NoParams = ""

type tRange struct {
	Latest   time.Time
	Earliest time.Time
}

func (r tRange) UnixNano() (top, bottom int64) {
	if r.Latest.IsZero() {
		top = time.Now().UTC().UnixNano()
	} else {
		top = r.Latest.UTC().UnixNano()
	}
	if !r.Earliest.IsZero() {
		bottom = r.Earliest.UTC().UnixNano()
	}
	return
}

func (r tRange) IsZero() bool {
	return r.Latest.IsZero() && r.Earliest.IsZero()
}

func ToURLQueryString(data any) string {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("expected a struct, got %v", v.Kind()))
	}

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
	sort.Strings(strs)

	return strings.Join(strs, "&")
}
