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
)

func StructToURLQueryString(data interface{}) string {
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
