# Homesearch API

A simple API for searching a database of homes (houses) based on given latitude/longitude and additional filtering criteria.

## Prerequisites

- A DynamoDB table (you can use Local Stack or download local Dynamodb Installer for this if you do not have an AWS Account):

```
  RetsSearchTable:
    Type: 'AWS::DynamoDB::Table'
    Properties:
      TableName: !Sub "${EnvironmentPrefix}rets-search"
      BillingMode: PAY_PER_REQUEST
      AttributeDefinitions:
        - AttributeName: id
          AttributeType: S
        - AttributeName: latitude_box
          AttributeType: N
        - AttributeName: longitude
          AttributeType: N
      KeySchema:
        - AttributeName: id
          KeyType: HASH
      GlobalSecondaryIndexes:
        - IndexName: latitude-longitude-index
          KeySchema:
            - AttributeName: latitude_box
              KeyType: HASH
            - AttributeName: longitude
              KeyType: RANGE
          Projection:
            ProjectionType: INCLUDE
            NonKeyAttributes:
            - latitude
            - associationAmenities
            - associationFee
            - bathroomsTotalInteger
            - bedroomsTotal
            - garageSpaces
            - livingArea
            - postalCode
            - searchablePrice
            - status
            - storiesTotal
            - webAvailable
            - yearBuilt
```



## Features

- Search for homes within a specified latitude/longitude range
- Filter search results using additional criteria such as price, number of rooms, etc.

## Usage

The API uses the `POST` method to receive query parameters in the request body, which should be in JSON format.

The API requires the following query parameters:

- `minLatitude`: the minimum latitude of the search range (float)
- `maxLatitude`: the maximum latitude of the search range (float)
- `minLongitude`: the minimum longitude of the search range (float)
- `maxLongitude`: the maximum longitude of the search range (float)

Additionally, you can include any other fields as query parameters to filter the search results. The field names should be in the format `min[FieldName]` for minimum values and `max[FieldName]` for maximum values. For example:

- `minPrice`: the minimum price of the homes to search for
- `maxPrice`: the maximum price of the homes to search for

Fields not beginning with `min` or `max` are matched using equals or `IN`, for example:

- `City: [ 'city1', 'city2' ]` will match where City is either 'city1' or 'city2'

The API will return a JSON response containing the search results.

## Deployment

The API is built with [Echo](https://echo.labstack.com/), a High performance, extensible, minimalist Go web framework. To deploy the API, you'll need to:

1. Use the existing Golang Echo app (main file: `main.go`)
2. Create the DynamoDB table locally on your machine (using Local Stack or local dynamodb) or on AWS.
3. Deploy the app using `go run main.go` 