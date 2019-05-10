package dynamo

import (
	"auth/util"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetProjects() (error, []Project) {
	var projects []Project
	tableName := "Project"

	params := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	// Make the DynamoDB Query API call
	result, err := d.Scan(params)
	util.Check(err)

	for _, i := range result.Items {
		project := Project{}
		err = dynamodbattribute.UnmarshalMap(i, &project)
		projects = append(projects, project)
		util.Check(err)
	}

	return nil, projects
}

func SaveProject(project Project) bool {
	av, err := dynamodbattribute.MarshalMap(project)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Project"),
	}

	_, err = d.PutItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

func SaveProjects(projects []Project) bool {
	res := true
	for _, project := range projects {
		res = SaveProject(project)
	}
	return res
}

func UpdateProject(project Project) error {
	tableName := "Project"

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":p": {
				S: aws.String(""),
			},
			":r": {
				S: aws.String(""),
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ProjectName": {
				S: aws.String(project.ProjectName),
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

func DeleteProject(project Project) error {
	tableName := "Project"

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ProjectName": {
				S: aws.String(project.ProjectName),
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

func GetProjectEnvs(projectName string) (error, Project) {
	tableName := "Project"

	result, err := d.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ProjectName": {
				S: aws.String(projectName),
			},
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		return err, Project{}
	}

	item := Project{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return nil, item
}

func UpdateProjectEnvs(project Project) error {
	tableName := "Project"

	av, err := dynamodbattribute.MarshalList(project.EnvList)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":l": {
				L: av,
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ProjectName": {
				S: aws.String(project.ProjectName),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set EnvList = :l"),
	}

	_, err = d.UpdateItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
