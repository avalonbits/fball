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
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/avalonbits/fball"
	"github.com/kr/pretty"
	"go.uber.org/ratelimit"

	_ "github.com/mattn/go-sqlite3"
)

var (
	key = flag.String("key", "", "API key for football-api.")
	db  = flag.String("db", "", "Path to sqlite database.")
)

func main() {
	flag.Parse()

	logger := log.New(os.Stderr, "fball - ", log.LstdFlags|log.Lshortfile)
	limit := ratelimit.New(10, ratelimit.Per(time.Minute))
	DB, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc&_journal_mode=WAL", *db))
	if err != nil {
		logger.Fatal(err)
	}
	defer DB.Close()

	c := fball.NewCorpus(
		fball.NewClient(*key, limit, &http.Client{Timeout: 10 * time.Second}, logger),
		logger,
		DB,
	)

	ctx := context.Background()
	tr, err := c.Timezone(ctx)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(pretty.Sprint(tr))

	cr, err := c.Country(ctx, fball.CountryParams{})
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(pretty.Sprint(cr))

	sr, err := c.Season(ctx)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(pretty.Sprint(sr))

	lir, err := c.LeagueInfo(ctx, fball.LeagueInfoParams{})
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(pretty.Sprint(lir))

	tir, err := c.TeamInfo(ctx, fball.TeamInfoParams{
		Country: "Brazil",
	})
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(pretty.Sprint(tir))

	tsr, err := c.TeamStats(ctx, fball.TeamStatsParams{
		League: "71",
		Season: "2020",
		Team:   "123",
	})
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(pretty.Sprint(tsr))

	vr, err := c.Venue(ctx, fball.VenueParams{
		Country: "Brazil",
	})
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(pretty.Sprint(vr))

	sp, err := c.Standings(ctx, fball.StandingsParams{
		League: "71",
		Season: "2020",
	})
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(pretty.Sprint(sp))

	rn, err := c.Round(ctx, fball.RoundParams{
		League:  "71",
		Season:  "2020",
		Current: "false",
	})
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(pretty.Sprint(rn))

	fix, err := c.FixtureInfo(ctx, fball.FixtureInfoParams{
		League: "71",
		Season: "2020",
	})
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(pretty.Sprint(fix))

	h2h, err := c.Head2Head(ctx, fball.Head2HeadParams{
		H2H:    "147-144",
		League: "71",
		Season: "2020",
	})
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(pretty.Sprint(h2h))

}
