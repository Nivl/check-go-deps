package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Nivl/check-deps/modutil"
	flag "github.com/spf13/pflag"
)

// Flags represents all the flags accepted by the CLI
type Flags struct {
	CheckOldPkgs  bool
	CheckIndirect bool
	IgnoredPkgs   []string
}

// https://www.gnu.org/software/libc/manual/html_node/Exit-Status.html
const (
	ExitSuccess = 0
	ExitFailure = 1
)

func main() {
	os.Exit(run(os.Args, os.Stderr))
}

func run(args []string, out io.Writer) (exitStatus int) {
	flags, err := parseFlags(args)
	if err != nil {
		fmt.Fprintln(out, fmt.Sprintf("could not parse the flags: %s", err.Error()))
		return ExitFailure
	}

	modules, err := modutil.ParseCwd()
	if err != nil {
		fmt.Fprintln(out, fmt.Sprintf("could not parse the go.mod file: %s", err.Error()))
		return ExitFailure
	}

	res := parseModules(flags, modules)
	res.print(out)

	if res.HasModules() {
		return ExitFailure
	}
	return ExitSuccess
}

func parseFlags(args []string) (*Flags, error) {
	flags := &Flags{}
	fset := &flag.FlagSet{}
	fset.BoolVar(&flags.CheckOldPkgs, "old", false, "check for modules without updates for the last 6 months")
	fset.BoolVar(&flags.CheckIndirect, "indirect", false, "check indirect modules")
	fset.StringSliceVarP(&flags.IgnoredPkgs, "ignore", "i", []string{}, "coma separated list of packages to ignore")
	return flags, fset.Parse(args)
}

func parseModules(f *Flags, modules []*modutil.Module) *Results {
	res := &Results{}
	for _, m := range modules {
		// skip ignored packages
		isIgnored := false
		for _, pkg := range f.IgnoredPkgs {
			if strings.HasPrefix(m.Path, pkg) {
				isIgnored = true
				break
			}
		}
		if isIgnored {
			continue
		}

		// skip indirects if we don't want them
		if m.Indirect && !f.CheckIndirect {
			continue
		}

		// Report if the package has been replaced
		if m.Replace != nil {
			res.Replaced = append(res.Replaced, m)
			continue
		}

		// Report if the package has an update available
		if m.Update != nil {
			// It's possible that a tag appears as an "update" from a commit, even
			// if that tag is older
			if m.Time == nil || m.Update.Time == nil || m.Time.Before(*m.Update.Time) {
				res.Updated = append(res.Updated, m)
				continue
			}
		}

		// Report if the package hasn't been updated in 6 months
		if f.CheckOldPkgs && m.Time != nil {
			sixMonths := 6 * 30 * 24 * time.Hour
			if time.Since(*m.Time) >= sixMonths {
				res.Old = append(res.Old, m)
				continue
			}
		}
	}

	return res
}
