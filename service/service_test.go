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

func (pm *publisherMock) Publish(queueName, exchange string, msg []byte, headers map[string]interface{}) error {
	pm.messages = append(pm.messages, message{
		body:    msg,
		headers: headers,
	})

	return nil
}
func TestSendMessage(t *testing.T) {

	tt := map[string]struct {
		messagesToSend int
	}{
		"one": {messagesToSend: 1},
		"two": {messagesToSend: 2},
		"ten": {messagesToSend: 10},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			msg := []byte("testing testing, 1.2.3")
			headers := map[string]interface{}{"hello": "goodbye"}

			mock := &publisherMock{}

			serv := NewService(mock)

			err := serv.SendMessage(msg, headers, tc.messagesToSend)

			if err != nil {
				t.Fatalf("wasn't expecting an error but got one: %v", err)
			}

			if len(mock.messages) != tc.messagesToSend {
				t.Fatalf("expecting there to be %v message but got: %v", tc.messagesToSend, len(mock.messages))
			}

			for i := 0; i < tc.messagesToSend; i++ {
				if !reflect.DeepEqual(mock.messages[i].body, msg) {
					t.Fatalf("expecting body '%s' but got '%s'", msg, mock.messages[0].body)
				}

				if !reflect.DeepEqual(mock.messages[i].headers, headers) {
					t.Fatalf("expected headers to be '%+v', but got '%+v'", headers, mock.messages[0].headers)
				}
			}
		})
	}
}
