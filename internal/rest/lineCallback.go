package routes

import (
	"log"

	g "github.com/dannyh79/brp-webhook/internal/groups"
	s "github.com/dannyh79/brp-webhook/internal/services"
	"github.com/gin-gonic/gin"
)

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

type MessageEvent struct {
	Event   `json:",inline"`
	Text    string           `json:"text"`
	Message MessageEventBody `json:"message"`
}

type MessageEventBody struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type LeaveEvent struct {
	Event `json:",inline"`
}

const RegisterMyGroupMsg = "請好好靈修每日推播靈修內容到這"

func lineEventsHandler(sCtx *s.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var b LineCallbackBody
		if err := ctx.ShouldBindJSON(&b); err != nil {
			log.Printf("Error in unmarshalling request body: %v", err)
			ctx.Next()
			return
		}

		var regs []*s.GroupDto
		var cancels []*s.GroupDto
		for _, e := range b.Events {
			switch e.Type {
			case "message":
				if e.Message.Text == RegisterMyGroupMsg && len(e.ReplyToken) > 0 {
					g := s.NewGroupDto(g.NewGroup(e.Source.GroupId), e.ReplyToken)
					regs = append(regs, g)
				}
			case "leave":
				g := s.NewGroupDto(g.NewGroup(e.Source.GroupId), "")
				cancels = append(cancels, g)
			default:
			}
		}

		var registered []*s.GroupDto
		for _, g := range regs {
			switch err := sCtx.RegistrationService.Execute(g); err {
			case nil:
				registered = append(registered, g)
			case s.ErrorGroupAlreadyRegistered:
				g.WasRegistered = true
				registered = append(registered, g)
			default:
				log.Printf("Error in registering group: %v", err)
			}
		}
		for _, g := range registered {
			if err := sCtx.ReplyService.Execute(g); err != nil {
				log.Printf("Error in replying to completed registration for group %v via LINE: %v", g.Id, err)
			}
		}

		ctx.Set("cancels", cancels)
		for _, g := range cancels {
			if err := sCtx.UnlistService.Execute(g); err != nil {
				log.Printf("Error in unlisting group: %v", err)
			}
		}

		ctx.Next()
	}
}
