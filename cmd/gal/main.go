//
// gal/cmd/gal/main.go
//
// Copyright 2022 Naohiro CHIKAMATSU
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/jessevdk/go-flags"
)

type options struct {
	Version bool `short:"v" long:"version" description:"Show gal command version"`
}

var osExit = os.Exit

const cmdName string = "gal"
const version string = "1.0.0"

const (
	exitSuccess int = iota // 0
	exitFailure
)

func main() {
	var opts options
	parseArgs(&opts)
	canUseGitCommand()
	if !exists(".git") {
		die(": .git directory does not exists")
	}
	makeAuthorsFile(getAuthors())
}

// getAuthors returns authors in this project.
// This method get authors name and mail address from git log.
func getAuthors() []string {
	out, err := exec.Command("git", "log", "--pretty=format:%an<%ae>").Output()
	if err != nil {
		die("")
	}

	list := strings.Split(string(out), "\n")
	list = removeDuplicate(list)
	sort.Strings(list)
	return list
}

// removeDuplicate removes duplicates in the slice.
func removeDuplicate(list []string) []string {
	results := make([]string, 0, len(list))
	encountered := map[string]bool{}
	for i := 0; i < len(list); i++ {
		if !encountered[list[i]] {
			encountered[list[i]] = true
			results = append(results, list[i])
		}
	}
	return results
}

// canUseGitCommand check whether git command install in the system.
// If not install, exit command.
func canUseGitCommand() {
	_, err := exec.LookPath("git")
	if err != nil {
		die("git command does not exist. You need git command to execute gal command")
	}
}

// exists check whether file or directory exists.
func exists(path string) bool {
	_, err := os.Stat(path)
	return (err == nil)
}

// die exit program with message.
func die(msg string) {
	fmt.Fprintln(os.Stderr, cmdName+": "+msg)
	osExit(exitFailure)
}

// If it can not create file, exit command.
func makeAuthorsFile(text []string) {
	file, err := os.Create("AUTHORS.md")
	if err != nil {
		die(err.Error())
	}
	defer file.Close()

	file.WriteString("# Authors List (in alphabetical order)\n")
	file.WriteString(strings.Join(text, "\n"))
	file.WriteString("\n")
}

// parseArgs parse command line arguments.
// In this method, process for version option, help option.
func parseArgs(opts *options) []string {
	p := newParser(opts)

	args, err := p.Parse()
	if err != nil {
		// If user specify --help option, help message already showed in p.Parse().
		// Moreover, help messages are stored inside err.
		if !strings.Contains(err.Error(), cmdName) {
			showHelp(p)
			showHelpFooter()
		} else {
			showHelpFooter()
		}
		osExit(exitFailure)
	}

	if opts.Version {
		showVersion(cmdName, version)
		osExit(exitSuccess)
	}
	return args
}

// showHelp print help messages.
func showHelp(p *flags.Parser) {
	p.WriteHelp(os.Stdout)
}

// showHelpFooter print author contact information.
func showHelpFooter() {
	fmt.Println("")
	fmt.Println("Contact:")
	fmt.Println("  If you find the bugs, please report the content of the error.")
	fmt.Println("  [GitHub Issue] https://github.com/nao1215/gal/issues")
}

// newParser return initialized flags.Parser.
func newParser(opts *options) *flags.Parser {
	parser := flags.NewParser(opts, flags.Default)
	parser.Name = cmdName
	parser.Usage = "[OPTIONS]"
	return parser
}

// showVersion show ubume command version information.
func showVersion(cmdName string, version string) {
	description := cmdName + " version " + version + " (under Apache License verison 2.0)"
	fmt.Fprintln(os.Stdout, description)
}
