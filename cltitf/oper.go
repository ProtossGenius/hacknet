package cltitf

// CltOperItf client operator.
type CltOperItf interface {
	SendForward(email string, data []byte)
	OnForward(data []byte)
}
