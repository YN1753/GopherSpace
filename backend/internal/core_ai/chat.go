package core_ai

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type CoreAiService struct {
	Model model.BaseChatModel
}

func NewCoreAiService(model model.BaseChatModel) *CoreAiService {
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
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(chunk)
	}
	return nil
}
