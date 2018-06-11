package gostmark

type MessageStatus string

const (
	Queued    MessageStatus = "queued"
	Sent      MessageStatus = "sent"
	Processed MessageStatus = "processed"
)
