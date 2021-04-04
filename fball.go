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

import "time"

type urlQueryStringer interface {
	urlQueryString() string
}

type noParams struct{}

func (np noParams) urlQueryString() string {
	return ""
}

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

type refreshPolicy time.Duration

func (rp refreshPolicy) Valid(now time.Time, tsnano int64) bool {
	return now.UTC().Sub(time.Unix(0, tsnano)) < time.Duration(rp)
}

const (
	rp_OneDay   = refreshPolicy(86400 * time.Second)
	rp_Infinite = refreshPolicy(1<<63 - 1)
)

const (
	ep_Timezone  = "/timezone"
	ep_Countries = "/countries"
)
