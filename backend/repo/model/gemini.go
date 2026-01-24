package model

import (
	"context"
	"sync"

	"connectrpc.com/connect"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/entity"
	pb "github.com/ling/muse/gen/muse"
	"google.golang.org/genai"
)

var DefaultSafeSetting = []*genai.SafetySetting{
	{
		Category:  genai.HarmCategoryHarassment,
		Threshold: genai.HarmBlockThresholdOff,
	},
	{
		Category:  genai.HarmCategoryHateSpeech,
		Threshold: genai.HarmBlockThresholdOff,
	},
	{
		Category:  genai.HarmCategorySexuallyExplicit,
		Threshold: genai.HarmBlockThresholdOff,
	},
	{
		Category:  genai.HarmCategoryDangerousContent,
		Threshold: genai.HarmBlockThresholdOff,
	},
}

type Gemini struct {
	clients sync.Map
}

func NewGemini() LLMModel {
	return &Gemini{
		clients: sync.Map{},
	}
}

func (g *Gemini) getClient(ctx context.Context, apiKey string) *genai.Client {
	if client, ok := g.clients.Load(apiKey); ok {
		return client.(*genai.Client)
	}
	client, _ := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	g.clients.Store(apiKey, client)
	return client
}

func (g *Gemini) GenerateContent(ctx context.Context, apiKey string, model string, preset entity.Preset, input []Message) ([]string, error) {
	if apiKey == "" {
		return nil, errs.NewStandardf(connect.CodeInvalidArgument, "API Key不能为空")
	}
	histories, genConfig := g.buildReq(apiKey, preset, input)
	client := g.getClient(ctx, apiKey)
	rsp, err := client.Models.GenerateContent(ctx, model, histories, genConfig)
	if err != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "生成内容失败: %v", err)
	}
	result := make([]string, preset.CandidateCount)
	for i, candidate := range rsp.Candidates {
		for _, part := range candidate.Content.Parts {
			if part.Text != "" {
				if part.Thought {
					continue
				}
				result[i] += part.Text
			}
		}
	}
	return result, nil
}

func (g *Gemini) StreamGenerateContent(ctx context.Context, apiKey, model string, preset entity.Preset,
	input []Message) <-chan StreamStruct {

	resultChan := make(chan StreamStruct)
	defer close(resultChan)
	if apiKey == "" {
		resultChan <- StreamStruct{
			Done:  true,
			Error: errs.NewStandardf(connect.CodeInvalidArgument, "API Key不能为空"),
		}
		return resultChan
	}

	go func() {
		histories, genConfig := g.buildReq(apiKey, preset, input)
		client := g.getClient(ctx, apiKey)
		for content, err := range client.Models.GenerateContentStream(ctx, model, histories, genConfig) {
			if err != nil {
				resultChan <- StreamStruct{
					Done:  true,
					Error: errs.NewStandardf(connect.CodeInvalidArgument, "API Key不能为空"),
				}
			}
			for i, candidate := range content.Candidates {
				for _, part := range candidate.Content.Parts {
					if part.Text != "" {
						if part.Thought {
							continue
						}
						resultChan <- StreamStruct{
							Content: part.Text,
							Index:   i,
						}
					}
				}
			}
			resultChan <- StreamStruct{
				Done: true,
			}
		}
	}()

	return resultChan
}

func (g *Gemini) buildReq(apiKey string, preset entity.Preset, input []Message) (
	histories []*genai.Content, genConfig *genai.GenerateContentConfig) {

	histories = make([]*genai.Content, len(input))
	for _, message := range input {
		role := genai.RoleUser
		if message.Role == pb.Role_Assistant {
			role = genai.RoleModel
		}
		histories = append(histories, genai.NewContentFromText(message.Content, genai.Role(role)))
	}

	topK := float32(preset.TopK)
	genConfig = &genai.GenerateContentConfig{
		Temperature:      &preset.Temperature,
		TopP:             &preset.TopP,
		TopK:             &topK,
		CandidateCount:   int32(preset.CandidateCount),
		PresencePenalty:  &preset.PresencePenalty,
		FrequencyPenalty: &preset.FrequencyPenalty,
		SafetySettings:   DefaultSafeSetting,
	}
	return
}
