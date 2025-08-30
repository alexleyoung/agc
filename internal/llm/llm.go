package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/viper"

	"google.golang.org/genai"
)

var TEMPERATURE float32 = 0

func Query(ctx context.Context, model string, history []*genai.Content, prompt string) (*genai.GenerateContentResponse, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  viper.Get("GEMINI_API_KEY").(string),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Printf("Error generating AI client: %v", err)
		return &genai.GenerateContentResponse{}, err
	}

	history = append(history, &genai.Content{
		Role:  "user",
		Parts: []*genai.Part{{Text: prompt}},
	})

	var result *genai.GenerateContentResponse
	for range MAX_STEPS {
		// prompt model
		result, err = client.Models.GenerateContent(ctx, model, history, &genai.GenerateContentConfig{
			SystemInstruction: genai.NewContentFromText(SYSTEM_PROMPT, genai.RoleUser),
			Tools: []*genai.Tool{
				{FunctionDeclarations: functionDeclarations},
			},
			Temperature: &TEMPERATURE,
		})
		if err != nil {
			log.Printf("Error getting model response: %v", err)
			return &genai.GenerateContentResponse{}, err
		}

		history = append(history, result.Candidates[0].Content)

		// check for function calls
		fns := result.FunctionCalls()
		if len(fns) > 0 {
			fn := fns[0]

			args, err := json.Marshal(fn.Args)
			if err != nil {
				return &genai.GenerateContentResponse{}, fmt.Errorf("failed to marshal args: %w", err)
			}

			out, err := executeFunctionCall(ctx, fn.Name, args)
			if err != nil {
				log.Printf("Error executing function: %v", err)
				return &genai.GenerateContentResponse{}, err
			}

			history = append(history, &genai.Content{
				Role:  "model",
				Parts: []*genai.Part{{Text: fn.Name + " results: "}, {Text: out}},
			})
			continue
		}
		return result, nil
	}
	return &genai.GenerateContentResponse{}, fmt.Errorf("Max steps reached without resolution")
}

func executeFunctionCall(ctx context.Context, name string, argsJSON []byte) (string, error) {
	switch name {
	case "create_event":
		return create_event_call(ctx, argsJSON)
	case "quick_add_event":
		return quick_add_event_call(ctx, argsJSON)
	case "list_calendars":
		return list_calendars_call(ctx)
	case "get_current_time":
		return get_current_time_call(ctx)
	}
	return "", fmt.Errorf("Unknown function: %s", name)
}
