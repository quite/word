package main

import (
	"fmt"
	"io"
	"log"
	"net/textproto"
	"os"
	"os/exec"
	"strings"

	"github.com/quite/word/config"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/dict"
)

// my shellscript ~/bin/d did this:
// # first eat whitespace-only lines, then http://sed.sourceforge.net/sed1line.txt
//   } | sed "s/^[[:space:]]*$//" | sed -e :a -e '/^\n*$/{$d;N;ba' -e '}'

func hasDict(dicts []dict.Dict, name string) bool {
	for _, dict := range dicts {
		if dict.Name == name {
			return true
		}
	}
	return false
}

var conf *config.Config

func main() {
	var err error

	conf, err = config.New()
	if err != nil {
		panic(err)
	}

	var out io.WriteCloser = os.Stdout
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		// https://stackoverflow.com/a/54198703/945568
		var cmd *exec.Cmd
		cmd, out = runPager()
		defer func() {
			out.Close()
			cmd.Wait()
		}()
	}

	var c *dict.Client
	if c, err = dict.Dial("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port)); err != nil {
		panic(err)
	}
	defer c.Close()

	dicts, err := c.Dicts()
	if err != nil {
		panic(err)
	}

	if len(os.Args) != 2 {
		fmt.Printf("expected 1 arg: word to look up\n")
		fmt.Printf("\ndatabases on %s:%d:\n", conf.Host, conf.Port)
		for _, dict := range dicts {
			fmt.Printf("  %s\n    %s\n", dict.Name, dict.Desc)
		}
		fmt.Printf("\ndatabases configured: %s\n", strings.Join(conf.Databases, " "))
		os.Exit(2)
	}
	word := os.Args[1]

	for _, name := range conf.Databases {
		if hasDict(dicts, name) {
			defs, err := c.Define(name, word)
			if err != nil {
				if err.(*textproto.Error).Code != 552 {
					panic(err)
				}
			}
			for _, def := range defs {
				// TODO custom?
				fmt.Fprintf(out, "\n# %s\n\n%s\n\n", name, strings.TrimSpace(string(def.Text)))
			}
		}
	}
}

func runPager() (*exec.Cmd, io.WriteCloser) {
	pager := []string{"less"}
	if env := os.Getenv("PAGER"); env != "" {
		pager = strings.Split(os.Getenv("PAGER"), " ")
	}
	if conf.Pager != "" {
		pager = strings.Split(conf.Pager, " ")
	}
	cmd := exec.Command(pager[0], pager[1:]...)
	out, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	return cmd, out
}
