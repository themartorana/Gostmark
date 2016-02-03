package postmark

type Header struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

func NewHeader(name, value string) Header {
	return Header{
		Name:  name,
		Value: value,
	}
}
