package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Config struct {
	Commands []Command
	Global
}

type Command struct {
	Name string `json:"Name"`
	Cmd  string `json:"Cmd"`
	Help string `json:"Help"`
	Type string `json:"Type"`
}

type Global struct {
	LogFile string `json:"LogFile"`
	Listen  string `json:"Listen"`
}

var config Config

func readConfig() {
	appdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	conffile, _ := os.Open(appdir + "/config.json")
	decoder := json.NewDecoder(conffile)
	//config := new(Config)
	err = decoder.Decode(&config)
	defer conffile.Close()

	fmt.Printf("Logfile filename: %v\n", config.LogFile)
	fmt.Printf("Listen: %v\n", config.Listen)

	if err != nil {
		fmt.Printf("Error read configuration file %v\n", err)
		log.Fatalf("Error read configuration file %v\n", err)

	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body><h3>Redacid's cmd Metrics Exporter</h3> <br /> <a href=/metrics>Metrics</a></body></html>")
}

func metrics(w http.ResponseWriter, r *http.Request) {

	for _, CCommand := range config.Commands {

		cmd := exec.Command("bash", "-c", CCommand.Cmd)
		//cmd.Stdin = strings.NewReader("some input")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		//fmt.Printf("%v - Cmd: %v Result: %v", CCommand.Name, CCommand.Cmd, out.String())
		fmt.Fprintln(w, "# HELP ", CCommand.Name, CCommand.Help)
		fmt.Fprintln(w, "# TYPE ", CCommand.Name, CCommand.Type)
		fmt.Fprintln(w, CCommand.Name, strings.Trim(out.String(), "\n"))
	}

}

func init() {
	readConfig()
}

func main() {

	f, err1 := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err1 != nil {
		fmt.Printf("error opening log file: %v %v\n", err1, config.LogFile)
		log.Fatalf("error opening log file: %v %v\n", err1, config.LogFile)
	}
	defer f.Close()
	log.SetOutput(f)

	http.HandleFunc("/", handler)
	http.HandleFunc("/metrics/", metrics)
	log.Fatal(http.ListenAndServe(config.Listen, nil))
}
