package watcher

import (
	"encoding/base64"

	gmail "google.golang.org/api/gmail/v1"
	"gopkg.in/telegram-bot-api.v4"
)

func getMessageHeaders(msg *gmail.Message, headers ...string) map[string]string {
	headerMap := make(map[string]string)

	for _, header := range headers {
		for _, msgHeader := range msg.Payload.Headers {
			if header == msgHeader.Name {
				headerMap[header] = msgHeader.Value
				break
			}
		}
	}
	return headerMap
}

func getMessageBodyText(msg *gmail.Message) string {
	var bodyBytes []byte
	bodyPart := msg.Payload.Parts[0]
	switch bodyPart.MimeType {
	case "text/plain":
		bodyBytes, _ = base64.URLEncoding.DecodeString(bodyPart.Body.Data)
	case "multipart/alternative":
		bodyBytes, _ = base64.URLEncoding.DecodeString(bodyPart.Parts[0].Body.Data)
	}
	return string(bodyBytes[:len(bodyBytes)])
}

func getMessageAttachments(msg *gmail.Message, srv *gmail.Service, userID, msgID string) (attachments []tgbotapi.FileBytes) {
	parts := msg.Payload.Parts
	if parts[0].MimeType == "multipart/alternative" {
		for _, part := range parts[1:] {
			id := part.Body.AttachmentId
			attachment, _ := srv.Users.Messages.Attachments.Get(userID, msgID, id).Do()
			base64Bytes := []byte(attachment.Data)
			attBytes := make([]byte, base64.URLEncoding.DecodedLen(len(base64Bytes)))
			base64.URLEncoding.Decode(attBytes, base64Bytes)
			attachments = append(attachments, tgbotapi.FileBytes{Name: part.Filename, Bytes: attBytes})
		}
	}
	return attachments
}

func isInbox(msg *gmail.Message) bool {
	for _, label := range msg.LabelIds {
		if label == "INBOX" {
			return true
		}
	}
	return false
}
