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

package statusmanager

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	v1 "github.com/radius-project/radius/pkg/armrpc/api/v1"
	"github.com/radius-project/radius/pkg/armrpc/rpctest"
	"github.com/radius-project/radius/pkg/components/database"
	"github.com/radius-project/radius/pkg/components/queue"
	"github.com/radius-project/radius/pkg/ucp/resources"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type asyncOperationsManagerTest struct {
	manager        StatusManager
	databaseClient *database.MockClient
	queueClient    *queue.MockClient
}

const (
	operationTimeoutDuration      = time.Hour * 2
	opererationRetryAfterDuration = time.Second * 10
	azureEnvResourceID            = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/radius-test-rg/providers/Applications.Core/environments/env0"
	ucpEnvResourceID              = "/planes/radius/local/resourceGroups/radius-test-rg/providers/Applications.Core/environments/env0"
	saveErr                       = "save error"
	enqueueErr                    = "enqueue error"
	deleteErr                     = "delete error"
	getErr                        = "get error"
)

func setup(tb testing.TB) (asyncOperationsManagerTest, *gomock.Controller) {
	ctrl := gomock.NewController(tb)
	sc := database.NewMockClient(ctrl)
	enq := queue.NewMockClient(ctrl)
	aom := New(sc, enq, "test-location")
	return asyncOperationsManagerTest{manager: aom, databaseClient: sc, queueClient: enq}, ctrl
}

var reqCtx = &v1.ARMRequestContext{
	ResourceID:     resources.MustParse("/planes/radius/local/resourceGroups/radius-test-rg/providers/Applications.Core/container/container0"),
	OperationID:    uuid.Must(uuid.NewRandom()),
	HomeTenantID:   "home-tenant-id",
	ClientObjectID: "client-object-id",
	OperationType:  rpctest.MustParseOperationType("APPLICATIONS.CORE/ENVIRONMENTS|PUT"),
	Traceparent:    "trace",
	AcceptLanguage: "lang",
}

var opID = uuid.New()

var testAos = &Status{
	AsyncOperationStatus: v1.AsyncOperationStatus{
		ID:        opID.String(),
		Name:      opID.String(),
		Status:    v1.ProvisioningStateUpdating,
		StartTime: time.Now().UTC(),
	},
	LinkedResourceID: uuid.New().String(),
	Location:         "test-location",
	RetryAfter:       opererationRetryAfterDuration,
	HomeTenantID:     "test-home-tenant-id",
	ClientObjectID:   "test-client-object-id",
}

func TestOperationStatusResourceID(t *testing.T) {
	resourceIDTests := []struct {
		resourceID          string
		operationID         uuid.UUID
		operationResourceID string
	}{
		{
			resourceID:          azureEnvResourceID,
			operationID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			operationResourceID: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/applications.core/locations/global/operationstatuses/00000000-0000-0000-0000-000000000001",
		}, {
			resourceID:          ucpEnvResourceID,
			operationID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			operationResourceID: "/planes/radius/local/providers/applications.core/locations/global/operationstatuses/00000000-0000-0000-0000-000000000001",
		},
	}

	sm := &statusManager{location: v1.LocationGlobal}

	for _, tc := range resourceIDTests {
		t.Run(tc.resourceID, func(t *testing.T) {
			rid, err := resources.ParseResource(tc.resourceID)
			require.NoError(t, err)
			url := sm.operationStatusResourceID(rid, tc.operationID)
			require.Equal(t, tc.operationResourceID, url)
		})
	}
}

func TestCreateAsyncOperationStatus(t *testing.T) {
	createCases := []struct {
		Desc       string
		SaveErr    error
		EnqueueErr error
		DeleteErr  error
	}{
		{
			Desc:       "create_success",
			SaveErr:    nil,
			EnqueueErr: nil,
			DeleteErr:  nil,
		},
		{
			Desc:       "create_save-error",
			SaveErr:    errors.New(saveErr),
			EnqueueErr: nil,
			DeleteErr:  nil,
		},
		{
			Desc:       "create_enqueue-error",
			SaveErr:    nil,
			EnqueueErr: errors.New(enqueueErr),
			DeleteErr:  nil,
		},
		{
			Desc:       "create_delete-error",
			SaveErr:    nil,
			EnqueueErr: errors.New(enqueueErr),
			DeleteErr:  errors.New(deleteErr),
		},
	}

	for _, tt := range createCases {
		t.Run(fmt.Sprint(tt.Desc), func(t *testing.T) {
			aomTest, mctrl := setup(t)
			defer mctrl.Finish()

			aomTest.databaseClient.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.SaveErr)

			// We can't expect an async operation to be queued if it is not saved to the DB.
			if tt.SaveErr == nil {
				aomTest.queueClient.EXPECT().Enqueue(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.EnqueueErr)
			}

			// If there is an error when enqueuing the message, the async operation should be deleted.
			if tt.EnqueueErr != nil {
				aomTest.databaseClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.DeleteErr)
			}

			options := QueueOperationOptions{
				OperationTimeout: operationTimeoutDuration,
				RetryAfter:       opererationRetryAfterDuration,
			}
			err := aomTest.manager.QueueAsyncOperation(context.TODO(), reqCtx, options)

			if tt.SaveErr == nil && tt.EnqueueErr == nil && tt.DeleteErr == nil {
				require.NoError(t, err)
			}

			if tt.SaveErr != nil {
				require.Error(t, err, saveErr)
			}

			if tt.EnqueueErr != nil {
				require.Error(t, err, enqueueErr)
			}

			if tt.DeleteErr != nil {
				require.Error(t, err, deleteErr)
			}
		})
	}
}

