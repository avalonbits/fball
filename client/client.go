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

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"git.cana.pw/avalonbits/fball"
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type limiter interface {
	Take() time.Time
}

type Client struct {
	key    string
	doer   Doer
	limit  limiter
	logger *log.Logger
}

func NewClient(key string, limit limiter, doer Doer, logger *log.Logger) *Client {
	if logger == nil {
		logger = log.Default()
	}
	return &Client{
		doer:   doer,
		key:    key,
		logger: logger,
		limit:  limit,
	}
}

func (c *Client) Timezone(ctx context.Context) ([]fball.TimezoneResponse, error) {
	tr := fball.TimezoneResponse{}
	if err := c.get(ctx, &tr, fball.EP_Timezone, nil); err != nil {
		return nil, err
	}
	return []fball.TimezoneResponse{tr}, nil
}

type CountryParams struct {
	Name   string
	Code   string
	Search string
}

func (cp CountryParams) URLQueryString() string {
	strs := []string{}
	if cp.Name != "" {
		strs = append(strs, "name="+template.URLQueryEscaper(cp.Name))
	}
	if cp.Code != "" {
		strs = append(strs, "code="+template.URLQueryEscaper(cp.Code))
	}
	if cp.Search != "" {
		strs = append(strs, "search="+template.URLQueryEscaper(cp.Search))
	}
	return strings.Join(strs, "&")
}

func (c *Client) Country(ctx context.Context, p CountryParams) ([]fball.CountryResponse, error) {
	cr := fball.CountryResponse{}
	if err := c.get(ctx, &cr, fball.EP_Countries, p); err != nil {
		return nil, err
	}
	return []fball.CountryResponse{cr}, nil
}

type response interface {
	Err() error
	When(int64)
}

const base = "https://v3.football.api-sports.io"

type urlQueryStringer interface {
	URLQueryString() string
}

func (c *Client) get(ctx context.Context, data response, endpoint string, params urlQueryStringer) error {
	if data == nil {
		return fmt.Errorf("inalid data: must be non-nil")
	}
	if len(endpoint) == 0 {
		return fmt.Errorf("invalid endpoint: empty string")
	}

	url := base + endpoint
	qstr := params.URLQueryString()
	if qstr != "" {
		url += "?"
		url += qstr
	}

	c.logger.Println("GET", url)
	c.limit.Take()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	now := time.Now().UTC().UnixNano()
	req.Header.Set("X-RapidAPI-Key", c.key)
	resp, err := c.doer.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(data); err != nil {
		return err
	}
	data.When(now)

	return data.Err()
}
