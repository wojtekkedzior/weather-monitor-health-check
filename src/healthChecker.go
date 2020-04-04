package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Info struct {
	TempFromDB   *Temperature
	TempFromJSON *Temperature
	HostInfo     *HostInfo
}

type Temperature struct {
	Id          int       `db:"id" json:"id"`
	Humidity    float32   `json:"humidity"`
	Pressure    float32   `json:"pressure"`
	Temperature float32   `json:"temperature"`
	Timestamp   time.Time `json:"timestamp"`
}

type HostInfo struct {
	Hostname string
}

func main() {
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/healthy", healthHandler)
	http.HandleFunc("/ready", readyHandler)
	http.HandleFunc("/process", processHandler)
	log.Fatal(http.ListenAndServe(":8090", nil))
}

// For K8
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	var hostInfo HostInfo
	var tempFromDB, tempFromJSON Temperature

	info := Info{TempFromDB: &tempFromDB, TempFromJSON: &tempFromJSON, HostInfo: &hostInfo}

	tempFromDB.getFromDB()
	tempFromJSON.getFromJava()
	hostInfo.getHostInfo()

	info.marshal(w)

	fmt.Println(info.TempFromJSON)
	fmt.Println(info.TempFromDB)
}

func (info *HostInfo) getHostInfo() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err.Error())
	}

	info.Hostname = hostname
}

func (temp *Temperature) getFromDB() {
	db, err := sql.Open("mysql", dbUserName+":"+dbPassword+"@tcp("+dbHost+":3306+)/iot?parseTime=true")

	defer db.Close()

	if err != nil {
		panic(err.Error())
	}

	err = db.QueryRow("SELECT * FROM temperature ORDER BY id DESC LIMIT 1").Scan(&temp.Id, &temp.Humidity, &temp.Pressure, &temp.Temperature, &temp.Timestamp)
}

func (temp *Temperature) getFromJava() {
	resp, err := http.Get("http://" + backend + ":8080/weather-monitor/temperature/getLatest")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		bodyString := string(bodyBytes)
		fmt.Println(bodyString)

		jerr := json.Unmarshal(bodyBytes, &temp)

		if jerr != nil {
			log.Fatal(jerr)
		}
	}
}

func (t *Info) marshal(w http.ResponseWriter) {
	js, err := json.Marshal(t)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func printFlagValue(f *flag.Flag) {
	fmt.Println("Flag: ", f)
}

var dbUserName, dbPassword, dbHost, backend string

func init() {
	flag.StringVar(&dbUserName, "db-username", "", "DB username")
	flag.StringVar(&dbPassword, "db-password", "", "DB password")
	flag.StringVar(&dbHost, "db-host", "localhost", "DB host")
	flag.StringVar(&backend, "backend", "localhost", "The Weather Monitor backend")

	flag.Parse()

	fmt.Println("Running with: ")
	flag.VisitAll(printFlagValue)
}
