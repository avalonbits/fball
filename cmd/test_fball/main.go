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

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/avalonbits/fball"
	"github.com/kr/pretty"
)

var (
	key = flag.String("key", "", "API key for football-api.")
)

func main() {
	flag.Parse()

	c := fball.NewClient(*key, &http.Client{Timeout: 10 * time.Second})

	ctx := context.Background()
	tr, err := c.Timezone(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.Sprint(tr))

	cr, err := c.Country(ctx, fball.CountryParams{})
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.Sprint(cr))

	sr, err := c.Season(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.Sprint(sr))

	lir, err := c.LeagueInfo(ctx, fball.LeagueInfoParams{})
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.Sprint(lir))

	tir, err := c.TeamInfo(ctx, fball.TeamInfoParams{
		Country: "Brazil",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.Sprint(tir))

	tsr, err := c.TeamStats(ctx, fball.TeamStatsParams{
		League: "71",
		Season: "2020",
		Team:   "123",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.Sprint(tsr))

	vr, err := c.Venue(ctx, fball.VenueParams{
		Country: "Brazil",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.Sprint(vr))

	sp, err := c.Standings(ctx, fball.StandingsParams{
		League: "71",
		Season: "2020",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.Sprint(sp))

	rn, err := c.Round(ctx, fball.RoundParams{
		League:  "71",
		Season:  "2020",
		Current: "false",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.Sprint(rn))

	fix, err := c.FixtureInfo(ctx, fball.FixtureInfoParams{
		League: "71",
		Season: "2020",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.Sprint(fix))

	h2h, err := c.Head2Head(ctx, fball.Head2HeadParams{
		H2H:    "147-144",
		League: "71",
		Season: "2020",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.Sprint(h2h))

	fsr, err := c.FixtureStats(ctx, fball.FixtureStatsParams{
		Fixture: "328362",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.Sprint(fsr))

	er, err := c.Event(ctx, fball.EventParams{
		Fixture: "328362",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.Sprint(er))

	lr, err := c.Lineup(ctx, fball.LineupParams{
		Fixture: "328362",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.Sprint(lr))

	psr, err := c.PlayerStats(ctx, fball.PlayerStatsParams{
		Fixture: "328362",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.Sprint(psr))
}
