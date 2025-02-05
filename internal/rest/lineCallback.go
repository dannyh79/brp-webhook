package routes

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
