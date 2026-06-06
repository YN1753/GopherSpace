package core_ai

import (
	"context"
	"log"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type RouterAgent struct {
	model model.BaseChatModel
}

func NewRouterAgent(model model.BaseChatModel) *RouterAgent {
	return &RouterAgent{model: model}
}

func (r *RouterAgent) Decide(ctx context.Context, input string) (string, error) {
	systemPrompt := `你是一个高精度的学生意图识别网关。你需要深度分析学生的输入，并严格从以下三个标签中选择一个返回。
绝对不要输出任何解释、引导语、标点符号、引号或多余的字，只需要输出标签名字本身。

可供选择的标签如下：
- 'reviewer'   ：当学生提交了具体的 Go 语言代码、寻求代码审查（Code Review）、交作业或者要求查找代码 Bug 时。
- 'consultant' ：当学生询问 Go 语言职业发展规划、校园课程大纲、期末考试重点，或者需要查询特定的 Go 标准库官方权威设计文档时。
- 'tutor'      ：当学生进行日常打招呼（如 Hi、你好）、纯闲聊，或者询问通用的语法概念（例如“什么是 channel”，但并未贴出具体复杂的业务代码）时。

记住：你的回答只能是 reviewer、consultant 或 tutor 中的一个，多一个字都会导致路由系统崩溃。`
	resp, err := r.model.Generate(ctx, []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(input),
	})
	if err != nil {
		return "tutor", err
	}
	decision := strings.TrimSpace(strings.ToLower(resp.Content))
	log.Printf("Router 决策: input=%q -> decision=%q", input, decision)
	if decision == "reviewer" || decision == "consultant" {
		return decision, nil
	}
	return "tutor", nil
}
