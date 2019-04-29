package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	flag "github.com/spf13/pflag"
)

// Module represents a single Go module
type Module struct {
	Path     string       // module path
	Version  string       // module version
	Versions []string     // available module versions (with -versions)
	Replace  *Module      // replaced by this module
	Time     *time.Time   // time version was created
	Update   *Module      // available update, if any (with -u)
	Main     bool         // is this the main module?
	Indirect bool         // is this module only an indirect dependency of main module?
	Dir      string       // directory holding files for this module, if any
	GoMod    string       // path to go.mod file for this module, if any
	Error    *ModuleError // error loading module
}

// ModuleError contains the error message that occurred when loading the module
type ModuleError struct {
	Err string
}

// Flags represents all the flags accepted by the CLI
type Flags struct {
	CheckOldPkgs  bool
	CheckIndirect bool
	IgnoredPkgs   []string
}

// Results contains all the modules that need to be reported
type Results struct {
	Updated  []*Module
	Replaced []*Module
	Old      []*Module
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

	// get an invalid JSON list of all modules
	out, err := Run("go", "list", "-m", "-u", "-json", "all")
	if err != nil {
		exitStatus = 1
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	// make list a valid JSON list
	out = "[" + out + "]"
	out = strings.ReplaceAll(out, "}\n{", "},\n{")

	// Parse the JSON list into or Go Slice
	modules := []*Module{}
	err = json.Unmarshal([]byte(out), &modules)
	if err != nil {
		exitStatus = 1
		fmt.Fprintln(os.Stderr, err.Error())
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
func checkModule(f *Flags, m *Module, r *Results) {
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
