package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"log"
	"os"
	"strings"
)

const version = "1.2.3"

func selfUpdate(slug string) error {
	selfupdate.EnableLog()

	previous := semver.MustParse(version)
	latest, err := selfupdate.UpdateSelf(previous, slug)
	if err != nil {
		return err
	}

	if previous.Equals(latest.Version) {
		fmt.Println("Current binary is the latest version", version)
	} else {
		fmt.Println("Update successfully done to version", latest.Version)
		fmt.Println("Release note:\n", latest.ReleaseNotes)
	}
	return nil
}

func usage() {
	_, _ = fmt.Fprintln(os.Stderr, "Usage: selfupdate-example [flags]\n")
	flag.PrintDefaults()
}
func confirmAndSelfUpdate(slug string) {
	latest, found, err := selfupdate.DetectLatest(slug)
	if err != nil {
		log.Println("Error occurred while detecting version:", err)
		return
	}

	v := semver.MustParse(version)
	if !found || latest.Version.LTE(v) {
		log.Println("Current version is the latest")
		return
	}

	fmt.Print("Do you want to update to", latest.Version, "? (y/n): ")
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')

	s := strings.TrimSpace (input)

	if err != nil || (s != "y" && s != "n") {
		log.Println("Invalid input")
		return
	}
	if input == "n\n" {
		return
	}

	exe, err := os.Executable()
	if err != nil {
		log.Println("Could not locate executable path")
		return
	}
	selfupdate.EnableLog()

	if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
		log.Println("Error occurred while updating binary:", err)
		return
	}
	log.Println("Successfully updated to version", latest.Version)
}
func main(){
	ver := flag.Bool("version", false, "Show version")
	slug := flag.String("slug", "rhysd/go-github-selfupdate", "Repository of this command")

	flag.Usage = usage
	flag.Parse()

	if *ver {
		fmt.Println(version)
		os.Exit(0)
	}

	confirmAndSelfUpdate(*slug)
	os.Exit(0)


	//usage()
}