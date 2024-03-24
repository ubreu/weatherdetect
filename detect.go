package detect

import (
	"fmt"
	"net/http"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

const precipitationUrl string = "https://data.geo.admin.ch/ch.meteoschweiz.messwerte-niederschlag-10min/ch.meteoschweiz.messwerte-niederschlag-10min_en.json"

func init() {
	functions.HTTP("DetectForLocation", DetectForLocation)
}

func DetectForLocation(w http.ResponseWriter, r *http.Request) {
	station := r.URL.Query()["station"]
	if len(station) > 0 {
	  fmt.Fprint(w, station[0])
	  fmt.Fprint(w, station[1])
	  fmt.Fprint(w, len(station))
	}    

	res, err := http.Get(precipitationUrl)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
	} else {
		fmt.Printf("client: got response!\n")
		fmt.Printf("client: status code: %d\n", res.StatusCode)
	}


}
