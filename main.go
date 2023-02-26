package main

import (
	"fmt"
	"net/http"
	"search-x/lib"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type RequestItem struct {
	Keywords       int    `json:"keywords"`
	AvailableOnly  int    `json:"availableOnly"`
	ForSaleTypes   string `json:"forSaleTypes"`
	PropertyTypes  string `json:"propertyTypes"`
	OtherAmenities string `json:"otherAmenities"`
	ViewTypes      string `json:"viewTypes"`
	PerPage        int    `json:"per_page"`
}

type ResponseItem struct {
	Id            string  `json:"id"`
	PhotoUri      string  `json:"photoUri"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	DisplayPrice  int     `json:"displayPrice"`
	Status        string  `json:"status"`
	Bedrooms      int     `json:"bedrooms"`
	FullBathrooms int     `json:"fullBathrooms"`
	HalfBathrooms int     `json:"halfBathrooms"`
	SquareFeet    int     `json:"squareFeet"`
	Address       string  `json:"address"`
	Unit          string  `json:"unit,omitempty"`
	City          string  `json:"city"`
	State         string  `json:"state"`
	Zip           int     `json:"zip"`
}

func prepare(json map[string]interface{}, nothing ResponseItem) map[string]interface{} {
	jsonCopy := make(map[string]interface{})
	for k, v := range json {
		jsonCopy[k] = v
	}
	fmt.Println("RAW REQUEST = ", jsonCopy)

	if jsonCopy["availableOnly"] == 1 {
		jsonCopy["webAvailable"] = true
	}
	delete(jsonCopy, "availableOnly")
	delete(jsonCopy, "forSaleTypes")
	delete(jsonCopy, "propertyType")

	if _, ok := jsonCopy["keywords"]; ok && lib.IsValidZipcode(jsonCopy["keywords"].(string)) {
		zip, _ := strconv.Atoi(jsonCopy["keywords"].(string))
		minLongitude, maxLongitude, minLatitude, maxLatitude := lib.BoundingRectangle(zip)
		jsonCopy["minLongitude"] = minLongitude
		jsonCopy["maxLongitude"] = maxLongitude
		jsonCopy["minLatitude"] = minLatitude
		jsonCopy["maxLatitude"] = maxLatitude
		jsonCopy["postalCode"] = zip
		delete(jsonCopy, "keywords")
	} else if _, ok := jsonCopy["north"]; ok {
		jsonCopy["minLongitude"] = jsonCopy["west"]
		jsonCopy["maxLongitude"] = jsonCopy["east"]
		jsonCopy["minLatitude"] = jsonCopy["south"]
		jsonCopy["maxLatitude"] = jsonCopy["north"]
		delete(jsonCopy, "west")
		delete(jsonCopy, "east")
		delete(jsonCopy, "south")
		delete(jsonCopy, "north")
	}

	// legacy fields never used
	delete(jsonCopy, "per_page")
	delete(jsonCopy, "locationType")

	fmt.Println("COOKED REQUEST = ", jsonCopy)
	return jsonCopy
}

func photoUri(item map[string]interface{}) string {
	if photoUriPath, ok := item["photoUriPath"]; ok {
		return "/main" + photoUriPath.(string)
	} else {
		return ""
	}
}

func intOrNone(i string) int {
	s, _ := strconv.Atoi(i)
	if s == 0 {
		return 0
	}
	result := s
	return result
}

func Search(c echo.Context) error {
	newest := make(map[string]interface{})
	var request ResponseItem
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	queryParameters := prepare(newest, request)
	result := lib.DynamoQuery(queryParameters)

	response := make([]ResponseItem, 0, len(result))
	for _, item := range result {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		respItem := ResponseItem{
			Id:            itemMap["id"].(string),
			PhotoUri:      photoUri(itemMap),
			Latitude:      itemMap["latitude"].(float64),
			Longitude:     itemMap["longitude"].(float64),
			DisplayPrice:  int(itemMap["listPrice"].(float64)),
			Status:        itemMap["status"].(string),
			Bedrooms:      int(itemMap["bedroomsTotal"].(float64)),
			FullBathrooms: int(itemMap["bathroomsTotalInteger"].(float64)) - int(itemMap["bathroomsHalf"].(float64)),
			HalfBathrooms: int(itemMap["bathroomsHalf"].(float64)),
			SquareFeet:    intOrNone(itemMap["livingArea"].(string)),
			Address:       strings.Split(itemMap["unitAddress"].(string), " #")[0],
			City:          itemMap["city"].(string),
			State:         itemMap["stateOrProvince"].(string),
			Zip:           int(itemMap["postalCode"].(float64)),
		}

		if unitAddress, ok := itemMap["unitAddress"].(string); ok && strings.Contains(unitAddress, " #") {
			respItem.Unit = strings.Split(unitAddress, " #")[1]
		}

		response = append(response, respItem)
	}
	return c.JSON(http.StatusOK, response)
}

func main() {
	e := echo.New()

	e.POST("/search-x.api", Search)

	e.Logger.Fatal(e.Start(":8080"))
}
