//
// Copyright 2017 Bryan T. Meyers <bmeyers@datadrake.com>
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

package cmd

import (
	"flag"
	"github.com/DataDrake/cuppa/providers"
	"github.com/DataDrake/cuppa/results"
	"github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/format"
	"github.com/DataDrake/waterlog/level"
	"log"
	"os"
)

// Releases fulfills the "release" subcommand
type Releases struct{}

// Flags builds a flagset for this command
func (r Releases) Flags() *flag.FlagSet {
	rcmd := flag.NewFlagSet("releases", flag.ExitOnError)
	rcmd.Usage = r.Usage
	return rcmd
}

// Short prints a quick description of this command
func (r Releases) Short() string {
	return "Get all stable releases"
}

// Usage prints the usage for this command
func (r Releases) Usage() {
	print("USAGE: " + os.Args[0] + " releases <URL>\n\n")
	print("DESCRIPTION: " + r.Short() + "\n\n")
	r.Flags().PrintDefaults()
}

/*
Execute releases for all providers
*/
func (r Releases) Execute() int {
	ps := providers.All()
	flags := r.Flags()
	flags.Parse(os.Args[2:])

	if flags.NArg() != 1 {
		return USAGE
	}

	w := waterlog.New(os.Stdout, "", log.Ltime)
	w.SetLevel(level.Info)
	w.SetFormat(format.Min)
	found := false
	for _, p := range ps {
		w.Infof("\033[1m%s\033[21m checking for match:\n", p.Name())
		name := p.Match(flags.Arg(0))
		if name == "" {
			w.Warnf("\033[1m%s\033[21m does not match.\n", p.Name())
			continue
		}
		rs, s := p.Releases(name)
		if s != results.OK {
			w.Warnf("Failed to fetch releases, code: %d\n", s)
			continue
		}
		found = true
		rs.PrintAll()
		w.Goodf("\033[1m%s\033[21m match(es) found.\n", p.Name())
	}
	if found {
		w.Goodln("Done")
	} else {
		w.Fatalln("No releases found.")
	}
	return GOOD
}
