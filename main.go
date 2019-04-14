package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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
	CheckOldPkgs   bool
	CheckIndirects bool
	IgnoredPkgs    []string
}

func main() {
	exitStatus := 0
	defer func() {
		os.Exit(exitStatus)
	}()

	flags := &Flags{}
	flag.BoolVar(&flags.CheckOldPkgs, "check-old", false, "check for modules without updates for the last 6 months")
	flag.BoolVar(&flags.CheckIndirects, "check-indirects", false, "check indirect modules")
	flag.StringSliceVarP(&flags.IgnoredPkgs, "ignore", "i", []string{}, "coma separated list of packages to ignore")
	flag.Parse()

	// get an invalid JSON list of all modules
	out, err := Run("go", "list", "-m", "-u", "-json", "all")
	if err != nil {
		exitStatus = 1
		log.Fatal(err)
	}

	// make list a valid JSON list
	out = "[" + out + "]"
	out = strings.ReplaceAll(out, "}\n{", "},\n{")

	// Parse the JSON list into or Go Slice
	modules := []*Module{}
	err = json.Unmarshal([]byte(out), &modules)
	if err != nil {
		exitStatus = 1
		log.Fatal(err)
	}

	// check every modules one-by-one
	for _, m := range modules {
		if checkModule(flags, m) {
			exitStatus = 1
		}
	}
}

// checkModule checks a single module and prints its status
func checkModule(f *Flags, m *Module) (needsUpdate bool) {
	for _, pkg := range f.IgnoredPkgs {
		if strings.HasPrefix(m.Path, pkg) {
			return false
		}
	}

	if m.Indirect && !f.CheckIndirects {
		return false
	}

	// Set the tags
	tag := ""
	if m.Indirect {
		tag += "[indirect]"
	}

	// Report if the package has been replaced
	if m.Replace != nil {
		write(tag, "%s has been replaced by %s", m.Path, m.Replace.Path)
		return true
	}

	// Report if the package has an update available
	if m.Update != nil {
		write(tag, "%s can be updated to %s", m.Path, m.Update.Version)
		return true
	}

	// Report if the package hasn't been updated in 6 months
	if f.CheckOldPkgs && m.Time != nil {
		sixMonths := 6 * 30 * 24 * time.Hour
		if time.Since(*m.Time) >= sixMonths {
			monthsPassed := time.Since(*m.Time) / (24 * time.Hour) / 30
			write(tag, "%s hasn't been updated in over %d months (since %s)", m.Path, monthsPassed, m.Time.Format("2006/01/02"))
			return true
		}
	}
	return false
}

func write(tag, format string, a ...interface{}) {
	if tag != "" {
		tag += " "
	}
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Printf(tag+format, a...)
}
