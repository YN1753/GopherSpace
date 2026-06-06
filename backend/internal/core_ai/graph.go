package core_ai

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type ChatRequest struct {
	UserMessage string
	History     []*schema.Message
}

type routerOutput struct {
	Decision string
	Request  ChatRequest
}

type ChatGraph struct {
	runnable compose.Runnable[ChatRequest, *schema.Message]
}

func (cg *ChatGraph) Stream(ctx context.Context, req ChatRequest) (*schema.StreamReader[*schema.Message], error) {
	return cg.runnable.Stream(ctx, req)
}

func NewChatGraph(
	ctx context.Context,
	router *RouterAgent,
	reviewer *ReviewerAgent,
	consultant *ConsultantAgent,
	tutor *TutorAgent,
	chatModel model.BaseChatModel,
) (*ChatGraph, error) {
	g := compose.NewGraph[ChatRequest, *schema.Message]()

	// router 节点：调用 RouterAgent 做意图分类
	err := g.AddLambdaNode("router", compose.InvokableLambda(
		func(ctx context.Context, req ChatRequest) (routerOutput, error) {
			decision, err := router.Decide(ctx, req.UserMessage)
			return routerOutput{Decision: decision, Request: req}, err
		},
	))
	if err != nil {
		return nil, err
	}

	// 三个 prompt 构建节点
	err = g.AddLambdaNode("reviewer_prompt", compose.InvokableLambda(
		func(ctx context.Context, output routerOutput) ([]*schema.Message, error) {
			return reviewer.PrepareMessages(output.Request.UserMessage, output.Request.History), nil
		}))
	if err != nil {
		return nil, err
	}

	err = g.AddLambdaNode("consultant_prompt", compose.InvokableLambda(
		func(ctx context.Context, output routerOutput) ([]*schema.Message, error) {
			return consultant.PrepareMessages(output.Request.UserMessage, output.Request.History), nil
		}))
	if err != nil {
		return nil, err
	}

	err = g.AddLambdaNode("tutor_prompt", compose.InvokableLambda(
		func(ctx context.Context, output routerOutput) ([]*schema.Message, error) {
			return tutor.PrepareMessages(output.Request.UserMessage, output.Request.History), nil
		}))
	if err != nil {
		return nil, err
	}

	// model 节点：共享的 LLM
	err = g.AddChatModelNode("model", chatModel)
	if err != nil {
		return nil, err
	}

	// 三个 prompt 节点都连到 model
	g.AddEdge("reviewer_prompt", "model")
	g.AddEdge("consultant_prompt", "model")
	g.AddEdge("tutor_prompt", "model")
	g.AddEdge("model", compose.END)

	// 分支：router 根据决策路由到对应 prompt
	err = g.AddBranch("router", compose.NewGraphBranch(
		func(ctx context.Context, output any) (string, error) {
			result, ok := output.(routerOutput)
			if !ok {
				return "", fmt.Errorf("不支持的类型 %T", output)
			}
			switch result.Decision {
			case "reviewer":
				return "reviewer_prompt", nil
			case "consultant":
				return "consultant_prompt", nil
			default:
				return "tutor_prompt", nil
			}
		},
		map[string]bool{
			"reviewer_prompt":   true,
			"consultant_prompt": true,
			"tutor_prompt":      true,
		},
	))
	if err != nil {
		return nil, err
	}

	g.AddEdge(compose.START, "router")
	runnable, err := g.Compile(ctx)
	if err != nil {
		return nil, err
	}
	return &ChatGraph{runnable: runnable}, nil
}
