package dynamo

import (
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type User struct {
	UserName string `dynamodbav:"UserName" json:"user_name" binding:"required"`
	Password string `dynamodbav:"Password" json:"password"`
	RoleName string `dynamodbav:"RoleName" json:"role_name"`
}

type Relation struct {
	UserName    string `dynamodbav:"UserName"  json:"user_name" binding:"required"`
	ProjectName string `dynamodbav:"ProjectName"  json:"project_name" binding:"required"`
}

type Env struct {
	EnvKey   string `json:"env_key" binding:"required"`
	EnvValue string `json:"env_value" binding:"required"`
}

type Project struct {
	ProjectName string `dynamodbav:"ProjectName" json:"project_name" binding:"required"`
	EnvList     []Env  `dynamodbav:"EnvList" json:"env_list"`
}

var d *dynamodb.DynamoDB
var once sync.Once

type MyProvider struct{}

func (m *MyProvider) Retrieve() (v credentials.Value, err error) {
	return
}
func (m *MyProvider) IsExpired() bool {
	return false
}

func connect() *dynamodb.DynamoDB {
	creds := credentials.NewCredentials(&MyProvider{})

	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Endpoint:    aws.String("http://localhost:4569"),
		Credentials: creds,
	})

	// Create DynamoDB client
	return dynamodb.New(sess)
}

func GetDBInstance() *dynamodb.DynamoDB {
	once.Do(func() {
		d = connect()
	})

	return d
}

func createDDBTable(input *dynamodb.CreateTableInput) {
	_, err := d.CreateTable(input)
	if err != nil {
		fmt.Println("Got error calling CreateTable:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Created the table", *input.TableName)
}

func createLocalDDB() {
	createDDBTable(getUserTable())
	createDDBTable(getRelationTable())
	createDDBTable(getProjectTable())
}

func getUserTable() *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("UserName"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("UserName"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("User"),
	}
}

func getRelationTable() *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("UserName"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("ProjectName"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("UserName"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("ProjectName"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("Relation"),
	}
}

func getProjectTable() *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ProjectName"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ProjectName"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("Project"),
	}
}

func listDDB() {
	// create the input configuration instance
	input := &dynamodb.ListTablesInput{}

	for {
		// Get the list of tables
		result, err := d.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}

		for _, n := range result.TableNames {
			fmt.Println(*n)
		}

		// assign the last read tablename as the start for our next call to the ListTables function
		// the maximum number of table names returned in a call is 100 (default), which requires us to make
		// multiple calls to the ListTables function to retrieve all table names
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}
}

func createItem(userStr string) {
	user := User{
		UserName: userStr,
		Password: userStr,
		RoleName: userStr,
	}

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		fmt.Println("Got error marshalling new user item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("User"),
	}

	_, err = d.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func New() {
	GetDBInstance()
	// createLocalDDB()
	// createDDBTable(getProjectTable())
	// createItem("user")
	// createItem("admin")
}
