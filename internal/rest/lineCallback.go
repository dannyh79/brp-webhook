package routes

import (
	"fmt"

	g "github.com/dannyh79/brp-webhook/internal/groups"
	s "github.com/dannyh79/brp-webhook/internal/services"
	"github.com/gin-gonic/gin"
)

type LineCallbackBody struct {
	Events []Event `json:"events"`
}

// Base event
type Event struct {
	Type    string           `json:"type"`
	Source  Source           `json:"source"`
	Message MessageEventBody `json:"message"`
}

type Source struct {
	Type    string `json:"type"`
	GroupId string `json:"groupId"`
	UserId  string `json:"userId"`
}

type MessageEvent struct {
	Event   `json:",inline"`
	Text    string           `json:"text"`
	Message MessageEventBody `json:"message"`
}

type MessageEventBody struct {
	Type       string `json:"type"`
	Text       string `json:"text"`
	ReplyToken string `json:"replyToken"`
}

type LeaveEvent struct {
	Event `json:",inline"`
}

const RegisterMyGroupMsg = "請好好靈修每日推播靈修內容到這"

type groupDto struct {
	*g.Group
	ReplyToken string
}

func msgEventsHandler(ctx *gin.Context) {
	defer ctx.Request.Body.Close()

	var b LineCallbackBody
	if err := ctx.ShouldBindJSON(&b); err != nil {
		fmt.Printf("Error in unmarshalling request body: %v", err)
		ctx.Next()
		return
	}

	var gs []*groupDto
	for _, e := range b.Events {
		if e.Type == "message" && e.Message.Text == RegisterMyGroupMsg && len(e.Message.ReplyToken) > 0 {
			gs = append(gs, &groupDto{Group: g.NewGroup(e.Source.GroupId), ReplyToken: e.Message.ReplyToken})
		}
	}

	ctx.Set("groups", gs)

	ctx.Next()
}

func groupRegistrationHandler(sCtx *s.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gsIf, exists := ctx.Get("groups")
		if !exists {
			ctx.Next()
			return
		}

		gs, ok := gsIf.([]*groupDto)
		if !ok {
			ctx.Next()
			return
		}

		var registered []*groupDto
		for _, g := range gs {
			if err := sCtx.RegistrationService.Execute(g.Group); err != nil {
				fmt.Printf("Error in registering group: %v", err)
			} else {
				registered = append(registered, g)
			}
		}

		ctx.Set("group", registered)

		ctx.Next()
	}
}

func replyHandler(sCtx *s.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gsIf, exists := ctx.Get("groups")
		if !exists {
			ctx.Next()
			return
		}

		gs, ok := gsIf.([]*groupDto)
		if !ok {
			ctx.Next()
			return
		}

		for _, g := range gs {
			if err := sCtx.ReplyService.Execute(g.ReplyToken); err != nil {
				fmt.Printf("Error in replying to completed registration for group %v via LINE: %v", g.Id, err)
			}
		}

		ctx.Next()
	}
}
