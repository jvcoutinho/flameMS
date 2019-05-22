package message

import "errors"

type Operation int

const (
	Initialize     Operation = 0
	Publish        Operation = 1
	Subscribe      Operation = 2
	CheckExistence Operation = 3
)

type Request struct {
	Requestor string
	Operation Operation
	QueueName string
	Body      interface{}
}

type Response struct {
	Error string
}

func (msg *Request) GetRequestor() string {
	return msg.Requestor
}

func (msg *Request) GetOperation() Operation {
	return msg.Operation
}

func (msg *Request) GetQueueName() string {
	return msg.QueueName
}

func (msg *Request) GetBody() interface{} {
	return msg.Body
}

func NewRequest(requestor string, operation Operation, queueName string, body interface{}) *Request {
	return &Request{requestor, operation, queueName, body}
}

func NewResponse(err error) *Response {
	if err != nil {
		return &Response{err.Error()}
	}
	return &Response{""}
}

func (msg *Response) GetError() error {
	return errors.New(msg.Error)
}

func (msg *Response) HasError() bool {
	return msg.Error != ""
}
