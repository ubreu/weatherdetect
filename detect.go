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

		ValuesOfMatchingFeatures := []float64{}
		for _, s := range stations {
			findFeature := func(f Feature) bool { return strings.EqualFold(f.ID, s) }
			mapValue := func(f Feature) float64 { return f.Properties.Value }
			MatchingFeatureValues := Map(Filter(md.Features, findFeature), mapValue)
			ValuesOfMatchingFeatures = append(ValuesOfMatchingFeatures, MatchingFeatureValues...)
		}
		log.Printf("values for matching features: %+v", ValuesOfMatchingFeatures)
		var rainDetected = 0
		for _, v := range ValuesOfMatchingFeatures {
			if v > 0.0 {
				rainDetected = 1
			}
		}
		u := DetectionResult{Rain: rainDetected}
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

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}
