package model

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/entity"
	pb "github.com/ling/muse/gen/muse"
	log "github.com/sirupsen/logrus"
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

func (g *Gemini) GenerateContent(ctx context.Context, apiKey, model, proxyUrl string, preset entity.Preset, input []Message) ([]string, error) {
	if apiKey == "" {
		return nil, errs.NewStandardf(connect.CodeInvalidArgument, "API Key不能为空")
	}
	histories, genConfig := g.buildReq(preset, input, proxyUrl)
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

func (g *Gemini) StreamGenerateContent(ctx context.Context, apiKey, model, proxyUrl string, preset entity.Preset, input []Message) <-chan StreamStruct {

	resultChan := make(chan StreamStruct)
	if apiKey == "" {
		go func() {
			resultChan <- StreamStruct{
				Done:  true,
				Error: errs.NewStandardf(connect.CodeInvalidArgument, "API Key不能为空"),
			}
		}()
		return resultChan
	}

	go func() {
		histories, genConfig := g.buildReq(preset, input, proxyUrl)
		client := g.getClient(ctx, apiKey)
		var globalErr error
		for rsp, err := range client.Models.GenerateContentStream(ctx, model, histories, genConfig) {
			if err != nil {
				log.Errorf("gemini输出失败: %v", err)
				select {
				case resultChan <- StreamStruct{
					Done:  true,
					Error: errs.NewStandardf(connect.CodeInternal, "发送消息失败: %v", err),
				}:
				case <-ctx.Done():
					log.Warn("ctx的时间已到!")
				}
			}
			if rsp == nil {
				log.Warn("Gemini回复的消息为空")
				continue
			}
			rspJson, err := json.Marshal(rsp)
			if err != nil {
				log.Warn("序列化gemini 回复失败")
			}
			log.Debugf("gemini回复: %s", string(rspJson))
			for i, candidate := range rsp.Candidates {
				if candidate == nil {
					log.Warn("Gemini回复的候选词为空")
					continue
				}
				if candidate.FinishReason != "" && candidate.FinishReason != genai.FinishReasonStop {
					log.Errorf("Gemini非正常输出: %s", candidate.FinishReason)
					globalErr = errs.Newf(pb.ErrCode_AIExceptionOutput,
						"Gemini非正常输出: %s", candidate.FinishReason)
					break
				}
				for _, part := range candidate.Content.Parts {
					if part == nil {
						log.Warn("Gemini的candidate.Content.Parts为空")
						continue
					}
					if part.Text != "" {
						if part.Thought {
							continue
						}
						select {
						case resultChan <- StreamStruct{
							Content: part.Text,
							Index:   i,
						}:
						case <-ctx.Done():
							log.Warn("ctx的时间已到!")
							return
						}
					}
				}
			}
		}
		log.Info("gemini输出完毕!")

		// 检测是否是由于上下文取消导致的结束
		if err := ctx.Err(); err != nil {
			log.Warnf("gemini输出中断: %v", err)
			select {
			case resultChan <- StreamStruct{
				Done:  true,
				Error: err,
			}:
			case <-time.After(time.Second): // 防止阻塞
			}
			return
		}

		select {
		case resultChan <- StreamStruct{
			Done:  true,
			Error: globalErr,
		}:
		case <-ctx.Done():
			log.Warn("ctx的时间已到!")
		}
	}()

	return resultChan
}

func (g *Gemini) buildReq(preset entity.Preset, input []Message, proxyUrl string) (
	histories []*genai.Content, genConfig *genai.GenerateContentConfig) {

	histories = make([]*genai.Content, 0, len(input))
	for _, message := range input {
		role := genai.RoleUser
		if message.Role == pb.Role_Assistant {
			role = genai.RoleModel
		}
		histories = append(histories, genai.NewContentFromText(message.Content, genai.Role(role)))
	}

	topK := float32(preset.TopK)
	genConfig = &genai.GenerateContentConfig{
		Temperature:    &preset.Temperature,
		TopP:           &preset.TopP,
		TopK:           &topK,
		CandidateCount: int32(preset.CandidateCount),
		//PresencePenalty:  &preset.PresencePenalty,
		//FrequencyPenalty: &preset.FrequencyPenalty,
		SafetySettings: DefaultSafeSetting,
		HTTPOptions:    &genai.HTTPOptions{BaseURL: proxyUrl},
	}
	return
}
