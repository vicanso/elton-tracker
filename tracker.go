// Copyright 2018 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tracker

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/vicanso/elton"
)

const (
	// HandleSuccess handle success
	HandleSuccess = iota
	// HandleFail handle fail
	HandleFail
)

var (
	defaultMaskFields  = regexp.MustCompile(`password`)
	errNoTrackFunction = errors.New("require on track function")
)

type (
	// Info tracker info
	Info struct {
		CID    string                 `json:"cid,omitempty"`
		Query  map[string]string      `json:"query,omitempty"`
		Params map[string]string      `json:"params,omitempty"`
		Form   map[string]interface{} `json:"form,omitempty"`
		Result int                    `json:"result,omitempty"`
		Err    error                  `json:"err,omitempty"`
	}
	// OnTrack on track function
	OnTrack func(*Info, *elton.Context)
	// Config tracker config
	Config struct {
		OnTrack OnTrack
		Mask    *regexp.Regexp
		Skipper elton.Skipper
	}
)

func convertMap(data map[string]string, mask *regexp.Regexp) map[string]string {
	if len(data) == 0 {
		return nil
	}
	m := make(map[string]string)
	for k, v := range data {
		if mask.MatchString(k) {
			m[k] = "***"
		} else {
			m[k] = v
		}
	}
	return m
}

// New create a tracker middleware
func New(config Config) elton.Handler {
	mask := config.Mask
	if mask == nil {
		mask = defaultMaskFields
	}
	if config.OnTrack == nil {
		panic(errNoTrackFunction)
	}
	skipper := config.Skipper
	if skipper == nil {
		skipper = elton.DefaultSkipper
	}
	return func(c *elton.Context) (err error) {
		if skipper(c) {
			return c.Next()
		}
		result := HandleSuccess
		query := convertMap(c.Query(), mask)
		params := convertMap(c.Params, mask)
		var form map[string]interface{}
		if len(c.RequestBody) != 0 {
			form = make(map[string]interface{})
			json.Unmarshal(c.RequestBody, &form)
			for k := range form {
				if mask.MatchString(k) {
					form[k] = "***"
				}
			}
		}
		err = c.Next()
		if err != nil {
			result = HandleFail
		}
		config.OnTrack(&Info{
			CID:    c.ID,
			Query:  query,
			Params: params,
			Form:   form,
			Result: result,
			Err:    err,
		}, c)
		return
	}
}
