package marshaller

import (
	"encoding/json"
	"fmt"

	"../message"
)

type MessageType int

const (
	Request  MessageType = 0
	Response MessageType = 1
)

type Marshaller interface {
	MarshalRequest(request message.Request) []byte
	MarshalResponse(response message.Response) []byte
	UnmarshalRequest(data []byte) message.Request
	UnmarshalResponse(data []byte) message.Response
}

type JSONMarshaller struct {
}

func (JSONMarshaller) MarshalRequest(request message.Request) []byte {
	marshalledData, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err.Error())
	}
	return marshalledData
}

func (JSONMarshaller) MarshalResponse(response message.Response) []byte {
	marshalledData, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err.Error())
	}
	return marshalledData
}

func (JSONMarshaller) UnmarshalRequest(data []byte) message.Request {

	unmarshalledData := message.Request{}
	err := json.Unmarshal(data, &unmarshalledData)
	if err != nil {
		fmt.Println(err.Error())
	}
	return unmarshalledData
}

func (JSONMarshaller) UnmarshalResponse(data []byte) message.Response {

	unmarshalledData := message.Response{}
	err := json.Unmarshal(data, &unmarshalledData)
	if err != nil {
		fmt.Println(err.Error())
	}
	return unmarshalledData
}
