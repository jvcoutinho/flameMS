package message

import "errors"

type Operation int

const (
	Register       Operation = 0
	Initialize     Operation = 1
	Publish        Operation = 2
	Subscribe      Operation = 3
	CheckExistence Operation = 4
	Topic          Operation = 5
	Stream         Operation = 6
)

type Request struct {
	Requestor string
	Operation Operation
	QueueName string
	Body      interface{}
	Priority  int
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

func (msg *Request) GetPriority() int {
	return msg.Priority
}

func NewRequest(requestor string, operation Operation, queueName string, body interface{}, priority int) *Request {
	return &Request{requestor, operation, queueName, body, priority}
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
