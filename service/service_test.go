package service

import (
	"reflect"
	"testing"
)

type message struct {
	headers map[string]interface{}
	body    []byte
}

type publisherMock struct {
	messages []message
}

func (pm *publisherMock) Publish(queueName, exchange string, body []byte, headers map[string]interface{}) error {
	pm.messages = append(pm.messages, message{
		body:    body,
		headers: headers,
	})

	return nil
}
func TestSendMessage(t *testing.T) {

	tt := map[string]struct {
		body           []byte
		headers        map[string]interface{}
		messagesToSend int
	}{
		"one":              {body: []byte("one"), headers: map[string]interface{}{"one": 1}, messagesToSend: 1},
		"two":              {body: []byte("two"), headers: map[string]interface{}{"two": 2}, messagesToSend: 2},
		"ten":              {body: []byte("ten"), headers: map[string]interface{}{"ten": 10}, messagesToSend: 10},
		"no headers":       {body: []byte("no headers"), messagesToSend: 1},
		"multiple headers": {body: []byte("lots of headers"), headers: map[string]interface{}{"one": 1, "two": 2}, messagesToSend: 1},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			body := []byte("testing testing, 1.2.3")
			headers := map[string]interface{}{"hello": "goodbye"}

			mock := &publisherMock{}

			serv := NewService(mock)

			err := serv.SendMessage(body, headers, tc.messagesToSend)

			if err != nil {
				t.Fatalf("wasn't expecting an error but got one: %v", err)
			}

			if len(mock.messages) != tc.messagesToSend {
				t.Fatalf("expecting there to be %v message but got: %v", tc.messagesToSend, len(mock.messages))
			}

			for i := 0; i < tc.messagesToSend; i++ {
				if !reflect.DeepEqual(mock.messages[i].body, body) {
					t.Fatalf("expecting body '%s' but got '%s'", body, mock.messages[0].body)
				}

				if !reflect.DeepEqual(mock.messages[i].headers, headers) {
					t.Fatalf("expected headers to be '%+v', but got '%+v'", headers, mock.messages[0].headers)
				}
			}
		})
	}
}

func TestSendMessageWithoutBody(t *testing.T) {
	mock := &publisherMock{}

	serv := NewService(mock)

	err := serv.SendMessage(nil, nil, 1)

	if err == nil {
		t.Fatalf("expecting error for not providing body, but didn't get one")
	}
}
