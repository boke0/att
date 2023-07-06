package messages

import "github.com/boke0/att/pkg/primitives"

type ArtMessage struct {
	InitializeMessage *ArtInitializeMessage
	TextMessage       *ArtTextMessage
}

type ArtInitializeMessage struct {
	InitiatorId  string                           `json:"initiator_id"`
	Suk          []byte                           `json:"suk"`
	SukSignature []byte                           `json:"suk_signature"`
	Keys         map[string]primitives.PrivateKey `json:"keys"`
}

type ArtTextMessage struct {
	SenderId              string                           `json:"sender_id"`
	EphemeralKey          []byte                           `json:"ephemeral_key"`
	EphemeralKeySignature []byte                           `json:"ephemeral_key_signature"`
	Keys                  map[string]primitives.PrivateKey `json:"keys"`
}
