package detect

import (
	"encoding/json"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"log"
	"net/http"
	"strings"
)

const precipitationUrl string = "https://data.geo.admin.ch/ch.meteoschweiz.messwerte-niederschlag-10min/ch.meteoschweiz.messwerte-niederschlag-10min_en.json"

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
	Rain int
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
	var md MeteoData
	res, err := http.Get(precipitationUrl)
	if err != nil {
		log.Printf("error making http request: %s\n", err)
	} else {
		err := json.NewDecoder(res.Body).Decode(&md)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		AllMatchingFeatures := []Feature{}
		for _, s := range stations {
			findFeature := func(f Feature) bool { return strings.EqualFold(f.ID, s) }
			MatchingFeatureValues := Filter(md.Features, findFeature)
			AllMatchingFeatures = append(AllMatchingFeatures, MatchingFeatureValues...)
		}
		log.Printf("values for matching features: %+v", AllMatchingFeatures)
		var rainDetected = 0
		for _, f := range AllMatchingFeatures {
			if f.Properties.Value > 0.0 {
				rainDetected = 1
			}
		}
		u := DetectionResult{Rain: rainDetected}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u)
	}
}

func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
