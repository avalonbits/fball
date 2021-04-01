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

type CommonResponse struct {
	Get        string      `json:"get"`
	Parameters interface{} `json:"parameters"`
	Errors     interface{} `json:"errors"`
	Results    int         `json:"results"`
	Paging     PagingToken `json:"paging"`
	Timestamp  int64
}

func (cr *CommonResponse) When(timestamp int64) {
	cr.Timestamp = timestamp
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
