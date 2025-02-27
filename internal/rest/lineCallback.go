package routes

import (
	"log"

	g "github.com/dannyh79/brp-webhook/internal/groups"
	"github.com/dannyh79/brp-webhook/internal/sentry"
	s "github.com/dannyh79/brp-webhook/internal/services"
	"github.com/gin-gonic/gin"
)

// FIXME: Use LINE SDK's discriminated union event types here
type LineCallbackBody struct {
	Events []Event `json:"events"`
}

// Base event
type Event struct {
	Type       string           `json:"type"`
	Source     Source           `json:"source"`
	Message    MessageEventBody `json:"message"`
	ReplyToken string           `json:"replyToken"`
}

type Source struct {
	Type    string `json:"type"`
	GroupId string `json:"groupId"`
	UserId  string `json:"userId"`
}

// https://developers.line.biz/en/reference/messaging-api/#message-event
type MessageEvent struct {
	Event   `json:",inline"`
	Text    string           `json:"text"`
	Message MessageEventBody `json:"message"`
}

type MessageEventBody struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// https://developers.line.biz/en/reference/messaging-api/#join-event
type JoinEvent struct {
	Event `json:",inline"`
}

// https://developers.line.biz/en/reference/messaging-api/#leave-event
type LeaveEvent struct {
	Event `json:",inline"`
}

func lineEventsHandler(sCtx *s.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sentry.TagBy(ctx, "lineEventsHandler")

		var b LineCallbackBody
		if err := ctx.ShouldBindJSON(&b); err != nil {
			log.Printf("Error in unmarshalling request body: %v", err)
			ctx.Next()
			return
		}

		for _, e := range b.Events {
			// NOTE: Current flow targets group-related events only
			if e.Source.GroupId == "" {
				continue
			}

			switch e.Type {
			case "join":
				handleJoinEvent(sCtx, e)
			case "message":
				handleMessageEvent(sCtx, e)
			case "leave":
				handleLeaveEvent(sCtx, e)
			}
		}

		ctx.Next()
	}
}

func handleJoinEvent(sCtx *s.ServiceContext, e Event) {
	if len(e.ReplyToken) == 0 {
		return
	}

	g := s.NewGroupDto(g.NewGroup(e.Source.GroupId), e.ReplyToken)
	if err := sCtx.WelcomeService.Execute(g); err != nil {
		log.Printf("Error in welcoming group %v via LINE: %v", g.Id, err)
	}
}

func handleMessageEvent(sCtx *s.ServiceContext, e Event) {
	if e.Message.Text != g.MsgRegisterMyGroup || len(e.ReplyToken) == 0 {
		return
	}

	g := s.NewGroupDto(g.NewGroup(e.Source.GroupId), e.ReplyToken)

	if err := sCtx.RegistrationService.Execute(g); err != nil {
		if err == s.ErrorGroupAlreadyRegistered {
			g.WasRegistered = true
		} else {
			log.Printf("Error in registering group: %v", err)
			return
		}
	}

	if err := sCtx.ReplyService.Execute(g); err != nil {
		log.Printf("Error in replying to completed registration for group %v via LINE: %v", g.Id, err)
	}
}

func handleLeaveEvent(sCtx *s.ServiceContext, e Event) {
	g := s.NewGroupDto(g.NewGroup(e.Source.GroupId), "")
	if err := sCtx.UnlistService.Execute(g); err != nil {
		log.Printf("Error in unlisting group: %v", err)
	}
}
