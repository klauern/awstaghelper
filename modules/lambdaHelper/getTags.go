package lambdaHelper

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"log"
	"strings"
)

// getInstances return all lambdas from specified region
func getInstances(session session.Session) []*lambda.FunctionConfiguration {
	client := lambda.New(&session)
	input := &lambda.ListFunctionsInput{}

	var result []*lambda.FunctionConfiguration

	err := client.ListFunctionsPages(input,
		func(page *lambda.ListFunctionsOutput, lastPage bool) bool {
			result = append(result, page.Functions...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get instances", err)
		return nil
	}
	return result
}

// ParseLambdasTags parse output from getInstances and return arn and specified tags.
func ParseLambdasTags(tagsToRead string, session session.Session) [][]string {
	instancesOutput := getInstances(session)
	client := lambda.New(&session)
	var rows [][]string
	headers := []string{"Arn"}
	headers = append(headers, strings.Split(tagsToRead, ",")...)
	rows = append(rows, headers)
	for _, lambdaOutput := range instancesOutput {
		lambdaTags, err := client.ListTags(&lambda.ListTagsInput{Resource: lambdaOutput.FunctionArn})
		if err != nil {
			fmt.Println("Not able to get lambda tags", err)
		}
		tags := map[string]string{}
		for key, value := range lambdaTags.Tags {
			tags[key] = *value
		}

		var resultTags []string
		for _, key := range strings.Split(tagsToRead, ",") {
			resultTags = append(resultTags, tags[key])
		}
		rows = append(rows, append([]string{*lambdaOutput.FunctionArn}, resultTags...))
	}
	return rows
}
