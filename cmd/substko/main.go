package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/drsigned/substko/internal/fingerprints"
	"github.com/drsigned/substko/internal/targets"
	"github.com/drsigned/substko/pkg/substko"
	"github.com/logrusorgru/aurora/v3"
)

type options struct {
	concurrency        int
	fingerprintsFile   string
	silent             bool
	targetsListFile    string
	noColor            bool
	updateFingerprints bool
	verbose            bool
}

var (
	co options
	so substko.Options
	au aurora.Aurora
)

func banner() {
	fmt.Fprintln(os.Stderr, aurora.BrightBlue(`
           _         _   _
 ___ _   _| |__  ___| |_| | _____
/ __| | | | '_ \/ __| __| |/ / _ \
\__ \ |_| | |_) \__ \ |_|   < (_) |
|___/\__,_|_.__/|___/\__|_|\_\___/ v1.0.0
`).Bold())
}

func init() {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	flag.StringVar(&co.fingerprintsFile, "f", dir+"/.substko/fingerprints.json", "")
	flag.IntVar(&co.concurrency, "c", 20, "")
	flag.BoolVar(&so.HTTPS, "https", false, "")
	flag.BoolVar(&co.silent, "silent", false, "")
	flag.StringVar(&co.targetsListFile, "l", "", "")
	flag.BoolVar(&co.noColor, "nc", false, "")
	flag.IntVar(&so.Timeout, "timeout", 10, "")
	flag.BoolVar(&co.updateFingerprints, "u", false, "")
	flag.BoolVar(&co.verbose, "v", false, "")

	flag.Usage = func() {
		banner()

		h := "USAGE:\n"
		h += "  substko [OPTIONS]\n"

		h += "\nOPTIONS:\n"
		h += "  -c               concurrency level (default: 20)\n"
		h += "  -f               path to fingerprints file\n"
		h += "  -https           force HTTPS connection (default: false)\n"
		h += "  -l               targets list\n"
		h += "  -nc              no color mode (default: false)\n"
		h += "  -silent          silent mode\n"
		h += "  -timeout         HTTP timeout in seconds (default: 10)\n"
		h += "  -u               download/update fingerprints\n"
		h += "  -v               verbose mode\n"

		fmt.Fprintf(os.Stderr, h)
	}

	flag.Parse()

	au = aurora.NewAurora(!co.noColor)
}

func main() {
	// Update/Download fingerprints
	if co.updateFingerprints {
		success, err := fingerprints.Update(co.fingerprintsFile)
		if err != nil {
			fmt.Println(err)
		}

		if success {
			fmt.Println("downloaded/updated: " + co.fingerprintsFile)
			os.Exit(0)
		}

		os.Exit(1)
	}

	// Load targets
	targets, err := targets.Load(co.targetsListFile)
	if err != nil {
		log.Fatalln(err)
	}

	// Load fingerprints
	fingerprints, err := fingerprints.Load(co.fingerprintsFile)
	if err != nil {
		log.Fatalln(err)
	}

	so.Fingerprints = fingerprints

	if !co.silent {
		banner()

		fmt.Println("")
		fmt.Println("[", au.BrightBlue("INFO"), "]", len(targets), "targets loaded")
		fmt.Println("[", au.BrightBlue("INFO"), "]", len(so.Fingerprints), "fingerprints loaded")
		fmt.Println("")
	}

	targetsChannel := make(chan string, co.concurrency)

	wg := new(sync.WaitGroup)

	for i := 0; i < co.concurrency; i++ {
		wg.Add(1)

		go func() {
			for target := range targetsChannel {
				if target == "" {
					continue
				}

				status, STKOType, at, err := substko.CheckSTKO(target, &so)
				if err != nil && co.verbose {
					fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
					continue
				}

				if status == "Vulnerable" {
					fmt.Println("[", au.BrightGreen(status), "]", target, "-", au.Green(STKOType+":"), au.Italic(au.Green(at)))
				} else if status == "Edge Case" {
					fmt.Println("[", au.BrightYellow(status), "]", target, "-", au.Yellow(STKOType+":"), au.Italic(au.Yellow(at)))
				} else {
					if !co.silent {
						fmt.Println("[", au.BrightRed(status), "]", target)
					}
				}
			}

			wg.Done()
		}()
	}

	for _, target := range targets {
		targetsChannel <- target
	}

	close(targetsChannel)
	wg.Wait()
}
