package internal

type BuildMessageState int8

const (
	BuildMessageStateReadMessageSignature BuildMessageState = iota
	BuildMessageStateReadMessageLength
	BuildMessageStateReadMessageData
)
