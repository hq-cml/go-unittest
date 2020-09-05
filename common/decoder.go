package common

type Decoder interface {
	Unmarshal([]byte, interface{}) error
}

func NewDecoder() Decoder {
	var d Decoder
	return d
}

type Movie struct {
	Name string
	Type string
	Score int
}