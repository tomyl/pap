package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"sort"

	"github.com/auroralaboratories/pulse"
	xdg "github.com/queria/golang-go-xdg"
)

const (
	clientname = "pap"
	confdir    = "pap"
	confname   = "profiles.json"
)

var flagAdd string
var flagList bool
var flagListAuto bool
var flagNext bool
var flagNextAuto bool
var flagRemove string
var flagNotifications bool
var flagVerbose bool

func init() {
	flag.StringVar(&flagAdd, "add", "", "Save current source/sink pair as this profile name.")
	flag.BoolVar(&flagList, "list", false, "List profiles.")
	flag.BoolVar(&flagListAuto, "list-auto", false, "List auto-generated profiles.")
	flag.BoolVar(&flagNext, "next", false, "Switch to next profile.")
	flag.BoolVar(&flagNextAuto, "next-auto", false, "Switch to next auto-generated profile.")
	flag.StringVar(&flagRemove, "remove", "", "Remove profile.")
	flag.BoolVar(&flagNotifications, "notify", false, "Use desktop notifications.")
	flag.BoolVar(&flagVerbose, "verbose", false, "Use verbose output.")
}

func success(pattern string, params ...interface{}) {
	message(0, pattern, params...)
}

func failure(pattern string, params ...interface{}) {
	message(1, pattern, params...)
}

func message(rc int, pattern string, params ...interface{}) {
	msg := fmt.Sprintf(pattern, params...)
	if flagNotifications {
		output, err := exec.Command("notify-send", msg).CombinedOutput()
		if err != nil {
			if len(output) > 0 {
				fmt.Print(output)
			}
			log.Fatal(err)
		}
	} else {
		fmt.Println(msg)
	}
	os.Exit(rc)
}

func verbose(pattern string, params ...interface{}) {
	if flagVerbose {
		msg := fmt.Sprintf(pattern, params...)
		fmt.Println(msg)
	}
}

type profile struct {
	Title  string
	Source *pulse.Source
	Sink   *pulse.Sink
}

func getClient() *pulse.Client {
	client, err := pulse.NewClient(clientname)

	if err != nil {
		failure("Failed to get client: %v", err)
	}

	return client
}

func getServerInfo(client *pulse.Client) pulse.ServerInfo {
	info, err := client.GetServerInfo()

	if err != nil {
		failure("Failed to get server info: %v", err)
	}

	return info
}

func getSources(client *pulse.Client) []pulse.Source {
	sources, err := client.GetSources()

	if err != nil {
		failure("Failed to get sources: %v", err)
	}

	return sources
}

func getSourceByName(sources []pulse.Source, name string, must bool) *pulse.Source {
	for _, source := range sources {
		if source.Name == name {
			return &source
		}
	}
	if must {
		failure("Failed to find source named %s", name)
	} else if flagVerbose {
		verbose("Failed to find source named %s", name)
	}
	return nil
}

func getSinkByName(sinks []pulse.Sink, name string, must bool) *pulse.Sink {
	for _, sink := range sinks {
		if sink.Name == name {
			return &sink
		}
	}
	if must {
		failure("Failed to find sink named %s", name)
	} else if flagVerbose {
		verbose("Failed to find sink named %s", name)
	}
	return nil
}

func getSinks(client *pulse.Client) []pulse.Sink {
	sinks, err := client.GetSinks()

	if err != nil {
		failure("Failed to get sinks: %v", err)
	}

	return sinks
}

func getProfilesPath() string {
	filepath, err := xdg.Data.Ensure(path.Join(confdir, confname))

	if err != nil {
		failure("Failed to get path: %v", err)
	}

	return filepath
}

func loadProfiles() []profile {
	profiles := make([]profile, 0)
	filepath := getProfilesPath()

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return profiles
	}

	verbose("Loading %s", filepath)
	buf, err := ioutil.ReadFile(filepath)

	if err != nil {
		failure("Failed to read %s: %v", filepath, err)
	}

	if err := json.Unmarshal(buf, &profiles); err != nil {
		failure("Failed to decode json: %v", err)
	}

	sort.Sort(byTitle(profiles))

	return profiles
}

func saveProfiles(profiles []profile) {
	filepath := getProfilesPath()
	verbose("Saving %s", filepath)

	fh, err := os.Create(filepath)

	if err != nil {
		failure("Failed to open %s: %v", filepath, err)
	}

	defer fh.Close()

	buf, err := json.Marshal(profiles)

	if err != nil {
		failure("Failed to encode json: %v", err)
	}

	if _, err := fh.Write(buf); err != nil {
		failure("Failed to write profiles: %v", err)
	}
}

func cmdListProfiles() {
	profiles := loadProfiles()
	client := getClient()
	info := getServerInfo(client)
	listProfiles(profiles, info)
}

