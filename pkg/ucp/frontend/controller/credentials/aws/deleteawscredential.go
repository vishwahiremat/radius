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
package aws

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	v1 "github.com/radius-project/radius/pkg/armrpc/api/v1"
	armrpc_controller "github.com/radius-project/radius/pkg/armrpc/frontend/controller"
	armrpcrest "github.com/radius-project/radius/pkg/armrpc/rest"
	"github.com/radius-project/radius/pkg/components/database"
	"github.com/radius-project/radius/pkg/components/secret"
	"github.com/radius-project/radius/pkg/ucp/datamodel"
	"github.com/radius-project/radius/pkg/ucp/datamodel/converter"
	"github.com/radius-project/radius/pkg/ucp/frontend/controller/credentials"
	"github.com/radius-project/radius/pkg/ucp/ucplog"
)

var _ armrpc_controller.Controller = (*DeleteAWSCredential)(nil)

// DeleteAWSCredential is the controller implementation to delete a UCP AWS credential.
type DeleteAWSCredential struct {
	armrpc_controller.Operation[*datamodel.AWSCredential, datamodel.AWSCredential]
	secretClient secret.Client
}

// NewDeleteAWSCredential creates a new DeleteAWSCredential controller which is used to delete AWS credentials from the
// secret store, and returns it along with any errors that may have occurred.
func NewDeleteAWSCredential(opts armrpc_controller.Options, secretClient secret.Client) (armrpc_controller.Controller, error) {
	return &DeleteAWSCredential{
		Operation: armrpc_controller.NewOperation(opts,
			armrpc_controller.ResourceOptions[datamodel.AWSCredential]{
				RequestConverter:  converter.AWSCredentialDataModelFromVersioned,
				ResponseConverter: converter.AWSCredentialDataModelToVersioned,
			}),
		secretClient: secretClient,
	}, nil
}

// Run() checks if the AWS Credential exists, deletes the associated secret, and then deletes the AWS Credential from storage.
// If the AWS Credential does not exist, it returns a No Content response. If an error occurs, it returns an error.
func (c *DeleteAWSCredential) Run(ctx context.Context, w http.ResponseWriter, req *http.Request) (armrpcrest.Response, error) {
	logger := ucplog.FromContextOrDiscard(ctx)
	serviceCtx := v1.ARMRequestContextFromContext(ctx)

	old, etag, err := c.GetResource(ctx, serviceCtx.ResourceID)
	if err != nil {
		return nil, err
	}

	if old == nil {
		return armrpcrest.NewNoContentResponse(), nil
	}

	secretName := credentials.GetSecretName(serviceCtx.ResourceID)

	// Delete the credential secret.
	err = c.secretClient.Delete(ctx, secretName)
	if errors.Is(err, &secret.ErrNotFound{}) {
		return armrpcrest.NewNoContentResponse(), nil
	} else if err != nil {
		return nil, err
	}

	if r, err := c.PrepareResource(ctx, req, nil, old, etag); r != nil || err != nil {
		return r, err
	}

	if err := c.DatabaseClient().Delete(ctx, serviceCtx.ResourceID.String()); err != nil {
		if errors.Is(&database.ErrNotFound{ID: serviceCtx.ResourceID.String()}, err) {
			return armrpcrest.NewNoContentResponse(), nil
		}
		return nil, err
	}

	logger.Info(fmt.Sprintf("Deleted AWS Credential %s successfully", serviceCtx.ResourceID))
	return armrpcrest.NewOKResponse(nil), nil
}
