package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dooleyonline/backend/internal/model"
	"google.golang.org/genai"
)

var (
	ImagePattern = `sample/([0-9]|10)\.webp`
	ConditionMax = 5.0
)

func generate(ctx context.Context, categories []string) ([]model.Item, error) {
	gemini, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	content, err := gemini.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-lite",
		genai.Text(Prompt),
		geminiConfig(categories),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	data := content.Text()

	var items []model.Item
	if err := json.Unmarshal([]byte(data), &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	return items, nil
}

func geminiConfig(categories []string) *genai.GenerateContentConfig {
	return &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type:     genai.TypeArray,
			MinItems: &NumItems,
			MaxItems: &NumItems,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"name":        {Type: genai.TypeString},
					"description": {Type: genai.TypeString},
					"images": {
						Type: genai.TypeArray,
						Items: &genai.Schema{
							Type:    genai.TypeString,
							Pattern: ImagePattern,
						},
					},
					"price": {
						Type:   genai.TypeNumber,
						Format: "float",
					},
					"condition": {
						Type:    genai.TypeInteger,
						Default: 0,
						Maximum: &ConditionMax,
					},
					"is_negotiable": {Type: genai.TypeBoolean},
					"category": {
						Type:   genai.TypeString,
						Format: "enum",
						Enum:   categories,
					},
					"subcategory": {
						Type: genai.TypeString,
						Enum: []string{"Other"},
					},
				},
				Required: []string{"name", "description", "images", "price", "condition", "is_negotiable", "category", "subcategory"},
			},
		},
	}
}