func TestDeleteAsyncOperationStatus(t *testing.T) {
	deleteCases := []struct {
		Desc      string
		DeleteErr error
	}{
		{
			Desc:      "delete_success",
			DeleteErr: nil,
		},
		{
			Desc:      "delete_error",
			DeleteErr: errors.New(deleteErr),
		},
	}

	for _, tt := range deleteCases {
		t.Run(fmt.Sprint(tt.Desc), func(t *testing.T) {
			aomTest, mctrl := setup(t)
			defer mctrl.Finish()

			aomTest.databaseClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.DeleteErr)
			rid, err := resources.ParseResource(azureEnvResourceID)
			require.NoError(t, err)
			err = aomTest.manager.Delete(context.TODO(), rid, uuid.New())

			if tt.DeleteErr != nil {
				require.Error(t, err, deleteErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetAsyncOperationStatus(t *testing.T) {
	getCases := []struct {
		Desc   string
		GetErr error
		Obj    *database.Object
	}{
		{
			Desc:   "get_success",
			GetErr: nil,
			Obj: &database.Object{
				Metadata: database.Metadata{ID: opID.String(), ETag: "etag"},
				Data:     testAos,
			},
		},
		{
			Desc:   "create_enqueue-error",
			GetErr: errors.New(getErr),
			Obj:    nil,
		},
	}

	for _, tt := range getCases {
		t.Run(fmt.Sprint(tt.Desc), func(t *testing.T) {
			aomTest, mctrl := setup(t)
			defer mctrl.Finish()

			aomTest.databaseClient.
				EXPECT().
				Get(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(tt.Obj, tt.GetErr)

			rid, err := resources.ParseResource(azureEnvResourceID)
			require.NoError(t, err)
			aos, err := aomTest.manager.Get(context.TODO(), rid, uuid.New())

			if tt.GetErr == nil {
				require.NoError(t, err)
				expected := &Status{}
				_ = tt.Obj.As(&expected)
				require.Equal(t, expected, aos)
			}

			if tt.GetErr != nil {
				require.Error(t, err, getErr)
			}
		})
	}
}

func TestUpdateAsyncOperationStatus(t *testing.T) {
	updateCases := []struct {
		Desc    string
		GetErr  error
		Obj     *database.Object
		SaveErr error
	}{
		{
			Desc:   "update_success",
			GetErr: nil,
			Obj: &database.Object{
				Metadata: database.Metadata{ID: opID.String(), ETag: "etag"},
				Data:     testAos,
			},
			SaveErr: nil,
		},
	}

	for _, tt := range updateCases {
		t.Run(fmt.Sprint(tt.Desc), func(t *testing.T) {
			aomTest, mctrl := setup(t)
			defer mctrl.Finish()

			aomTest.databaseClient.
				EXPECT().
				Get(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(tt.Obj, tt.GetErr)

			if tt.GetErr == nil {
				aomTest.databaseClient.
					EXPECT().
					Save(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(tt.SaveErr)
			}

			testAos.Status = v1.ProvisioningStateSucceeded
			rid, err := resources.ParseResource(azureEnvResourceID)
			require.NoError(t, err)
			err = aomTest.manager.Update(context.TODO(), rid, opID, v1.ProvisioningStateAccepted, nil, nil)

			if tt.GetErr == nil && tt.SaveErr == nil {
				require.NoError(t, err)
			}

			if tt.GetErr != nil {
				require.Error(t, err, getErr)
			}

			if tt.SaveErr != nil {
				require.Error(t, err, saveErr)
			}
		})
	}
}
