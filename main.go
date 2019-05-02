package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Nivl/check-deps/modutil"
	"github.com/olekukonko/tablewriter"
	flag "github.com/spf13/pflag"
)

// Flags represents all the flags accepted by the CLI
type Flags struct {
	CheckOldPkgs  bool
	CheckIndirect bool
	IgnoredPkgs   []string
}

// Results contains all the modules that need to be reported
type Results struct {
	Updated  []*modutil.Module
	Replaced []*modutil.Module
	Old      []*modutil.Module
}

// HasModules checks if the results contains any modules
func (r *Results) HasModules() bool {
	return len(r.Updated) > 0 ||
		len(r.Replaced) > 0 ||
		len(r.Old) > 0
}

func main() {
	exitStatus := 0
	defer func() {
		os.Exit(exitStatus)
	}()

	flags := &Flags{}
	flag.BoolVar(&flags.CheckOldPkgs, "old", false, "check for modules without updates for the last 6 months")
	flag.BoolVar(&flags.CheckIndirect, "indirect", false, "check indirect modules")
	flag.StringSliceVarP(&flags.IgnoredPkgs, "ignore", "i", []string{}, "coma separated list of packages to ignore")
	flag.Parse()

	modules, err := modutil.ParseCwd()
	if err != nil {
		exitStatus = 1
		fmt.Fprintln(os.Stderr, fmt.Sprintf("could not parse the go.mod file: %s", err.Error()))
		return
	}

	// check every modules one-by-one
	res := &Results{}
	for _, m := range modules {
		checkModule(flags, m, res)
	}

	if res.HasModules() {
		exitStatus = 1
	}

	// Print the result
	// Updated
	if len(res.Updated) > 0 {
		table := tablewriter.NewWriter(os.Stderr)
		table.SetHeader([]string{"Module", "Current Version", "New Version", "Indirect"})
		for _, m := range res.Updated {
			table.Append([]string{
				m.Path,
				m.Version,
				m.Update.Version,
				strconv.FormatBool(m.Indirect),
			})
		}
		table.Render()
	}

	// Replaced
	if len(res.Replaced) > 0 {
		table := tablewriter.NewWriter(os.Stderr)
		table.SetHeader([]string{"Module", "Replaced By", "Indirect"})
		for _, m := range res.Replaced {
			table.Append([]string{
				m.Path,
				m.Replace.Path,
				strconv.FormatBool(m.Indirect),
			})
		}
		table.Render()
	}

	// Old
	if len(res.Old) > 0 {
		table := tablewriter.NewWriter(os.Stderr)
		table.SetHeader([]string{"Module", "Last update", "Indirect"})
		for _, m := range res.Old {
			monthsPassed := time.Since(*m.Time) / (24 * time.Hour) / 30
			table.Append([]string{
				m.Path,
				fmt.Sprintf("%d months ago (%s)", monthsPassed, m.Time.Format("2006/01/02")),
				strconv.FormatBool(m.Indirect),
			})
		}
		table.Render()
	}
}

// checkModule checks a single module and prints its status
func checkModule(f *Flags, m *modutil.Module, r *Results) {
	for _, pkg := range f.IgnoredPkgs {
		if strings.HasPrefix(m.Path, pkg) {
			return
		}
	}

	if m.Indirect && !f.CheckIndirect {
		return
	}

	// Report if the package has been replaced
	if m.Replace != nil {
		r.Replaced = append(r.Replaced, m)
		return
	}

	// Report if the package has an update available
	if m.Update != nil {
		// It's possible that a tag appears as an "update" from a commit, even
		// if that tag is older
		if (m.Time != nil && m.Update.Time != nil) && m.Time.After(*m.Update.Time) {
			return
		}
		r.Updated = append(r.Updated, m)
		return
	}

	// Report if the package hasn't been updated in 6 months
	if f.CheckOldPkgs && m.Time != nil {
		sixMonths := 6 * 30 * 24 * time.Hour
		if time.Since(*m.Time) >= sixMonths {
			r.Old = append(r.Old, m)
			return
		}
	}
}
