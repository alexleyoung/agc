package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/viper"

	"google.golang.org/genai"
)

var client *genai.Client

func CreateChat(ctx context.Context, model string, history []*genai.Content) (*genai.Chat, error) {
	var err error
	if client == nil {
		client, err = genai.NewClient(ctx, &genai.ClientConfig{
			APIKey:  viper.Get("GEMINI_API_KEY").(string),
			Backend: genai.BackendGeminiAPI,
		})
		if err != nil {
			log.Printf("Error generating AI client: %v", err)
			return nil, err
		}
	}

	chat, err := client.Chats.Create(ctx, model, config, history)
	if err != nil {
		log.Printf("Error creating chat: %v", err)
		return nil, err
	}
	return chat, nil
}

func Query(ctx context.Context, chat *genai.Chat, prompt string) (*genai.GenerateContentResponse, error) {
	msg := genai.Part{Text: prompt}
	for range MAX_STEPS {
		// prompt model
		result, err := chat.SendMessage(ctx, msg)
		if err != nil {
			log.Printf("Error getting model response: %v", err)
			return &genai.GenerateContentResponse{}, err
		}

		// check for function calls
		fns := result.FunctionCalls()
		if len(fns) > 0 {
			fn := fns[0]

			args, err := json.Marshal(fn.Args)
			if err != nil {
				return &genai.GenerateContentResponse{}, fmt.Errorf("failed to marshal args: %w", err)
			}

			fmt.Printf("Calling function: %s\n", fn.Name)
			out, err := executeFunctionCall(ctx, fn.Name, args)
			if err != nil {
				log.Printf("Error executing function: %v", err)
				return &genai.GenerateContentResponse{}, err
			}

			// set mesage as function response
			msg = genai.Part{
				FunctionResponse: &genai.FunctionResponse{
					Name: fn.Name,
					Response: map[string]any{
						"output": out,
					},
				},
			}
			continue
		}
		return result, nil
	}
	return &genai.GenerateContentResponse{}, fmt.Errorf("Max steps reached without resolution")
}
