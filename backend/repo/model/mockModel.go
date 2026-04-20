package model

import (
	"context"
	"time"

	"github.com/ling/muse/entity"
)

const mockLastTime = 60
const mockInterval = 1

type MockModel struct {
}

func NewMockModel() LLMModel {
	return &MockModel{}
}

func (m *MockModel) GenerateContent(ctx context.Context, apiKey, model, proxyUrl string, preset entity.Preset, input []Message) ([]string, error) {
	panic("implement me")
}

func (m *MockModel) StreamGenerateContent(ctx context.Context, apiKey, model, proxyUrl string, preset entity.Preset, input []Message) <-chan StreamStruct {
	send := make(chan StreamStruct)
	startTime := time.Now().Unix()
	go func() {
		ticker := time.NewTicker(mockInterval * time.Second)
		for time.Now().Unix()-startTime <= mockLastTime {
			select {
			case <-ticker.C:
				send <- StreamStruct{
					Content: "我爱你",
					Index:   0,
					Done:    false,
					Error:   nil,
				}
			}
		}
		send <- StreamStruct{
			Content: "我爱你",
			Index:   0,
			Done:    true,
			Error:   nil,
		}
	}()
	return send
}
