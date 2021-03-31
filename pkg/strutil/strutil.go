/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package strutil

import (
	"encoding/csv"
	"strings"

	"github.com/containerd/containerd/errdefs"
	"github.com/pkg/errors"
)

// ConvertKVStringsToMap is from https://github.com/moby/moby/blob/v20.10.0-rc2/runconfig/opts/parse.go
//
// ConvertKVStringsToMap converts ["key=value"] to {"key":"value"}
func ConvertKVStringsToMap(values []string) map[string]string {
	result := make(map[string]string, len(values))
	for _, value := range values {
		kv := strings.SplitN(value, "=", 2)
		if len(kv) == 1 {
			result[kv[0]] = ""
		} else {
			result[kv[0]] = kv[1]
		}
	}

	return result
}

// InStringSlice checks whether a string is inside a string slice.
// Comparison is case insensitive.
//
// From https://github.com/containerd/containerd/blob/7c6d710bcfc81a30ac1e8cbb2e6a4c294184f7b7/pkg/cri/util/strings.go#L21-L30
func InStringSlice(ss []string, str string) bool {
	for _, s := range ss {
		if strings.EqualFold(s, str) {
			return true
		}
	}
	return false
}

func DedupeStrSlice(in []string) []string {
	m := make(map[string]struct{})
	var res []string
	for _, s := range in {
		if _, ok := m[s]; !ok {
			res = append(res, s)
			m[s] = struct{}{}
		}
	}
	return res
}

// ParseCSVMap parses a string like "foo=x,bar=y" into a map
func ParseCSVMap(s string) (map[string]string, error) {
	csvR := csv.NewReader(strings.NewReader(s))
	ra, err := csvR.ReadAll()
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse %q", s)
	}
	if len(ra) != 1 {
		return nil, errors.Wrapf(errdefs.ErrInvalidArgument, "expected a single line, got %d lines", len(ra))
	}
	fields := ra[0]
	m := make(map[string]string)
	for _, field := range fields {
		kv := strings.SplitN(field, "=", 2)
		if len(kv) == 2 {
			m[kv[0]] = kv[1]
		} else {
			m[kv[0]] = ""
		}
	}
	return m, nil
}
