package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

	//fmt.Println(string(body))

	plugin = &Plugin{}
	err = json.Unmarshal(body, plugin)
	return
}

func printDependencies(pluginName string, dependencyMap map[string]string) {
	plugin, err := getPlugin(pluginName)
	if err != nil {
		panic(err)
	}

	dependencyMap[plugin.Name] = plugin.Version

	for _, dependent := range plugin.Dependencies {
		if _, ok := dependencyMap[dependent.Name]; !ok {
			dependencyMap[dependent.Name] = dependent.Version
			fmt.Println(dependent.Name + "=" + dependent.Version)

			printDependencies(dependent.Name, dependencyMap)
		}
	}
}

func main() {
	name := flag.String("name", "hugo", "Plugin Name")
	out := flag.String("out", "", "outputfile")

	flag.Parse()

	dependencyMap := make(map[string]string)
	printDependencies(*name, dependencyMap)

	f, err := os.Create(*out)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	for name, ver := range dependencyMap {
		f.WriteString(name + "=" + ver + "\n")
	}
}
