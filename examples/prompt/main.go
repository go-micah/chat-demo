package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/go-micah/go-bedrock/providers"
)

func main() {

	// get prompt from command line args

	prompt := os.Args[1]

	// connect to AWS

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	svc := bedrockruntime.NewFromConfig(cfg)

	// prepare prompt

	accept := "*/*"
	contentType := "application/json"
	model := "anthropic.claude-3-haiku-20240307-v1:0"
	body := "{\"messages\":[{\"role\":\"user\",\"content\":[{\"type\":\"text\",\"text\":\"" + prompt + "\"}]}],\"anthropic_version\":\"bedrock-2023-05-31\",\"max_tokens\":2000,\"temperature\":1,\"top_k\":250,\"top_p\":0.999,\"stop_sequences\":[]}"

	// invoke Amazon Bedrock

	resp, err := svc.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		Accept:      &accept,
		ModelId:     &model,
		ContentType: &contentType,
		Body:        []byte(body),
	})
	if err != nil {
		log.Fatalf("error from Bedrock, %v", err)
	}

	// print response to stdout

	var out providers.AnthropicClaudeMessagesInvokeModelOutput

	err = json.Unmarshal(resp.Body, &out)
	if err != nil {
		log.Fatalf("unable to unmarshal response from Bedrock: %v", err)
	}

	fmt.Println(out.Content[0].Text)

}
