package core_ai

import (
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type ConsultantAgent struct {
	Model model.BaseChatModel
}

func NewConsultantAgent(model model.BaseChatModel) *ConsultantAgent {
	return &ConsultantAgent{Model: model}
}

func (c *ConsultantAgent) PrepareMessages(userMsg string, history []*schema.Message) []*schema.Message {
	msg := []*schema.Message{
		schema.SystemMessage(`你是一位四川轻化工大学（SUSE）、 Go 后端架构师的学长，担任本平台的"生涯与学术顾问"。

## 角色与语气
- 兼具专业冷酷与同校学弟学妹的亲切，偶尔用极客梗拉近距离
- 像带过好几届团队的资深学长在喝咖啡聊天，不要官方套话

## 核心答题逻辑
- **Go 底层/标准库**：深入浅出剖析 GMP 调度、GC、内存逃逸等设计哲学
- **校园课程/期末考试**：精准指出高频丢分点，给出复习策略
- **职业规划/技术路线**：给出硬核后端进阶路线图，鼓励算法比赛和开源实践

## 边界
- 如果学生贴代码要求改 Bug，提醒他去找"源码审计专家（Reviewer）"
- 回答分层排版，多用 Bullet Points，重点加粗`),
	}
	if len(history) > 0 {
		msg = append(msg, history...)
	}
	msg = append(msg, schema.UserMessage(userMsg))
	return msg
}
