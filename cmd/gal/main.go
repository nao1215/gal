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
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
)

type options struct {
	LOC     bool `short:"l" long:"loc" description:"Sort AUTHOR names by amount of modified LOC (descendig order, default: alphabetical order)"`
	Commit  bool `short:"c" long:"commit" description:"Sort AUTHOR names by amount of commits (descendig order, default: alphabetical order)"`
	Version bool `short:"v" long:"version" description:"Show gal command version"`
}

var osExit = os.Exit

const cmdName string = "gal"
const version string = "1.1.1"

type user struct {
	name string
	mail string
}

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
	header := getHeader(opts)
	authors := getAuthors(opts)
	makeAuthorsFile(header, authors)
}

// getAuthors returns authors in this project.
// This method get authors name and mail address from git log.
func getAuthors(opts options) []string {
	if opts.Commit {
		return getAuthorsCommitOrder()
	} else if opts.LOC {
		return getAuthorsLOCOrder()
	}
	return getAuthorsAlphabeticalOrder()
}

// getAuthorsAlphabeticalOrder returs authors name and authors mail address in alphabetical order
func getAuthorsAlphabeticalOrder() []string {
	out, err := exec.Command("git", "log", "--pretty=format:%an<%ae>").Output()
	if err != nil {
		die(err.Error())
	}

	list := strings.Split(string(out), "\n")
	list = removeDuplicate(list)
	sort.Strings(list)
	return list
}

func getAuthorsCommitOrder() []string {
	// https://stackoverflow.com/questions/51966053/what-is-wrong-with-invoking-git-shortlog-from-go-exec
	// > If no revisions are passed on the command line and either standard input is not a terminal or there
	// > is no current branch, git shortlog will output a summary of the log read from standard input, without
	// > reference to the current repository.
	out, err := exec.Command("git", "shortlog", "--numbered", "--summary", "--email", "main").Output()
	if err != nil {
		out, err = exec.Command("git", "shortlog", "--numbered", "--summary", "--email", "master").Output()
		if err != nil {
			die("this repository don't have main branch or master branch: " + err.Error())
		}
	}
	// Before:     4	CHIKAMATSU Naohiro <n.chika156@gmail.com>
	// After :CHIKAMATSU Naohiro <n.chika156@gmail.com>
	var result []string
	re := regexp.MustCompile(`\s+\d+\s`)
	for _, v := range strings.Split(string(out), "\n") {
		if v == "" {
			break
		}
		result = append(result, re.ReplaceAllString(v, ""))
	}
	return result
}

func getAuthorsLOCOrder() []string {
	// key=user, value=amount of modified LOC
	userInfo := map[user]int{}

	users := getUsers()
	for _, v := range users {
		out, err := exec.Command("git", "log", "--author="+v.mail, "--numstat", "--pretty=", "--no-merges", "main").Output()
		if err != nil {
			out, err = exec.Command("git", "log", "--author=\""+v.mail+"\"", "--numstat", "--pretty=", "--no-merges", "master").Output()
			if err != nil {
				die("this repository don't have main branch or master branch: " + err.Error())
			}
		}

		list := strings.Split(string(out), "\n")
		sum := 0
		for _, line := range list {
			list := strings.Fields(line)
			if len(list) == 3 { // 0=append line num, 1=delete line num, 2=file name
				sum = sum + atoi(list[0]) + atoi(list[1])
			}
		}
		userInfo[v] = sum
	}

	type kv struct {
		Key   user
		Value int
	}

	var ss []kv
	for k, v := range userInfo {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	result := []string{}
	for _, kv := range ss {
		result = append(result, kv.Key.name+"<"+kv.Key.mail+">")
	}
	return result
}

func sortUserInfo(m map[user]int) map[user]int {

	return m
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		die("can not get amount of modified LOC")
	}
	return i
}

func getUsers() []user {
	users := []user{}
	authors := getAuthors(options{})

	rex := regexp.MustCompile(`<[^<]*@.*>$`) // e-mail address
	for _, v := range authors {
		mailWithAngleBrackets := rex.FindString(v)
		tmp := strings.Replace(mailWithAngleBrackets, "<", "", 1)
		mail := strings.Replace(tmp, ">", "", 1)

		u := user{
			name: strings.Replace(v, mailWithAngleBrackets, "", 1),
			mail: mail,
		}
		users = append(users, u)
	}
	return users
}

func getHeader(opts options) string {
	if opts.Commit {
		return "# Authors List (in descending order, amount of commits)\n"
	} else if opts.LOC {
		return "# Authors List (in descendig order, amount of modified LOC)\n"
	}
	return "# Authors List (in alphabetical order)\n"
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
func makeAuthorsFile(header string, text []string) {
	file, err := os.Create("AUTHORS.md")
	if err != nil {
		die(err.Error())
	}
	defer file.Close()

	file.WriteString(header)
	for _, v := range text {
		if !strings.Contains(v, "dependabot") {
			file.WriteString(v + "\n")
		}
	}
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

	if opts.Commit && opts.LOC {
		die("can not specify --commit option and --loc at same time")
	}
	return args
}

// showHelp print help messages.
func showHelp(p *flags.Parser) {
	p.WriteHelp(os.Stdout)
}

// showHelpFooter print author contact information.
func showHelpFooter() {
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
	description := cmdName + " version " + version + " (under Apache License version 2.0)"
	fmt.Fprintln(os.Stdout, description)
}
