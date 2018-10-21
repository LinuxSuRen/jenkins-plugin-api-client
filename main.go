package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"log"
)

type Dependency struct {
	Name     string `json:"name"`
	Implied  bool   `json:"implied"`
	Optional bool   `json:"optional"`
	Title    string `json:"title"`
	Version  string `json:"version"`
}

type Plugin struct {
	BuildDate         string       `json:"buildDate"`
	Dependencies      []Dependency `json:"dependencies"`
	Excerpt           string       `json:"excerpt"`
	FirstRelease      string       `json:"firstRelease"`
	Gav               string       `json:"gav"`
	Name              string       `json:"name"`
	PreviousTimestamp string       `json:"previousTimestamp"`
	PreviousVersion   string       `json:"previousVersion"`
	ReleaseTimestamp  string       `json:"releaseTimestamp"`
	RequireCore       string       `json:"RequireCore"`
	Title             string       `json:"title"`
	Url               string       `json:"url"`
	Version           string       `json:"version"`
}

func getPlugin(name string) (plugin *Plugin, err error) {
	resp, err := http.Get("https://plugins.jenkins.io/api/plugin/" + name)
	if err != nil {
		return plugin, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	plugin = &Plugin{}
	err = json.Unmarshal(body, plugin)
	if err != nil {
		log.Println("error when unmarshal:", string(body))
	}
	return
}

func collectDependencies(pluginName string, dependencyMap map[string]string) {
	plugin, err := getPlugin(pluginName)
	if err != nil {
		log.Println("can't get the plugin by name:", pluginName)
		panic(err)
	}

	dependencyMap[plugin.Name] = plugin.Version

	for _, dependent := range plugin.Dependencies {
		if _, ok := dependencyMap[dependent.Name]; !ok {
			dependencyMap[dependent.Name] = dependent.Version
			fmt.Println(dependent.Name + "=" + dependent.Version)

			collectDependencies(dependent.Name, dependencyMap)
		}
	}
}

func print(dependencyMap map[string]string, f *os.File) {
	for name, ver := range dependencyMap {
		f.WriteString(name + "=" + ver + "\n")
	}
}

type FlagArray []string

func (i *FlagArray) String() string {
	return "plugin name list"
}

func (i *FlagArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var names FlagArray

	flag.Var(&names, "names", "Plugin Name list")
	out := flag.String("out", "", "outputfile")
	flag.Parse()

	f, err := os.Create(*out)
	if err != nil {
		panic(err)
	}

	dependencyMap := make(map[string]string)
	for _, name := range names {
		collectDependencies(name, dependencyMap)
		print(dependencyMap, f)
	}

	defer f.Close()
}
