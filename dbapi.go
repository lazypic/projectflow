package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	partitionKey = "ID"
)

func tableStruct(tableName string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(partitionKey),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(partitionKey),
				KeyType:       aws.String("HASH"),
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest), // ondemand
		TableName:   aws.String(tableName),
	}
}

func validTable(db dynamodb.DynamoDB, tableName string) bool {
	input := &dynamodb.ListTablesInput{}
	isTableName := false
	// 한번에 최대 100개의 테이블만 가지고 올 수 있다.
	// 한 리전에 최대 256개의 테이블이 존재할 수 있다.
	// https://docs.aws.amazon.com/ko_kr/amazondynamodb/latest/developerguide/Limits.html
	for {
		result, err := db.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					fmt.Fprintf(os.Stderr, "%s %s\n", dynamodb.ErrCodeInternalServerError, err.Error())
				default:
					fmt.Fprintf(os.Stderr, "%s\n", aerr.Error())
				}
			} else {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
			return false
		}

		for _, n := range result.TableNames {
			if *n == tableName {
				isTableName = true
				break
			}
		}
		if isTableName {
			break
		}
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}
	return isTableName
}

func hasItem(db dynamodb.DynamoDB, tableName string, primarykey string) (bool, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			partitionKey: {
				S: aws.String(primarykey),
			},
		},
	}
	result, err := db.GetItem(input)
	if err != nil {
		return false, err
	}
	if result.Item == nil {
		return false, nil
	}
	return true, nil
}

// AddProject 는 사용자를 추가하는 함수이다.
func AddProject(db dynamodb.DynamoDB) error {
	hasBool, err := hasItem(db, *flagTable, *flagID)
	if err != nil {
		return err
	}
	if hasBool {
		return errors.New("The data already exists. Can not add data")
	}
	p := Project{
		ID:            *flagID,
		Budget:        *flagBudget,
		StartDate:     *flagStartDate,
		UpdateDate:    *flagUpdateDate,
		EndDate:       *flagEndDate,
		ProjectStatus: *flagProjectStatus,
		MonetaryUnit:  *flagMonetaryUnit,
	}

	dynamodbJSON, err := dynamodbattribute.MarshalMap(p)
	if err != nil {
		return err
	}

	data := &dynamodb.PutItemInput{
		Item:      dynamodbJSON,
		TableName: aws.String(*flagTable),
	}
	_, err = db.PutItem(data)
	if err != nil {
		return err
	}
	return nil
}

// SetProject 는 프로젝트 자료구조를 수정하는 함수이다.
func SetProject(db dynamodb.DynamoDB) error {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(*flagTable),
		Key: map[string]*dynamodb.AttributeValue{
			partitionKey: {
				S: aws.String(*flagID),
			},
		},
	}
	result, err := db.GetItem(input)
	if err != nil {
		return err
	}
	p := Project{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &p)
	if err != nil {
		return err
	}
	if *flagBudget != 0 && p.Budget != *flagBudget {
		p.Budget = *flagBudget
	}
	if *flagStartDate != "" && p.StartDate != *flagStartDate {
		p.StartDate = *flagStartDate
	}
	if *flagEndDate != "" && p.EndDate != *flagEndDate {
		p.EndDate = *flagEndDate
	}
	if *flagProjectStatus != "" && p.ProjectStatus != *flagProjectStatus {
		p.ProjectStatus = *flagProjectStatus
	}
	if *flagMonetaryUnit != "KRW" && p.MonetaryUnit != *flagMonetaryUnit {
		p.MonetaryUnit = *flagMonetaryUnit
	}
	p.UpdateDate = *flagUpdateDate
	dynamodbJSON, err := dynamodbattribute.MarshalMap(p)
	if err != nil {
		return err
	}
	data := &dynamodb.PutItemInput{
		Item:      dynamodbJSON,
		TableName: aws.String(*flagTable),
	}
	_, err = db.PutItem(data)
	if err != nil {
		return err
	}
	return nil
}

// RmProject 는 프로젝트 자료구조를 사용자를 삭제하는 함수이다.
func RmProject(db dynamodb.DynamoDB) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(*flagTable),
		Key: map[string]*dynamodb.AttributeValue{
			partitionKey: {
				S: aws.String(*flagID),
			},
		},
	}
	_, err := db.DeleteItem(input)
	if err != nil {
		return err
	}
	return nil
}

// GetProjects 는 사용자를 가지고오는 함수이다.
func GetProjects(db dynamodb.DynamoDB, word string) error {
	proj := expression.NamesList(
		expression.Name("ID"),
		expression.Name("Budget"),
		expression.Name("MonetaryUnit"),
		expression.Name("StartDate"),
		expression.Name("EndDate"),
		expression.Name("UpdateDate"),
		expression.Name("ProjectStatus"),
	)
	filt1 := expression.Name("ID").Contains(word)
	filt2 := expression.Name("StartDate").Contains(word)
	filt3 := expression.Name("EndDate").Contains(word)

	expr, err := expression.NewBuilder().
		WithFilter(filt1.Or(filt2).Or(filt3)).
		WithProjection(proj).
		Build()
	if err != nil {
		return err
	}
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(*flagTable),
	}
	result, err := db.Scan(params)
	if err != nil {
		return err
	}
	for _, i := range result.Items {
		p := Project{}
		err = dynamodbattribute.UnmarshalMap(i, &p)
		if err != nil {
			return err
		}
		fmt.Println(p)
	}
	return nil
}
