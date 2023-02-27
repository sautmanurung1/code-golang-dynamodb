package lib

import (
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const retsSearchTable = "rets-search"

var dynamodbClient *dynamodb.DynamoDB // initialize the client outside the function

// Initialize the client
//func InitDynamoDBClient() {
//	sess := session.Must(session.NewSession())
//	dynamodbClient = dynamodb.New(sess)
//}

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

	var values []int
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

	var result []interface{}
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
	convLatitudeBox := float64(latitudeBox)
	query := CreateQuery(queryParameters, convLatitudeBox)
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
		var items []map[string]interface{}
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

func BuildBaseQuery(queryParameters map[string]interface{}, latitudeBox float64) (*dynamodb.QueryInput, map[string]*dynamodb.AttributeValue) {
	var indexName string
	var keyName string

	if queryParameters["webAvailable"] == true {
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
	}

	return query, map[string]*dynamodb.AttributeValue{
		":latitude_box": {N: aws.String(fmt.Sprintf("%.f", latitudeBox))},
		":minLongitude": {N: aws.String(fmt.Sprint(queryParameters["minLongitude"].(float64)))},
		":maxLongitude": {N: aws.String(fmt.Sprint(queryParameters["maxLongitude"].(float64)))},
	}
}

func AddFilter(queryParameters map[string]interface{}, query *dynamodb.QueryInput, key string) {
	if key[:3] == "min" {
		field := key[3:]
		field = strings.ToLower(field[:1]) + field[1:]
		query.FilterExpression = aws.String(fmt.Sprintf("%s >= :%s_min AND ", field, field))
		query.ExpressionAttributeValues[fmt.Sprintf(":%s_min", field)] = &dynamodb.AttributeValue{
			N: aws.String(fmt.Sprintf("%v", queryParameters[key])),
		}
	} else if key[:3] == "max" {
		field := key[3:]
		field = strings.ToLower(field[:1]) + field[1:]
		query.FilterExpression = aws.String(fmt.Sprintf("%s <= :%s_max AND ", field, field))
		query.ExpressionAttributeValues[fmt.Sprintf(":%s_max", field)] = &dynamodb.AttributeValue{
			N: aws.String(fmt.Sprintf("%v", queryParameters[key])),
		}
	} else {
		field := key
		value := queryParameters[key]
		if values, ok := value.([]string); ok {
			if len(values) > 0 {
				var exprValues []string
				for idx, v := range values {
					query.ExpressionAttributeValues[fmt.Sprintf(":%s_%d", field, idx)] = &dynamodb.AttributeValue{
						S: aws.String(v),
					}
					exprValues = append(exprValues, fmt.Sprintf(":%s_%d", field, idx))
				}
				query.FilterExpression = aws.String(fmt.Sprintf("%s IN (%s) AND ", field, strings.Join(exprValues, ",")))
			}
		} else {
			switch v := value.(type) {
			case string:
				query.ExpressionAttributeValues[fmt.Sprintf(":%s", field)] = &dynamodb.AttributeValue{
					S: aws.String(v),
				}
				query.FilterExpression = aws.String(fmt.Sprintf("%s = :%s AND ", field, field))
			case bool:
				query.ExpressionAttributeValues[fmt.Sprintf(":%s", field)] = &dynamodb.AttributeValue{
					BOOL: aws.Bool(v),
				}
				query.FilterExpression = aws.String(fmt.Sprintf("%s = :%s AND ", field, field))
			case int:
				query.ExpressionAttributeValues[fmt.Sprintf(":%s", field)] = &dynamodb.AttributeValue{
					N: aws.String(fmt.Sprintf("%d", v)),
				}
				query.FilterExpression = aws.String(fmt.Sprintf("%s = :%s AND ", field, field))
			case float64:
				query.ExpressionAttributeValues[fmt.Sprintf(":%s", field)] = &dynamodb.AttributeValue{
					N: aws.String(fmt.Sprintf("%f", v)),
				}
				query.FilterExpression = aws.String(fmt.Sprintf("%s = :%s AND ", field, field))
			default:
				panic(fmt.Sprintf("Unsupported data type for %s: %T", field, value))
			}
		}
	}
}

func CreateQuery(queryParameters map[string]interface{}, latitudeBox float64) map[string]*dynamodb.AttributeValue {
	query, query2 := BuildBaseQuery(queryParameters, latitudeBox)
	for key := range queryParameters {
		if key == "minLongitude" || key == "maxLongitude" {
			continue
		}
		AddFilter(queryParameters, query, key)
	}

	fmt.Println("query = ", query)
	return query2
}
