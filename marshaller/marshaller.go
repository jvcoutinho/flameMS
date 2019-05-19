package marshaller

import (
	"encoding/json"
	"fmt"

	"../message"
)

type Marshaller interface {
	Marshal(msg message.Message) []byte
	Unmarshal(data []byte) message.Message
}

type JSONMarshaller struct {
}

func (marshaller *JSONMarshaller) Marshal(msg message.Message) []byte {
	marshalledData, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err.Error())
	}
	return marshalledData
}

func (marshaller *JSONMarshaller) Unmarshal(data []byte) message.Message {
	unmarshalledData := message.Message{}
	err := json.Unmarshal(data, &unmarshalledData)
	if err != nil {
		fmt.Println(err.Error())
	}
	return unmarshalledData
}
