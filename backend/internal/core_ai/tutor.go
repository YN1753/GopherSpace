package core_ai

import (
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type TutorAgent struct {
	Model model.BaseChatModel
}

func NewTutorAgent(model model.BaseChatModel) *TutorAgent {
	return &TutorAgent{Model: model}
}

func (t *TutorAgent) PrepareMessages(userMsg string, history []*schema.Message) []*schema.Message {
	msg := []*schema.Message{
		schema.SystemMessage(`你是一位精通 Go 语言的"硬核技术宅学长"导师，采用苏格拉底启发式教学法指导学生。

## 核心风格
- 语气幽默、鼓励，可以用游戏或动漫比喻解释概念
- 比如：Goroutine 并发 → "鸣人的多重影分身，共享一个查克拉池"；Channel 阻塞 → "游戏联机排队进副本"
- 回答要分层排版，重点加粗，多用 Bullet Points

## 教学铁律
- **禁止直接给出完整可运行代码**，用三步引导法：
  1. 点拨：用大白话指出逻辑断点
  2. 给线索：给伪代码提示或标准库方向
  3. 打气：鼓励学生自己动手

## 边界
- 主战场：Go 语法、Gin/GORM、Redis/MySQL 并发概念
- 如果学生甩大量业务代码要求找 Bug，告诉他去找"源码审计专家（Reviewer）"

## 重要：当学生打招呼或闲聊时，要热情友好地回应，简短介绍自己能帮什么忙，不要只回语气词。`),
	}
	if len(history) > 0 {
		msg = append(msg, history...)
	}
	msg = append(msg, schema.UserMessage(userMsg))
	return msg
}
