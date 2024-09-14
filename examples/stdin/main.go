package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/go-micah/go-bedrock/providers"
	"github.com/mattn/go-isatty"
)

func main() {

	// get prompt from command line args

	prompt := os.Args[1]

	// capture stdin if there is anything

	var document string

	if isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd()) {
		// do nothing
	} else {
		stdin, err := io.ReadAll(os.Stdin)

		if err != nil {
			panic(err)
		}
		document = string(stdin)
	}

	if document != "" {
		document = "<document>\n\n" + document + "\n\n</document>\n\n"
		prompt = document + prompt
	}

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

	textPrompt := providers.AnthropicClaudeContent{
		Type: "text",
		Text: prompt,
	}

	content := []providers.AnthropicClaudeContent{
		textPrompt,
	}

	body := providers.AnthropicClaudeMessagesInvokeModelInput{
		Messages: []providers.AnthropicClaudeMessage{
			{
				Role:    "user",
				Content: content,
			},
		},
		MaxTokens:     2000,
		TopP:          0.999,
		TopK:          int(250),
		Temperature:   1,
		StopSequences: []string{},
	}

	bodyString, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("unable to marshal body: %v", err)
	}

	// invoke Amazon Bedrock

	resp, err := svc.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		Accept:      &accept,
		ModelId:     &model,
		ContentType: &contentType,
		Body:        bodyString,
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
