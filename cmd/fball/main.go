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
	"log"
	"net/http"
	"os"
	"time"

	"git.cana.pw/avalonbits/fball/client"
	"git.cana.pw/avalonbits/fball/corpus"
	"github.com/kr/pretty"
	"go.uber.org/ratelimit"
)

var (
	key = flag.String("key", "", "API key for football-api.")
	db  = flag.String("db", "", "Path to sqlite database.")
)

func main() {
	flag.Parse()

	logger := log.New(os.Stderr, "fball - ", log.LstdFlags|log.Lshortfile)
	limit := ratelimit.New(10, ratelimit.Per(time.Minute))
	c := corpus.New(
		client.NewClient(*key, limit, &http.Client{Timeout: 10 * time.Second}, logger),
		nil,
	)

	ctx := context.Background()
	tr, err := c.Timezone(ctx)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(pretty.Sprint(tr))
}