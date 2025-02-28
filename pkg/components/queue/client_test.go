/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package queue

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

const (
	// defaultTestDequeueInterval is the default duration for the dequeue interval for inmemory test queue
	defaultTestDequeueInterval = time.Duration(5) * time.Millisecond
)

func TestStartDequeuer(t *testing.T) {
	mctrl := gomock.NewController(t)
	defer mctrl.Finish()

	mockCli := NewMockClient(mctrl)

	lastDequeueCh := make(chan struct{})

	firstCall := mockCli.EXPECT().Dequeue(gomock.Any(), gomock.Any()).Return(&Message{
		Metadata: Metadata{
			ID:            "testID",
			DequeueCount:  1,
			EnqueueAt:     time.Now(),
			ExpireAt:      time.Now().Add(10 * time.Hour),
			NextVisibleAt: time.Now().Add(5 * time.Minute),
		},
		ContentType: JSONContentType,
		Data:        []byte("{}"),
	}, nil)

	secondCall := mockCli.EXPECT().Dequeue(gomock.Any(), gomock.Any()).Return(nil, ErrInvalidMessage).After(firstCall.Call)
	mockCli.EXPECT().Dequeue(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, cfg QueueClientConfig) (*Message, error) {
			close(lastDequeueCh)
			return nil, ErrInvalidMessage
		}).AnyTimes().After(secondCall)

	ctx, cancel := context.WithCancel(context.TODO())
	msgCh, err := StartDequeuer(ctx, mockCli, WithDequeueInterval(defaultTestDequeueInterval))
	require.NoError(t, err)

	recvCnt := 0
	doneCh := make(chan struct{})

	// Consumer
	go func() {
		for msg := range msgCh {
			t.Logf("Dequeued Message ID: %s", msg.ID)
			recvCnt++
		}
		close(doneCh)
	}()

	<-lastDequeueCh
	cancel()
	<-doneCh

	require.Equal(t, 1, recvCnt)
}
