package services

import (
	g "github.com/dannyh79/brp-webhook/internal/groups"
)

const lineReplyApiEndpoint = "https://api.line.me/v2/bot/message/reply"

const msgWelcome = g.MsgWelcome
const msgOk = g.MsgRegistrationOk
const msgAlreadyRegistered = g.MsgAlreadyRegistered
