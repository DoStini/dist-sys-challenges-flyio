package node

type GenericPayload struct {
	Type  string `json:"type"`
	MsgID int    `json:"msg_id"`
}

type ResponseMessage struct {
	Type      string `json:"type"`
	InReplyTo int    `json:"in_reply_to"`
}

type InitRequest struct {
	NodeID  string   `json:"node_id"`
	NodeIDs []string `json:"node_ids"`
}
