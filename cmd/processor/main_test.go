package main

import (
    "context"
    "database/sql"
    "testing"
    "fmt"

    "github.com/golang/mock/gomock"
    "cloud.google.com/go/pubsub"
    "github.com/stretchr/testify/assert"
)

// SubscriptionReceiver is an interface that mocks pubsub.Subscription's methods
type SubscriptionReceiver interface {
    Receive(ctx context.Context, f func(ctx context.Context, msg *pubsub.Message)) error
}

// MockDB is used to mock the database interactions for testing
type MockDB struct {
    // Mock DB connection
    ExecFunc func(query string, args ...interface{}) (sql.Result, error)
}

func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
    return m.ExecFunc(query, args...)
}

// MockSubscription is used to mock pubsub.Subscription for testing
type MockSubscription struct {
    ReceiveFunc func(ctx context.Context, handler func(ctx context.Context, msg *pubsub.Message)) error
}

func (m *MockSubscription) Receive(ctx context.Context, handler func(ctx context.Context, msg *pubsub.Message)) error {
    return m.ReceiveFunc(ctx, handler)
}

func TestMainLoop(t *testing.T) {
    t.Log("Starting test of main_loop")
    // Mock context
    ctx := context.Background()

    // Create mock controller
    mockCtrl := gomock.NewController(t)
    defer mockCtrl.Finish()

    // Mock database interaction
    mockDB := &MockDB{
        ExecFunc: func(query string, args ...interface{}) (sql.Result, error) {
            t.Logf("In MockDB insert %s, %s", query, args)
            // Simulate a successful database insert
            if query == "INSERT OR REPLACE INTO scans (ip, port, service, response, scan_time) VALUES (?, ?, ?, ?, ?)" &&
				args[3] != "abdul" {
                t.Log("Succesful Insert")
                return nil, nil
            }
            t.Log("In MockDB error condition")
            return nil, fmt.Errorf("database insert failed")
        },
    }

    // Mock Pub/Sub subscription
    mockSub := &MockSubscription{
        ReceiveFunc: func(ctx context.Context, handler func(ctx context.Context, msg *pubsub.Message)) error {
            // Create a mock Pub/Sub message
            mockMsg1 := &pubsub.Message{
                Data: []byte(`{"ip":"192.168.1.1","port":80,"service":"HTTP","data_version":1,"data":{"response_bytes_utf8":[97,98,99,100]}}`),
            }

            // Create second mock Pub/Sub message (error case)
            mockMsg2 := &pubsub.Message{
                Data: []byte(`{"ip":"192.168.1.1","port":80,"service":"HTTP","data_version":2,"data":{"response_str":"abdul"}}`),
            }

            t.Log("Sending first message")
            handler(ctx, mockMsg1) // Simulate receiving message 1
            t.Log("Sending second message")
            handler(ctx, mockMsg2) // Simulate receiving message 2

            return nil
        },
    }

    // Test that the first message was processed correctly and inserted into the database and the second message panics

    assert.Panics(t, func() { main_loop(ctx, mockSub, mockDB) }, "Function should panic")

}
