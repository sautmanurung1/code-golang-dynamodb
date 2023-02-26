package lib

import (
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const retsSearchTable = "rets-search"

var dynamodbClient *dynamodb.DynamoDB // initialize the client outside the function

// Initialize the client
func InitDynamoDBClient() {
	session := session.Must(session.NewSession())
	dynamodbClient = dynamodb.New(session)
}

func Deserialize(dynamodbJSON map[string]*dynamodb.AttributeValue) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	err := dynamodbattribute.UnmarshalMap(dynamodbJSON, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func LatitudeBoxValues(minLat, maxLat float64) []int {
	minLatInt := int(minLat * 10)
	maxLatInt := int(maxLat * 10)

	values := []int{}
	for i := minLatInt; i <= maxLatInt; i++ {
		values = append(values, i)
	}
	return values
}

// DynamoQuery performs a query using the provided parameters
func DynamoQuery(queryParameters map[string]interface{}) []interface{} {
	latitudeBoxes := LatitudeBoxValues(
		queryParameters["minLatitude"].(float64),
		queryParameters["maxLatitude"].(float64))

	result := []interface{}{}
	var wg sync.WaitGroup
	wg.Add(len(latitudeBoxes))

	for _, latitudeBox := range latitudeBoxes {
		go ExecuteQueryThread(queryParameters, latitudeBox)
		wg.Done()
	}

	wg.Wait()
	return result
}

// ExecuteQueryThread is a function to be run in a goroutine to execute a query
func ExecuteQueryThread(queryParameters map[string]interface{}, latitudeBox int) {
	query := CreateQuery(queryParameters, latitudeBox)
	var result []interface{}
	consumedCapacity := ExecuteQuery(result, query)
	fmt.Println("CONSUMED_CAPACITY = ", consumedCapacity)
}

// ExecuteQuery executes the query and adds any results to the given list
func ExecuteQuery(result []interface{}, query map[string]*dynamodb.AttributeValue) float64 {
	var consumedCapacity float64 = 0
	lastKey := map[string]*dynamodb.AttributeValue{}
	queryInput := &dynamodb.QueryInput{
		KeyConditionExpression:    aws.String(*query["KeyConditionExpression"].S),
		ExpressionAttributeValues: query,
		TableName:                 aws.String(retsSearchTable),
	}
	if len(lastKey) > 0 {
		queryInput.ExclusiveStartKey = lastKey
	}
	for {
		var response *dynamodb.QueryOutput
		var err error // declare err variable
		if len(lastKey) > 0 {
			queryInput.ExclusiveStartKey = lastKey // set lastKey
		}
		response, err = dynamodbClient.Query(queryInput) // call Query method
		if err != nil {
			panic(err.Error())
		}

		// Deserialize response items
		items := []map[string]interface{}{}
		for _, item := range response.Items {
			result, err := Deserialize(item)
			if err != nil {
				panic(err.Error())
			}

			items = append(items, result)
		}

		consumedCapacity += *response.ConsumedCapacity.CapacityUnits

		// Add the items to the given result list
		result = append(result, items)

		// Check if there are more records
		lastKey = response.LastEvaluatedKey
		if len(lastKey) == 0 {
			break
		}
	}

	return consumedCapacity
}

func BuildBaseQuery(queryParams map[string]interface{}, latitudeBox string) (*dynamodb.QueryInput, int) {
	var indexName, keyName string
	if queryParams["webAvailable"].(bool) {
		indexName = "latitude-longitude-webavailable-index"
		keyName = "latitude_box_webavailable"
	} else {
		indexName = "latitude-longitude-index"
		keyName = "latitude_box"
	}

	query := &dynamodb.QueryInput{
		TableName:              aws.String(retsSearchTable),
		IndexName:              aws.String(indexName),
		ReturnConsumedCapacity: aws.String("TOTAL"),
		KeyConditionExpression: aws.String(fmt.Sprintf("%s = :latitude_box AND longitude BETWEEN :minLongitude AND :maxLongitude", keyName)),
		FilterExpression:       aws.String(""),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":latitude_box": {
				N: aws.String(latitudeBox),
			},
			":minLongitude": {
				N: aws.String(fmt.Sprint(queryParams["minLongitude"])),
			},
			":maxLongitude": {
				N: aws.String(fmt.Sprint(queryParams["maxLongitude"])),
			},
		},
	}

	for key, value := range queryParams {
		if key == "minLongitude" || key == "maxLongitude" || key == "webAvailable" {
			continue
		}

		addFilter(queryParams, query, key, value)
		// remove trailing 'AND '
		filterLength := len(*query.FilterExpression)
		query.FilterExpression = aws.String((*query.FilterExpression)[:filterLength-4])
	}

	return query, 0
}

func addFilter(queryParams map[string]interface{}, query *dynamodb.QueryInput, key string, value interface{}) {
	if key[:3] == "min" {
		field := key[1:]
		field = field[:1] + strings.ToLower(field[1:])

		*query.FilterExpression += fmt.Sprintf("%s >= :%s_min AND ", field, field)
		query.ExpressionAttributeValues[":min"+field] = &dynamodb.AttributeValue{
			N: aws.String(fmt.Sprint(value)),
		}
	} else if key[:3] == "max" {
		field := key[1:]
		field = field[:1] + strings.ToLower(field[1:])
		*query.FilterExpression += fmt.Sprintf("%s <= :%s_max AND ", field, field)
		query.ExpressionAttributeValues[":max"+field] = &dynamodb.AttributeValue{
			N: aws.String(fmt.Sprint(value)),
		}
	}
}
