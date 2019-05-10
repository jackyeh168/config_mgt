package dynamo

import (
	"auth/util"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetUsers() (error, []User) {
	var users []User
	tableName := "User"

	params := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	// Make the DynamoDB Query API call
	result, err := d.Scan(params)
	util.Check(err)

	for _, i := range result.Items {
		user := User{}
		err = dynamodbattribute.UnmarshalMap(i, &user)
		users = append(users, user)
		util.Check(err)
	}

	return nil, users
}

func SaveUser(user User) bool {
	if len(user.Password) < 2 || user.Password[:2] != "$2" {
		user.Password = util.Encrypt(user.Password)
	}

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("User"),
	}

	_, err = d.PutItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

func SaveUsers(users []User) bool {
	res := true
	for _, user := range users {
		res = SaveUser(user)
	}
	return res
}

func UpdateUser(user User) error {
	if len(user.Password) < 2 || user.Password[:2] != "$2" {
		user.Password = util.Encrypt(user.Password)
	}

	tableName := "User"

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":p": {
				S: aws.String(user.Password),
			},
			":r": {
				S: aws.String(user.RoleName),
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"UserName": {
				S: aws.String(user.UserName),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set Password = :p, RoleName = :r"),
	}

	_, err := d.UpdateItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func DeleteUser(user User) error {
	tableName := "User"

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"UserName": {
				S: aws.String(user.UserName),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := d.DeleteItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
