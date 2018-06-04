//
// Copyright 2016-2018 Bryan T. Meyers <bmeyers@datadrake.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package pypi

import (
	"github.com/DataDrake/cuppa/results"
	"time"
)

// ConvertURLS translates PyPi URLs to Cuppa results
func ConvertURLS(cr []URL, name, version string) *results.Result {
	r := &results.Result{}
	r.Name = name
	r.Version = version
	u := cr[len(cr)-1]
	r.Published, _ = time.Parse(DateFormat, u.UploadTime)
	r.Location = u.URL
	return r
}

// Releases holds one or more Source URLs
type Releases struct {
	Releases map[string][]URL `json:"releases"`
}

// Convert turns PyPi releases into a Cuppa results set
func (crs *Releases) Convert(name string) *results.ResultSet {
	rs := results.NewResultSet(name)
	for ver, rel := range crs.Releases {
		r := ConvertURLS(rel, name, ver)
		if r != nil {
			rs.AddResult(r)
		}
	}
	return rs
}
