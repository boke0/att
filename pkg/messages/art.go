package messages

type ArtMessage struct {
	InitializeMessage *ArtInitializeMessage
	TextMessage       *ArtTextMessage
}

type ArtInitializeMessage struct {
	InitiatorId       string                  `json:"initiator_id"`
	SetupKey          string                  `json:"setup_key"`
	SetupKeySignature string                  `json:"setup_key_signature"`
	Keys              map[string]ArtPublicKey `json:"keys"`
}

type ArtTextMessage struct {
	SenderId              string                           `json:"sender_id"`
	EphemeralKey          string                           `json:"ephemeral_key"`
	EphemeralKeySignature string                           `json:"ephemeral_key_signature"`
	Keys                  map[string]ArtPublicKey `json:"keys"`
	Payload               string                           `json:"payload"`
}

type ArtPublicKey struct {
	SenderId           string `json:"sender_id"`
	PublicKey          string `json:"public_key"`
	PublicKeySignature string `json:"public_key_signature"`
}
