package dynamo

import (
	"auth/util"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetRelations() (error, []Relation) {
	var relations []Relation
	tableName := "Relation"

	params := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	// Make the DynamoDB Query API call
	result, err := d.Scan(params)
	util.Check(err)

	for _, i := range result.Items {
		relation := Relation{}
		err = dynamodbattribute.UnmarshalMap(i, &relation)
		relations = append(relations, relation)
		util.Check(err)
	}

	return nil, relations
}

func SaveRelation(relation Relation) bool {
	av, err := dynamodbattribute.MarshalMap(relation)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Relation"),
	}

	_, err = d.PutItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

func SaveRelations(relations []Relation) bool {
	res := true
	for _, relation := range relations {
		res = SaveRelation(relation)
	}
	return res
}

func DeleteRelation(relation Relation) error {
	tableName := "Relation"

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"UserName": {
				S: aws.String(relation.UserName),
			},
			"ProjectName": {
				S: aws.String(relation.ProjectName),
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
