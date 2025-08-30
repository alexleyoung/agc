package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/viper"

	"google.golang.org/genai"
)

func Query(ctx context.Context, model string, history []*genai.Content, prompt string) (*genai.GenerateContentResponse, []*genai.Content, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  viper.Get("GEMINI_API_KEY").(string),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Printf("Error generating AI client: %v", err)
		return &genai.GenerateContentResponse{}, history, err
	}

	var result *genai.GenerateContentResponse
	for range MAX_STEPS {
		// prompt model
		chat, err := client.Chats.Create(ctx, model, CONFIG, history)

		msg := genai.Part{Text: prompt}
		result, err = chat.SendMessage(ctx, msg)
		if err != nil {
			log.Printf("Error getting model response: %v", err)
			return &genai.GenerateContentResponse{}, history, err
		}

		history = append(history, &genai.Content{Role: "user", Parts: []*genai.Part{&msg}})
		history = append(history, result.Candidates[0].Content)

		// check for function calls
		fns := result.FunctionCalls()
		if len(fns) > 0 {
			fn := fns[0]

			args, err := json.Marshal(fn.Args)
			if err != nil {
				return &genai.GenerateContentResponse{}, history, fmt.Errorf("failed to marshal args: %w", err)
			}

			out, err := executeFunctionCall(ctx, fn.Name, args)
			if err != nil {
				log.Printf("Error executing function: %v", err)
				return &genai.GenerateContentResponse{}, history, err
			}

			history = append(history, &genai.Content{
				Role:  "model",
				Parts: []*genai.Part{{Text: fn.Name + " results: "}, {Text: out}},
			})
			continue
		}
		return result, history, nil
	}
	return &genai.GenerateContentResponse{}, history, fmt.Errorf("Max steps reached without resolution")
}
