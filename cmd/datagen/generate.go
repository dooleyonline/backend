package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"

	itemsvc "github.com/dooleyonline/backend/internal/service/item"
	"google.golang.org/genai"
)

var (
	imageEnum = []string{
		"68328431-1a1d-4f35-8e1c-35a57421dc88",
		"689735dd-836c-4c60-b6c8-ac1e087f20d1",
		"b2dd76a2-33b1-4e98-b98c-c111a4d3810a",
		"d08f2271-44a7-4a98-8433-ae367d19cf4b",
		"dbcf4580-6544-4596-ad20-39cd56bff099",
		"b982c08b-342f-4ac5-930c-8ef871b40216",
		"90e385ca-2d76-4fd7-a9b2-a1d348654a76",
	}
)

func generate(ctx context.Context, categories []string) ([]itemsvc.MutationParams, error) {
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

	var items []itemsvc.MutationParams
	if err := json.Unmarshal([]byte(data), &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	return items, nil
}

func geminiConfig(categories []string) *genai.GenerateContentConfig {
	randImages := make([]string, 5)
	for i := range 5 {
		img := rand.IntN(len(imageEnum))
		randImages[i] = imageEnum[img]
	}

	return &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type:     genai.TypeArray,
			MinItems: &NumItems,
			MaxItems: &NumItems,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"name": {
						Description: "Name of the secondhand item",
						Type:        genai.TypeString,
					},
					"description": {
						Description: "Description of the item. Should address how old the item is and add more context to the item condition.",
						Type:        genai.TypeString,
					},
					"images": {
						Description: "Array of image IDs. This should be randomly chosen from the given enum.",
						Type:        genai.TypeArray,
						Items: &genai.Schema{
							Type:      genai.TypeString,
							Enum:      randImages,
							MinLength: genai.Ptr(int64(1)),
							MaxLength: genai.Ptr(int64(3)),
						},
					},
					"price": {
						Description: "Price of the item",
						Type:        genai.TypeNumber,
						Format:      "float",
					},
					"condition": {
						Description: "Condition of the item on a scale of 5",
						Type:        genai.TypeInteger,
						Default:     0,
						Maximum:     genai.Ptr(float64(5)),
					},
					"is_negotiable": {
						Description: "Whether the item price is negotiable or not.",
						Type:        genai.TypeBoolean,
					},
					"category": {
						Description: "Category of the item.",
						Type:        genai.TypeString,
						Enum:        categories,
					},
					"subcategory": {
						Description: "Subcategory of the item. This value must always be 'Other' for now.",
						Type:        genai.TypeString,
						Enum:        []string{"Other"},
					},
				},
				Required: []string{"name", "description", "images", "price", "condition", "is_negotiable", "category", "subcategory"},
			},
		},
	}
}
