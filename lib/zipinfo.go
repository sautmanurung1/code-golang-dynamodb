package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

var data map[string]interface{}

// init is called before the program starts
func init() {
	// Get the absolute path of the directory where the script is located
	scriptDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	// Open the file using the relative file path
	filePath := filepath.Join(scriptDir, "ZIP_CODES.geojson")
	dataBytes, _ := ioutil.ReadFile(filePath)

	// Parse the file contents as JSON
	err := json.Unmarshal(dataBytes, &data)
	if err != nil {
		return
	}
}

func IsValidZipcode(zipcode string) bool {
	zipcodePattern := regexp.MustCompile(`^\d{5}(?:[-\s]\d{4})?$`)
	return zipcodePattern.MatchString(zipcode)
}

func BoundingRectangle(zipCode int) (float64, float64, float64, float64) {
	convZipCode := strconv.Itoa(zipCode)
	// Find the feature with the given zip code
	feature := map[string]interface{}{}
	for _, feat := range data["features"].([]interface{}) {
		feat := feat.(map[string]interface{})
		if feat["properties"].(map[string]interface{})["ZIP"].(string) == convZipCode {
			feature = feat
			break
		}
	}

	if feature != nil {
		// Extract the bounding rectangle coordinates
		coords := feature["geometry"].(map[string]interface{})["coordinates"].([]interface{})[0].([]interface{})
		lons, lats := []float64{}, []float64{}
		for _, coord := range coords {
			lon, lat := coord.([]interface{})[0].(float64), coord.([]interface{})[1].(float64)
			lons, lats = append(lons, lon), append(lats, lat)
		}
		minLon, maxLon := minMax(lons)
		minLat, maxLat := minMax(lats)
		return minLon, maxLon, minLat, maxLat
	} else {
		return 0, 0, 0, 0
	}
}

func minMax(arr []float64) (float64, float64) {
	min, max := arr[0], arr[0]
	for _, val := range arr {
		if val < min {
			min = val
		}
		if val > max {
			max = val
		}
	}
	return min, max
}
