package model

import (
	"context"

	"github.com/ling/muse/entity"
	pb "github.com/ling/muse/gen/muse"
)

type LLMModel interface {
	GenerateContent(ctx context.Context, apiKey string, model string, preset entity.Preset, input []Message) ([]string, error)
	StreamGenerateContent(ctx context.Context, apiKey, model string, preset entity.Preset, input []Message) <-chan StreamStruct
}

type StreamStruct struct {
	Content string
	Index   int // 多个候选回复时，用于标记是第几个候选回复
	Done    bool
	Error   error
}

type Message struct {
	Role    pb.Role
	Content string
}
