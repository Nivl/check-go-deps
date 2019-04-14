package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
)

// Module represent a single Go module
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

func main() {
	checkOldPkgs := flag.Bool("check-old", false, "check for modules without updates for the last 6 months")
	checkIndirect := flag.Bool("check-indirect", false, "check indirect modules")
	ignoredPkgs := flag.StringSliceP("ignore", "i", []string{}, "coma separated list of packages to ignore")
	flag.Parse()

	out, err := Run("go", "list", "-m", "-u", "-json", "all")
	if err != nil {
		log.Fatal(err)
	}

	// make the output valid JSON
	out = "[" + out + "]"
	out = strings.ReplaceAll(out, "}\n{", "},\n{")

	modules := []*Module{}
	err = json.Unmarshal([]byte(out), &modules)
	if err != nil {
		log.Fatal(err)
	}

	for _, m := range modules {
		skip := false
		for _, pkg := range *ignoredPkgs {
			if strings.HasPrefix(m.Path, pkg) {
				skip = true
			}
		}
		if skip {
			continue
		}

		tag := ""
		if m.Indirect {
			if !*checkIndirect {
				continue
			}
			tag = "[indirect] "
		}

		// Report if the package has been replaced
		if m.Replace != nil {
			fmt.Printf(tag+"%s has been replaced by %s\n", m.Path, m.Replace.Path)
			continue
		}

		// Report if the package has an update available
		if m.Update != nil {
			fmt.Printf(tag+"%s can be updated to %s\n", m.Path, m.Update.Version)
			continue
		}

		// Report if the package hasn't been updated in 6 months
		if *checkOldPkgs && m.Time != nil {
			sixMonths := 6 * 30 * 24 * time.Hour
			if time.Since(*m.Time) >= sixMonths {
				fmt.Printf(tag+"%s hasn't been updated in over 6 months (%s)\n", m.Path, m.Time.String())
				continue
			}
		}
	}
}