func listProfiles(profiles []profile, info pulse.ServerInfo) {
	for _, profile := range profiles {
		def := ""
		if (profile.Source == nil || profile.Source.Name == info.DefaultSourceName) && profile.Sink.Name == info.DefaultSinkName {
			def = "*"
		}
		fmt.Printf("%s%s\n", profile.Title, def)
		if flagVerbose {
			if profile.Source != nil {
				fmt.Printf("\tsource %s (%s)\n", profile.Source.Name, profile.Source.Description)
			}
			fmt.Printf("\tsink   %s (%s)\n", profile.Sink.Name, profile.Sink.Description)
		}
	}
}

func cmdListAutoProfiles() {
	client := getClient()
	info := getServerInfo(client)
	profiles := getAutoProfiles(client)
	listProfiles(profiles, info)
}

func getAutoProfiles(client *pulse.Client) []profile {
	sources := getSources(client)
	sinks := getSinks(client)
	cardsource := make(map[int]pulse.Source, 0)

	for _, source := range sources {
		cardsource[source.CardIndex] = source
	}

	profiles := make([]profile, 0)

	for _, sink := range sinks {
		source := cardsource[sink.CardIndex]
		sink := sink
		profiles = append(profiles, profile{sink.Description, &source, &sink})
	}

	return profiles
}

func cmdAddProfile(title string) {
	client := getClient()
	info := getServerInfo(client)

	if info.DefaultSourceName == "" {
		failure("No default source!")
	}

	if info.DefaultSinkName == "" {
		failure("no default sink")
	}

	profiles := loadProfiles()

	for _, profile := range profiles {
		if profile.Source.Name == info.DefaultSourceName && profile.Sink.Name == info.DefaultSinkName {
			failure("Already added as profile \"%s\".", profile.Title)
		}
	}

	sources := getSources(client)
	sinks := getSinks(client)

	source := getSourceByName(sources, info.DefaultSourceName, true)
	sink := getSinkByName(sinks, info.DefaultSinkName, true)

	profiles = append(profiles, profile{title, source, sink})
	saveProfiles(profiles)

	success("Added profile %s.", title)
}

func cmdRemoveProfile(title string) {
	oldProfiles := loadProfiles()
	newProfiles := make([]profile, 0)

	for _, profile := range oldProfiles {
		if profile.Title != title {
			newProfiles = append(newProfiles, profile)
		}
	}

	if len(oldProfiles) == len(newProfiles) {
		failure("Found no profile named %s.", title)
	}

	saveProfiles(newProfiles)
	success("Removed profile %s.", title)
}

func cmdNextProfile() {
	profiles := loadProfiles()
	client := getClient()
	nextProfile(client, profiles)
}

func cmdNextAutoProfile() {
	client := getClient()
	profiles := getAutoProfiles(client)
	nextProfile(client, profiles)
}

func nextProfile(client *pulse.Client, profiles []profile) {
	if len(profiles) == 0 {
		failure("No profiles!")
	}

	activeidx := len(profiles) - 1
	info := getServerInfo(client)

	for i, profile := range profiles {
		if (profile.Source == nil || profile.Source.Name == info.DefaultSourceName) && profile.Sink.Name == info.DefaultSinkName {
			activeidx = i
		}
	}

	startidx := activeidx
	sources := getSources(client)
	sinks := getSinks(client)

	for {
		activeidx++

		if activeidx >= len(profiles) {
			activeidx = 0
		}

		if activeidx == startidx {
			failure("Found no usable profile")
		}

		active := profiles[activeidx]

		if active.Source != nil {
			source := getSourceByName(sources, active.Source.Name, false)

			if source == nil {
				continue
			}

			if err := client.SetDefaultSource(source.Name); err != nil {
				failure("Failed to set default source %s: %v", source.Name, err)
			}

			if source.Muted {
				if err := source.Unmute(); err != nil {
					failure("Failed to unmute source %s: %v", source, err)
				}
			}
		}

		if active.Sink != nil {
			sink := getSinkByName(sinks, active.Sink.Name, false)

			if sink == nil {
				continue
			}

			if err := client.SetDefaultSink(sink.Name); err != nil {
				failure("Failed to set default sink %s: %v", sink.Name, err)
			}
		}

		success("Activated profile %s.", active.Title)
		return
	}
}

func main() {
	flag.Parse()

	switch {
	case flagAdd != "":
		cmdAddProfile(flagAdd)
	case flagList:
		cmdListProfiles()
	case flagListAuto:
		cmdListAutoProfiles()
	case flagNext:
		cmdNextProfile()
	case flagNextAuto:
		cmdNextAutoProfile()
	case flagRemove != "":
		cmdRemoveProfile(flagRemove)
	default:
		fmt.Printf("pap - a simple pulseaudio profile manager\n")
		flag.PrintDefaults()
	}
}

type byTitle []profile

func (a byTitle) Len() int           { return len(a) }
func (a byTitle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byTitle) Less(i, j int) bool { return a[i].Title < a[j].Title }
