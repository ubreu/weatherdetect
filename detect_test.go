package detect

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFilter(t *testing.T) {
	features := []Feature{
		{ID: "SMA"},
		{ID: "BRZ"},
		{ID: "GVE"},
	}

	result := Filter(features, func(f Feature) bool { return f.ID == "SMA" })
	if len(result) != 1 || result[0].ID != "SMA" {
		t.Errorf("expected [SMA], got %v", result)
	}

	result = Filter(features, func(f Feature) bool { return false })
	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

func TestFilterCaseInsensitive(t *testing.T) {
	features := []Feature{{ID: "SMA"}, {ID: "brz"}}
	result := Filter(features, func(f Feature) bool { return f.ID == "sma" })
	// Filter itself is case-sensitive; EqualFold is used in detect()
	if len(result) != 0 {
		t.Errorf("Filter is case-sensitive, expected 0 matches for lowercase 'sma'")
	}
}

func TestDetectForLocation_MissingStation(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	DetectForLocation(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestDetectForLocation_WithMockServer(t *testing.T) {
	mockData := `{
		"creation_time": "2024-01-01T00:00:00Z",
		"type": "FeatureCollection",
		"map_short_name": "precipitation",
		"features": [
			{
				"type": "Feature",
				"id": "SMA",
				"properties": {
					"station_name": "Zürich",
					"station_symbol": 1,
					"value": 1.5,
					"unit": "mm",
					"altitude": "556",
					"measurement_height": "0"
				}
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockData))
	}))
	defer server.Close()

	// Test detect() directly with mock server URL
	result := detect([]string{"SMA"}, server.URL, 0.0, 1000.0)
	if result != 1 {
		t.Errorf("expected rain detected (1), got %d", result)
	}

	// Value 1.5 is NOT in sunshine range (3.0, 10.0]
	result = detect([]string{"SMA"}, server.URL, 3.0, 10.0)
	if result != 0 {
		t.Errorf("expected no sunshine (0), got %d", result)
	}
}

func TestDetectForLocation_JSONResponse(t *testing.T) {
	mockData := `{
		"creation_time": "2024-01-01T00:00:00Z",
		"type": "FeatureCollection",
		"map_short_name": "test",
		"features": []
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockData))
	}))
	defer server.Close()

	result := detect([]string{"NONEXISTENT"}, server.URL, 0.0, 1000.0)
	if result != 0 {
		t.Errorf("expected 0 for unknown station, got %d", result)
	}
}

func TestDetectForLocation_InvalidServer(t *testing.T) {
	result := detect([]string{"SMA"}, "http://localhost:0/invalid", 0.0, 1000.0)
	if result != 0 {
		t.Errorf("expected 0 on HTTP error, got %d", result)
	}
}

func TestDetectionResultJSONTags(t *testing.T) {
	r := DetectionResult{Rain: 1, Sunshine: 0}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var m map[string]int
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if _, ok := m["rain"]; !ok {
		t.Error("expected lowercase 'rain' key in JSON output")
	}
	if _, ok := m["sunshine"]; !ok {
		t.Error("expected lowercase 'sunshine' key in JSON output")
	}
}
