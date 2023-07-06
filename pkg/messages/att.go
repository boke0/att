package messages

type AttMessage struct {
	InitializeMessage *AttInitializeMessage
	TextMessage       *AttTextMessage
}

type AttInitializeMessage struct {
	InitiatorId            string                              `json:"initiator_id"`
	InitiatorRoomUserId    string                              `json:"initiator_room_user_id"`
	InitializeKey          string                              `json:"initialize_key"`
	InitializeKeySignature string                              `json:"initialize_key_signature"`
	Users                  map[string]AttUserInitializeMessage `json:"users"`
}

type AttUserInitializeMessage struct {
	RoomUserId             string `json:"room_user_id"`
	SignedPrekeyId string `json:"signed_prekey_id"`
}

type AttTextMessage struct {
	SenderId              string                  `json:"initiator_id"`
	EphemeralKey          string                  `json:"ephemeral_key"`
	EphemeralKeySignature string                  `json:"ephemeral_key_signature"`
	Keys                  map[string]AttPublicKey `json:"keys"`
	Payload               string
}

type AttPublicKey struct {
	SenderId           string `json:"sender_id"`
	PublicKey          string `json:"public_key"`
	PublicKeySignature string `json:"public_key_signature"`
}
