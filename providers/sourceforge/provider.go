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

package sourceforge

import (
	"encoding/xml"
	"fmt"
	"github.com/DataDrake/cuppa/results"
	"net/http"
	"os"
	"regexp"
	"time"
)

const (
	// API is the format string for a sourceforge RSS feed
	API = "https://sourceforge.net/projects/%s/rss?path=/%s"
)

// TarballRegex matches SourceForge sources
var TarballRegex = regexp.MustCompile("https?://.*sourceforge.net/projects?/(.+)/files/(.+/)?(.+?)-([\\d]+(?:.\\d+)*\\w*?)\\.(?:zip|tar\\..+z.*)(?:\\/download)?$")

// ProjectRegex matches SourceForge sources
var ProjectRegex = regexp.MustCompile("https?://.*sourceforge.net/projects?/(.+)/(?:files/)?(.+?/)?(.+?)-([\\d]+(?:.\\d+)*\\w*?).+$")

// Provider is the upstream provider interface for SourceForge
type Provider struct{}

// Item represents an entry in the RSS Feed
type Item struct {
	XMLName xml.Name `xml:"item"`
	Link    string   `xml:"link"`
	Date    string   `xml:"pubDate"`
}

// Feed represents the RSS feed itself
type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Items   []Item   `xml:"channel>item"`
}

// toResults converts a Feed to a ResultSet
func (f *Feed) toResults(name string) *results.ResultSet {
	rs := results.NewResultSet(name)
	for _, item := range f.Items {
		sm := TarballRegex.FindStringSubmatch(item.Link)
		if len(sm) != 5 {
			continue
		}
		pub, _ := time.Parse(time.RFC1123, item.Date+"C")
		r := results.NewResult(name, sm[4], item.Link, pub)
		rs.AddResult(r)
	}
	return rs
}

// Latest finds the newest release for a SourceForge package
func (c Provider) Latest(name string) (r *results.Result, s results.Status) {
	rs, s := c.Releases(name)
	if s != results.OK {
		return
	}
	r = rs.First()
	return
}

// Match checks to see if this provider can handle this kind of query
func (c Provider) Match(query string) string {
	sm := TarballRegex.FindStringSubmatch(query)
	if len(sm) != 5 {
		sm = ProjectRegex.FindStringSubmatch(query)
		if len(sm) != 5 {
			return ""
		}
	}
	return sm[0]
}

// Name gives the name of this provider
func (c Provider) Name() string {
	return "SourceForge"
}

// Releases finds all matching releases for a SourceForge package
func (c Provider) Releases(name string) (rs *results.ResultSet, s results.Status) {
	sm := TarballRegex.FindStringSubmatch(name)
	if len(sm) != 5 {
		sm = ProjectRegex.FindStringSubmatch(name)
		sm[1], sm[3] = sm[3], sm[1]
	}
	// Query the API
	resp, err := http.Get(fmt.Sprintf(API, sm[1], sm[2]))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		s = results.Unavailable
		return
	}
	defer resp.Body.Close()
	// Translate Status Code
	switch resp.StatusCode {
	case 200:
		s = results.OK
	case 404:
		s = results.NotFound
	default:
		s = results.Unavailable
	}

	// Fail if not OK
	if s != results.OK {
		return
	}

	dec := xml.NewDecoder(resp.Body)
	feed := &Feed{}
	err = dec.Decode(feed)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		s = results.Unavailable
		return
	}
	rs = feed.toResults(sm[3])
	if rs.Len() == 0 {
		s = results.NotFound
	}
	return
}
