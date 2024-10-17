package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"jazz/backend/configs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoDBCache Implementation of Cache using DynamoDB.
type DynamoDBCache struct {
	client *dynamodb.DynamoDB
	table  string
}

// NewDynamoDBCache initializes a new DynamoDB Cache.
func NewDynamoDBCache() *DynamoDBCache {
	cacheConfig := configs.GetCacheConfig()["stores"].(map[string]interface{})["dynamodb"].(map[string]interface{})

	awsRegion := cacheConfig["region"].(string)
	awsAccessKey := cacheConfig["key"].(string)
	awsSecretKey := cacheConfig["secret"].(string)
	tableName := cacheConfig["table"].(string)

	if awsAccessKey == "" || awsSecretKey == "" {
		fmt.Println("AWS credentials are not set. Falling back to default cache.")
		return nil
	}

	awsConfig := &aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
	}

	if endpoint, ok := cacheConfig["endpoint"].(string); ok && endpoint != "" {
		awsConfig.Endpoint = aws.String(endpoint)
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		fmt.Println("Failed to create AWS session:", err)
		return nil
	}

	client := dynamodb.New(sess)

	fmt.Printf("Connected to DynamoDB at table: %s in region: %s\n", tableName, awsRegion)
	return &DynamoDBCache{
		client: client,
		table:  tableName,
	}
}

// Set stores a value in DynamoDB.
func (d *DynamoDBCache) Set(key string, value interface{}, expiration time.Duration) error {
	expirationTime := time.Now().Add(expiration).Unix()
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	item := map[string]*dynamodb.AttributeValue{
		"Key": {
			S: aws.String(key),
		},
		"Value": {
			S: aws.String(string(valueBytes)),
		},
		"Expiration": {
			N: aws.String(fmt.Sprintf("%d", expirationTime)),
		},
	}

	_, err = d.client.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.table),
		Item:      item,
	})

	return err
}

// Remember stores a value in DynamoDB using a callback if the value does not already exist.
func (d *DynamoDBCache) Remember(key string, expiration time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	var result interface{}

	// Check if value is already in the cache
	value, err := d.Get(key)
	if err == nil && value != nil {
		return value, nil
	}

	// If value is not cached, execute the callback
	result, err = callback()
	if err != nil {
		return result, err
	}

	// Cache the value
	if err := d.Set(key, result, expiration); err != nil {
		return result, err
	}

	return result, nil
}

// Forget removes a value from DynamoDB.
func (d *DynamoDBCache) Forget(key string) error {
	_, err := d.client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(d.table),
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(key),
			},
		},
	})
	return err
}

// Get retrieves a value from DynamoDB.
func (d *DynamoDBCache) Get(key string) (interface{}, error) {
	result, err := d.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(d.table),
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(key),
			},
		},
	})

	if err != nil || result.Item == nil {
		return nil, nil
	}

	return *result.Item["Value"].S, nil
}
