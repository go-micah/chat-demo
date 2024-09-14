package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/go-micah/go-bedrock/providers"
)

func main() {

	// connect to AWS

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	svc := bedrockruntime.NewFromConfig(cfg)

	var messages []providers.AnthropicClaudeMessage

	// initial prompt
	fmt.Printf("Hi there. You can ask me stuff!\n")

	// create a TTY loop
	for {
		var chunks string

		// get user input
		prompt := stringPrompt(">")

		// check for special words

		if prompt == "quit\n" {
			os.Exit(0)
		}

		// prepare prompt

		accept := "*/*"
		contentType := "application/json"
		model := "anthropic.claude-3-haiku-20240307-v1:0"

		textPrompt := providers.AnthropicClaudeContent{
			Type: "text",
			Text: prompt,
		}

		message := providers.AnthropicClaudeMessage{
			Role: "user",
			Content: []providers.AnthropicClaudeContent{
				textPrompt,
			},
		}

		messages = append(messages, message)

		body := providers.AnthropicClaudeMessagesInvokeModelInput{
			Messages:      messages,
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

		// invoke with streaming response
		resp, err := svc.InvokeModelWithResponseStream(context.TODO(), &bedrockruntime.InvokeModelWithResponseStreamInput{
			Accept:      &accept,
			ModelId:     &model,
			ContentType: &contentType,
			Body:        bodyString,
		})
		if err != nil {
			log.Fatalf("error from Bedrock, %v", err)
		}

		var out providers.AnthropicClaudeMessagesInvokeModelOutput

		stream := resp.GetStream().Reader
		events := stream.Events()

		for {
			event := <-events
			if event != nil {
				if v, ok := event.(*types.ResponseStreamMemberChunk); ok {
					// v has fields
					err := json.Unmarshal([]byte(v.Value.Bytes), &out)
					if err != nil {
						log.Printf("unable to decode response:, %v", err)
						continue
					}
					if out.Type == "content_block_delta" {
						fmt.Printf("%v", out.Delta.Text)
						chunks = chunks + out.Delta.Text
					}
				} else if v, ok := event.(*types.UnknownUnionMember); ok {
					// catchall
					fmt.Print(v.Value)
				}
			} else {
				break
			}
		}
		stream.Close()

		if stream.Err() != nil {
			log.Fatalf("error from Bedrock, %v", stream.Err())
		}
		fmt.Println()

		textPrompt = providers.AnthropicClaudeContent{
			Type: "text",
			Text: chunks,
		}

		message = providers.AnthropicClaudeMessage{
			Role: "assistant",
			Content: []providers.AnthropicClaudeContent{
				textPrompt,
			},
		}

		messages = append(messages, message)
	}
}

// function to capture command line input
func stringPrompt(label string) string {

	var s string
	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}

	return s
}
