package detect

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

const precipitationUrl string = "https://data.geo.admin.ch/ch.meteoschweiz.messwerte-niederschlag-10min/ch.meteoschweiz.messwerte-niederschlag-10min_en.json"
const sunshineUrl string = "https://data.geo.admin.ch/ch.meteoschweiz.messwerte-sonnenscheindauer-10min/ch.meteoschweiz.messwerte-sonnenscheindauer-10min_en.json"

// converted using https://mholt.github.io/json-to-go/
type Feature struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Properties struct {
		StationName       string  `json:"station_name"`
		StationSymbol     int     `json:"station_symbol"`
		Value             float64 `json:"value"`
		Unit              string  `json:"unit"`
		Altitude          string  `json:"altitude"`
		MeasurementHeight string  `json:"measurement_height"`
	} `json:"properties"`
}

type MeteoData struct {
	CreationTime string    `json:"creation_time"`
	Type         string    `json:"type"`
	Features     []Feature `json:"features"`
}

type DetectionResult struct {
	Rain     int
	Sunshine int
}

func init() {
	functions.HTTP("DetectForLocation", DetectForLocation)
}

func DetectForLocation(w http.ResponseWriter, r *http.Request) {
	stations := r.URL.Query()["station"]
	if len(stations) == 0 {
		http.Error(w, "missing required station parameter", http.StatusBadRequest)
		return
	}

	u := DetectionResult{
		Rain:     detectRain(stations),
		Sunshine: detectSunshine(stations),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func detectRain(stations []string) int {
	var detected = 0
	var md MeteoData
	res, err := http.Get(precipitationUrl)
	if err != nil {
		log.Printf("error making http request: %s\n", err)
	} else {
		err := json.NewDecoder(res.Body).Decode(&md)
		if err != nil {
			log.Printf("error decoding http request: %s\n", err)
			return 0
		}

		AllMatchingFeatures := []Feature{}
		for _, s := range stations {
			findFeature := func(f Feature) bool { return strings.EqualFold(f.ID, s) }
			MatchingFeatureValues := Filter(md.Features, findFeature)
			AllMatchingFeatures = append(AllMatchingFeatures, MatchingFeatureValues...)
		}
		log.Printf("rain values for matching features: %+v", AllMatchingFeatures)
		for _, f := range AllMatchingFeatures {
			if f.Properties.Value > 0.0 && f.Properties.Value < 1000.0 {
				detected = 1
			}
		}
	}
	return detected
}

func detectSunshine(stations []string) int {
	var detected = 0
	var md MeteoData
	res, err := http.Get(sunshineUrl)
	if err != nil {
		log.Printf("error making http request: %s\n", err)
	} else {
		err := json.NewDecoder(res.Body).Decode(&md)
		if err != nil {
			log.Printf("error decoding http request: %s\n", err)
			return 0
		}

		AllMatchingFeatures := []Feature{}
		for _, s := range stations {
			findFeature := func(f Feature) bool { return strings.EqualFold(f.ID, s) }
			MatchingFeatureValues := Filter(md.Features, findFeature)
			AllMatchingFeatures = append(AllMatchingFeatures, MatchingFeatureValues...)
		}
		log.Printf("sunshine values for matching features: %+v", AllMatchingFeatures)
		for _, f := range AllMatchingFeatures {
			if f.Properties.Value > 4.0 && f.Properties.Value <= 10.0 {
				detected = 1
			}
		}
	}
	return detected
}

func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
