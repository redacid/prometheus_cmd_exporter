package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Config struct {
	Commands []Command
	Global
}

type Command struct {
	Name string `json:"Name"`
	Cmd  string `json:"Cmd"`
}

type Global struct {
	LogFile string `json:"logFile"`
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
	fmt.Fprintln(w, "# Redacid's cmd Metrics Exporter %v", config.LogFile)

	for _, CCommand := range config.Commands {
		fmt.Fprintln(w, CCommand.Name, " ", "Cmd: ", CCommand.Cmd)
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
	log.Fatal(http.ListenAndServe(":8888", nil))
}
