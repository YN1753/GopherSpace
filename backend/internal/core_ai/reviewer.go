package core_ai

import (
	"context"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type ReviewerAgent struct {
	Model model.BaseChatModel
}

func NewReviewerAgent(model model.BaseChatModel) *ReviewerAgent {
	return &ReviewerAgent{
		Model: model,
	}
}

func (r *ReviewerAgent) PrepareMessages(userMsg string, history []*schema.Message) []*schema.Message {
	msg := []*schema.Message{
		schema.SystemMessage(
			`你是一位资深的go代码审查专家。职责：检查程序的bug，内存泄露，竞态条件，指出不符合Go惯用用法的写法，给出具体的改进建议。
			语气：专业，直接，建设性。`),
	}
	if len(history) > 0 {
		msg = append(msg, history...)
	}
	msg = append(msg, schema.UserMessage(userMsg))
	return msg
}

func (r *ReviewerAgent) Stream(ctx context.Context, userMsg string, history []*schema.Message) (*schema.StreamReader[*schema.Message], error) {
	msg := r.PrepareMessages(userMsg, history)
	return r.Model.Stream(ctx, msg)
}
