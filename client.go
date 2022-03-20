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
	"log"
	"net/http"
	"time"
)

// Doer is an interface for perfomring http requests.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type limiter interface {
	Take() time.Time
}

// Client is an api-football.com client.
type Client struct {
	key    string
	doer   Doer
	limit  limiter
	logger *log.Logger
}

// NewClient creates an api-football.com client. The key is the one provided by the
// service when you register it. limit will rate limit any request and if no logger
// is provided, a default logger is used.
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

// Response is an interface for api-football.com responses.
type Response interface {
	// Err returns the error from the response, if any.
	Err() error

	// When returns the timestamp for the response.
	When() int64

	// setWhen sets the timestamp for the response.
	setWhen(int64)
}

const base = "https://v3.football.api-sports.io"

// Get will perform a GET request against the api-football service.
// The response is returned in the data out param.
func (c *Client) Get(ctx context.Context, endpoint string, queryStr string, data Response) error {
	if data == nil {
		return fmt.Errorf("inalid data: must be non-nil")
	}
	if len(endpoint) == 0 {
		return fmt.Errorf("invalid endpoint: empty string")
	}

	url := base + endpoint
	if queryStr != "" {
		url += "?"
		url += queryStr
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
		c.logger.Println(err)
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(data); err != nil {
		return err
	}
	data.setWhen(now)

	return data.Err()
}
