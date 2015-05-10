package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gotstago/reference/handlers"
	"log"
	"net/http"
	"os"
	"strings"
)

var logFile *os.File

func main() {
	var err error
	logFile, err = os.Create("logfile.txt")
	if err != nil {
		log.Fatal("Log file create:", err)
		return
	}
	defer logFile.Close()
	helloHandler := new(handlers.HelloHandler)
	getFileHandler := new(handlers.GetFileHandler)
	fmt.Fprintf(logFile, "starting again here\n")
	http.HandleFunc("/templates/provider/ftp", getFileHandler.ServeHTTP)
	http.HandleFunc("/hellos", helloHandler.ServeHTTP)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		fmt.Fprintf(logFile, "%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		fmt.Fprintf(os.Stdout, "%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		//data, err := query(city)
		data, err := environmentCanada{}.temperature(city)
		if err != nil {
			fmt.Fprintf(logFile, "%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})
	http.ListenAndServe(":8080", nil)
}

/*type HelloHandler struct{}

func (e HelloHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	sayParam := r.FormValue("say")

	if sayParam == "Nothing" {
		rw.WriteHeader(404)
	} else {
		//rw.Write([]byte(sayParam))
		rw.Write([]byte("hello!\n"))
	}
}*/

// func hello(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("hello!\n"))
// }

func query(city string) (weatherData, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city)
	if err != nil {
		fmt.Println("no info for ", city)
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var d weatherData

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		fmt.Println("no info for ", city)
		return weatherData{}, err
	}

	return d, nil
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

type weatherProvider interface {
	temperature(city string) (float64, error) // in Kelvin, naturally
}

type environmentCanada struct{}

type XMLLink struct {
	XMLName xml.Name `xml:"link"`
	Type    string   `xml:"type,attr"`
	Href    string   `xml:"href,attr"`
}

type XMLEntry struct {
	XMLName xml.Name `xml:"entry"`
	Title   string   `xml:"title"`
	Link    XMLLink  `xml:"link"`
}

type XMLFeed struct {
	XMLName xml.Name   `xml:"feed"`
	Entries []XMLEntry `xml:"entry"`
}

func (w environmentCanada) temperature(city string) (float64, error) {
	resp, err := http.Get("http://weather.gc.ca/rss/city/ns-31_e.xml")
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	var xmlFeed = XMLFeed{}
	/*var d struct {
		Main struct {
			Kelvin float64 `json:"temp"`
		} `json:"main"`
	}*/

	if err := xml.NewDecoder(resp.Body).Decode(&xmlFeed); err != nil {
		return 0, err
	}
	// Display The first strap
	fmt.Printf("Title: %s  Link: %s", xmlFeed.Entries[1].Title, xmlFeed.Entries[1].Link.Href)
	fmt.Printf("Title: %s  Link: %s", xmlFeed.Entries[1].Title, xmlFeed.Entries[1].Link.Href)

	//log.Printf("environmentCanada: %s: %.2f", xmlFeeds, d.Main.Kelvin)
	return 0, nil
}

type openWeatherMap struct{}

func (w openWeatherMap) temperature(city string) (float64, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	var d struct {
		Main struct {
			Kelvin float64 `json:"temp"`
		} `json:"main"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}

	log.Printf("openWeatherMap: %s: %.2f", city, d.Main.Kelvin)
	return d.Main.Kelvin, nil
}
