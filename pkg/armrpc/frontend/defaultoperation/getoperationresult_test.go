// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package defaultoperation

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	v1 "github.com/project-radius/radius/pkg/armrpc/api/v1"
	manager "github.com/project-radius/radius/pkg/armrpc/asyncoperation/statusmanager"
	ctrl "github.com/project-radius/radius/pkg/armrpc/frontend/controller"
	radiustesting "github.com/project-radius/radius/pkg/corerp/testing"
	"github.com/project-radius/radius/pkg/ucp/store"
	"github.com/stretchr/testify/require"
)

func TestGetOperationResultRun(t *testing.T) {
	mctrl := gomock.NewController(t)
	defer mctrl.Finish()

	mStorageClient := store.NewMockStorageClient(mctrl)
	ctx := context.Background()

	rawDataModel := radiustesting.ReadFixture("operationstatus_datamodel.json")
	osDataModel := &manager.Status{}
	_ = json.Unmarshal(rawDataModel, osDataModel)

	rawExpectedOutput := radiustesting.ReadFixture("operationstatus_output.json")
	expectedOutput := &v1.AsyncOperationStatus{}
	_ = json.Unmarshal(rawExpectedOutput, expectedOutput)

	t.Run("get non-existing resource", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := radiustesting.GetARMTestHTTPRequest(ctx, http.MethodGet, testHeaderfile, nil)
		ctx := radiustesting.ARMTestContextFromRequest(req)

		mStorageClient.
			EXPECT().
			Get(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, id string, _ ...store.GetOptions) (*store.Object, error) {
				return nil, &store.ErrNotFound{}
			})

		ctl, err := NewGetOperationResult(ctrl.Options{
			StorageClient: mStorageClient,
		})

		require.NoError(t, err)
		resp, err := ctl.Run(ctx, req)
		require.NoError(t, err)
		_ = resp.Apply(ctx, w, req)
		require.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})

	opResTestCases := []struct {
		desc              string
		provisioningState v1.ProvisioningState
		respCode          int
		headersCheck      bool
	}{
		{
			"not-in-terminal-state",
			v1.ProvisioningStateAccepted,
			http.StatusAccepted,
			true,
		},
		{
			"put-succeeded-state",
			v1.ProvisioningStateSucceeded,
			http.StatusNoContent,
			false,
		},
		{
			"delete-succeeded-state",
			v1.ProvisioningStateSucceeded,
			http.StatusNoContent,
			false,
		},
		{
			"put-failed-state",
			v1.ProvisioningStateFailed,
			http.StatusNoContent,
			false,
		},
		{
			"delete-failed-state",
			v1.ProvisioningStateFailed,
			http.StatusNoContent,
			false,
		},
	}

	for _, tt := range opResTestCases {
		t.Run(tt.desc, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := radiustesting.GetARMTestHTTPRequest(ctx, http.MethodGet, testHeaderfile, nil)
			ctx := radiustesting.ARMTestContextFromRequest(req)

			osDataModel.Status = tt.provisioningState

			mStorageClient.
				EXPECT().
				Get(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, id string, _ ...store.GetOptions) (*store.Object, error) {
					return &store.Object{
						Metadata: store.Metadata{ID: id},
						Data:     osDataModel,
					}, nil
				})

			ctl, err := NewGetOperationResult(ctrl.Options{
				StorageClient: mStorageClient,
			})

			require.NoError(t, err)
			resp, err := ctl.Run(ctx, req)
			require.NoError(t, err)
			_ = resp.Apply(ctx, w, req)
			require.Equal(t, tt.respCode, w.Result().StatusCode)

			if tt.headersCheck {
				require.NotNil(t, w.Header().Get("Location"))
				require.Equal(t, req.URL.String(), w.Header().Get("Location"))

				require.NotNil(t, w.Header().Get("Retry-After"))
				require.Equal(t, v1.DefaultRetryAfter, w.Header().Get("Retry-After"))
			}
		})
	}
}