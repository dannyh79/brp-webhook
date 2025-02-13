package routes_test

import (
	"net/http"

	s "github.com/dannyh79/brp-webhook/internal/services"
)

var tcs = []struct {
	name         string
	expectStatus int

	reqBody             map[string]interface{}
	noSignatureHead     bool
	hasInvalidSignature bool

	expectedWelcomes      int
	shouldWelcomeFail     bool
	expectedUnlistings    int
	shouldUnlistFail      bool
	expectedRegistrations int
	shouldRegisterFail    bool
	registerFailError     error
	expectedReplies       int
	shouldReplyFail       bool
}{
	//#region authentication for LINE requests
	{
		name:         "Returns 400 when the request is missing body",
		expectStatus: http.StatusBadRequest,
	},
	{
		name:            "Returns 401 when the request is missing signature header",
		expectStatus:    http.StatusUnauthorized,
		reqBody:         map[string]interface{}{},
		noSignatureHead: true,
	},
	{
		name:                "Returns 401 when the request is not authorized",
		expectStatus:        http.StatusUnauthorized,
		hasInvalidSignature: true,
		reqBody:             map[string]interface{}{},
	},
	//#endregion

	{
		name:         `Returns 200 when there is no events in request body`,
		expectStatus: http.StatusOK,
		reqBody: map[string]interface{}{
			"events": []map[string]interface{}{},
		},
		expectedWelcomes:      0,
		expectedUnlistings:    0,
		expectedRegistrations: 0,
		expectedReplies:       0,
	},
	{
		name:         `Returns 200 when receiving "我需要好好靈修" text message event from a group`,
		expectStatus: http.StatusOK,
		reqBody: map[string]interface{}{
			"events": []map[string]interface{}{
				{
					"type": "message",
					"message": map[string]interface{}{
						"text": "我需要好好靈修",
					},
					"source": map[string]interface{}{
						"groupId": "C1234",
					},
					"replyToken": "test-reply-token",
				},
			},
		},
		expectedUnlistings:    0,
		expectedRegistrations: 1,
		expectedReplies:       1,
	},
	{
		name:         `Returns 200 when receiving "我需要好好靈修" text message event, leave event, and join event from groups`,
		expectStatus: http.StatusOK,
		reqBody: map[string]interface{}{
			"events": []map[string]interface{}{
				{
					"type": "message",
					"message": map[string]interface{}{
						"text": "我需要好好靈修",
					},
					"source": map[string]interface{}{
						"groupId": "C1234",
					},
					"replyToken": "test-reply-token",
				},
				{
					"type": "leave",
					"source": map[string]interface{}{
						"groupId": "C2345",
					},
				},
				{
					"type": "join",
					"source": map[string]interface{}{
						"groupId": "C3456",
					},
					"replyToken": "another-test-reply-token",
				},
			},
		},
		expectedWelcomes:      1,
		expectedUnlistings:    1,
		expectedRegistrations: 1,
		expectedReplies:       1,
	},
	{
		name:         `Returns 200 when receiving "我需要好好靈修" text message event from a group - registration failed from having existing record`,
		expectStatus: http.StatusOK,
		reqBody: map[string]interface{}{
			"events": []map[string]interface{}{
				{
					"type": "message",
					"message": map[string]interface{}{
						"text": "我需要好好靈修",
					},
					"source": map[string]interface{}{
						"groupId": "C1234",
					},
					"replyToken": "test-reply-token",
				},
			},
		},
		expectedUnlistings:    0,
		expectedRegistrations: 1,
		shouldRegisterFail:    true,
		registerFailError:     s.ErrorGroupAlreadyRegistered,
		expectedReplies:       1,
	},
	{
		name:         `Returns 200 when receiving "我需要好好靈修" text message event from a group - registration failed for any other reason`,
		expectStatus: http.StatusOK,
		reqBody: map[string]interface{}{
			"events": []map[string]interface{}{
				{
					"type": "message",
					"message": map[string]interface{}{
						"text": "我需要好好靈修",
					},
					"source": map[string]interface{}{
						"groupId": "C1234",
					},
					"replyToken": "test-reply-token",
				},
			},
		},
		expectedUnlistings:    0,
		expectedRegistrations: 1,
		shouldRegisterFail:    true,
		expectedReplies:       0,
	},
	{
		name:         `Returns 200 when receiving "我需要好好靈修" text message event from a group - reply failed`,
		expectStatus: http.StatusOK,
		reqBody: map[string]interface{}{
			"events": []map[string]interface{}{
				{
					"type": "message",
					"message": map[string]interface{}{
						"text": "我需要好好靈修",
					},
					"source": map[string]interface{}{
						"groupId": "C1234",
					},
					"replyToken": "test-reply-token",
				},
			},
		},
		expectedUnlistings:    0,
		expectedRegistrations: 1,
		shouldRegisterFail:    false,
		expectedReplies:       1,
		shouldReplyFail:       true,
	},
	{
		name:         `Returns 200 when receiving leave event`,
		expectStatus: http.StatusOK,
		reqBody: map[string]interface{}{
			"events": []map[string]interface{}{
				{
					"type": "leave",
					"source": map[string]interface{}{
						"groupId": "C1234",
					},
				},
			},
		},
		expectedUnlistings:    1,
		shouldUnlistFail:      false,
		expectedRegistrations: 0,
		shouldRegisterFail:    false,
		expectedReplies:       0,
		shouldReplyFail:       false,
	},
	{
		name:         `Returns 200 when receiving leave event - unlist failed`,
		expectStatus: http.StatusOK,
		reqBody: map[string]interface{}{
			"events": []map[string]interface{}{
				{
					"type": "leave",
					"source": map[string]interface{}{
						"groupId": "C1234",
					},
				},
			},
		},
		expectedUnlistings:    1,
		shouldUnlistFail:      true,
		expectedRegistrations: 0,
		shouldRegisterFail:    false,
		expectedReplies:       0,
		shouldReplyFail:       false,
	},
	{
		name:         `Returns 200 when receiving join event`,
		expectStatus: http.StatusOK,
		reqBody: map[string]interface{}{
			"events": []map[string]interface{}{
				{
					"type": "join",
					"source": map[string]interface{}{
						"groupId": "C1234",
					},
					"replyToken": "test-reply-token",
				},
			},
		},
		expectedWelcomes:      1,
		shouldWelcomeFail:     false,
		expectedUnlistings:    0,
		shouldUnlistFail:      false,
		expectedRegistrations: 0,
		shouldRegisterFail:    false,
		expectedReplies:       0,
		shouldReplyFail:       false,
	},
	{
		name:         `Returns 200 when receiving join event - welcome failed`,
		expectStatus: http.StatusOK,
		reqBody: map[string]interface{}{
			"events": []map[string]interface{}{
				{
					"type": "join",
					"source": map[string]interface{}{
						"groupId": "C1234",
					},
					"replyToken": "test-reply-token",
				},
			},
		},
		expectedWelcomes:      1,
		shouldWelcomeFail:     true,
		expectedUnlistings:    0,
		shouldUnlistFail:      true,
		expectedRegistrations: 0,
		shouldRegisterFail:    false,
		expectedReplies:       0,
		shouldReplyFail:       false,
	},
}
