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

	fmt.Printf("Logfile filename %v\n", config.LogFile)

	if err != nil {
		fmt.Printf("Error read configuration file %v\n", err)
		log.Fatalf("Error read configuration file %v\n", err)

	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "Redacid's cmd Metrics Exporter %v\n", config.LogFile)
	fmt.Fprintf(w, "Redacid's cmd Metrics Exporter \n <a href=/metrics>Metrics</a>")
}

func metrics(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "# This is METRICS")
	fmt.Fprintln(w, "# Redacid's cmd Metrics Exporter ", config.LogFile)

	for _, CCommand := range config.Commands {

		cmd := exec.Command("bash", "-c", CCommand.Cmd)
		//cmd.Stdin = strings.NewReader("some input")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("in all caps: %q\n", out.String())
		fmt.Printf("%v - Cmd: %v Result: %v", CCommand.Name, CCommand.Cmd, out.String())
		fmt.Fprintln(w, "# HELP ", CCommand.Name, " ", CCommand.Help)
		fmt.Fprintln(w, "# TYPE ", CCommand.Name, " ", CCommand.Type)
		fmt.Fprintln(w, CCommand.Name, " ", out.String())
	}

}

func init() {
	readConfig()
	//fmt.Printf("Logfile filename in init %v\n", config.LogFile)
}

func main() {

	fmt.Printf("Logfile filename in main %v\n", config.LogFile)
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
