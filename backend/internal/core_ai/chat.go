package core_ai

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/claude"
	"github.com/cloudwego/eino/schema"
)

type CoreAiService struct {
	Model *claude.ChatModel
}

func NewCoreAiService(model *claude.ChatModel) *CoreAiService {
	return &CoreAiService{Model: model}
}

func (c *CoreAiService) ChatWithModel(ctx context.Context, message []*schema.Message) error {
	stream, err := c.Model.Stream(ctx, message)
	if err != nil {
		return errors.New("请求大模型失败")
	}
	defer stream.Close()
	for {
		chunk, err := stream.Recv()
		if err != nil {
			break
		}
		fmt.Println(chunk)
	}
	return nil
}
