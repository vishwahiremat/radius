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

// package storetest contains SHARED testing logic that is common to our data-store implementations.
package storetest

import (
	"encoding/json"
	"testing"

	"github.com/radius-project/radius/pkg/components/database"
	"github.com/radius-project/radius/pkg/ucp/resources"
	"github.com/radius-project/radius/pkg/ucp/util/etag"
	"github.com/radius-project/radius/test/testcontext"
	"github.com/stretchr/testify/require"
)

const (
	ResourceType1       = "System.Resources/resourceType1"
	ResourceType2       = "System.Resources/resourceType2"
	NestedResourceType1 = "System.Resources/resourceType1/nestedType"

	ResourcePath1       = "System.Resources/resourceType1/resource1"
	ResourcePath2       = "System.Resources/resourceType2/resource2"
	ResourcePath3       = "System.Resources/resourceType2/Resource3"
	NestedResourcePath1 = "System.Resources/resourceType1/resource1/nestedType/nested1"
	NestedResourcePath2 = "System.Resources/resourceType1/resource1/nestedType/nested2"
	NestedResourcePath3 = "System.Resources/resourceType1/resource2/nestedType/nested3"
	NestedResourcePath4 = "System.Resources/resourceType1/resource2/nestedType/nested4"

	RadiusScope         = "/planes/radius/local/"
	PlaneScope          = "/planes"
	ResourceGroup1Scope = "/planes/radius/local/resourceGroups/group1"
	ResourceGroup2Scope = "/planes/radius/local/resourceGroups/group2"
	ARMResourceScope    = "/subscriptions/abc/resourceGroups/group3"
)

var ResourceGroup1ID = parseOrPanic(ResourceGroup1Scope)
var ResourceGroup2ID = parseOrPanic(ResourceGroup2Scope)
var Resource1ID = parseOrPanic(ResourceGroup1Scope + "/providers/" + ResourcePath1)
var Resource2ID = parseOrPanic(ResourceGroup2Scope + "/providers/" + ResourcePath2)
var Resource3ID = parseOrPanic(ResourceGroup2Scope + "/providers/" + ResourcePath3)
var NestedResource1ID = parseOrPanic(ResourceGroup1Scope + "/providers/" + NestedResourcePath1)
var NestedResource2ID = parseOrPanic(ResourceGroup1Scope + "/providers/" + NestedResourcePath2)
var NestedResource3ID = parseOrPanic(ResourceGroup1Scope + "/providers/" + NestedResourcePath3)
var NestedResource4ID = parseOrPanic(ResourceGroup1Scope + "/providers/" + NestedResourcePath4)
var ARMResourceID = parseOrPanic(ARMResourceScope + "/providers/" + ResourcePath1)
var RadiusPlaneID = parseOrPanic(RadiusScope)

var ResourceGroup1Data = map[string]any{
	"value": "1",
	"properties": map[string]any{
		"group": "1",
	},
}

var ResourceGroup2Data = map[string]any{
	"value": "2",
	"properties": map[string]any{
		"group": "2",
	},
}

var Data1 = map[string]any{
	"value": "1",
	"properties": map[string]any{
		"resource": "1",
	},
}
var Data2 = map[string]any{
	"value": "2",
	"properties": map[string]any{
		"resource": "2",
	},
}
var Data3 = map[string]any{
	"value": "3",
	"properties": map[string]any{
		"resource": "3",
	},
}
var NestedData1 = map[string]any{
	"value": "n1",
	"properties": map[string]any{
		"resource": "n1",
	},
}
var NestedData2 = map[string]any{
	"value": "n2",
	"properties": map[string]any{
		"resource": "n2",
	},
}
var NestedData3 = map[string]any{
	"value": "n3",
	"properties": map[string]any{
		"resource": "n3",
	},
}
var NestedData4 = map[string]any{
	"value": "n4",
	"properties": map[string]any{
		"resource": "n4",
	},
}

var RadiusPlaneData = map[string]any{
	"value:": "1",
	"properties": map[string]any{
		"plane": "1",
	},
}

// MarshalOrPanic takes in any type and returns a byte slice, panicking if an error occurs while marshalling.
func MarshalOrPanic(in any) []byte {
	b, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}

	return b
}

func parseOrPanic(id string) resources.ID {
	parsed, err := resources.Parse(id)
	if err != nil {
		panic(err)
	}

	return parsed
}

func createObject(id resources.ID, data any) database.Object {
	return database.Object{
		Metadata: database.Metadata{
			ID: id.String(),
		},
		Data: data,
	}
}

func compareObjects(t *testing.T, expected *database.Object, actual *database.Object) {
	t.Helper()

	// Compare everything except ETags
	expectedCopy := *expected
	expectedCopy.ETag = ""

	actualCopy := *actual
	actualCopy.ETag = ""

	require.Equal(t, expectedCopy, actualCopy)
}

// CompareObjectLists compares two slices of store.Objects, ignoring their ETags.
func CompareObjectLists(t *testing.T, expected []database.Object, actual []database.Object) {
	t.Helper()

	expectedCopy := []database.Object{}
	expectedCopy = append(expectedCopy, expected...)

	actualCopy := []database.Object{}
	actualCopy = append(actualCopy, actual...)

	// Compare everything except ETags
	for i := range expectedCopy {
		expectedCopy[i].ETag = ""
	}

	for i := range actualCopy {
		actualCopy[i].ETag = ""
	}

	require.ElementsMatch(t, expectedCopy, actualCopy)
}

// This function tests the database Client's Get, Save and Delete methods by creating, updating and deleting objects with
// different IDs and scopes, and checks the results of various query scenarios with different filters and scopes. It also
// checks that the expected objects are returned.
func RunTest(t *testing.T, client database.Client, clear func(t *testing.T)) {
	ctx, cancel := testcontext.NewWithCancel(t)
	t.Cleanup(cancel)

	t.Run("get_not_found", func(t *testing.T) {
		clear(t)

		obj, err := client.Get(ctx, Resource1ID.String())
		require.ErrorIs(t, err, &database.ErrNotFound{ID: Resource1ID.String()})
		require.Nil(t, obj)
	})

	t.Run("delete_not_found", func(t *testing.T) {
		clear(t)

		err := client.Delete(ctx, Resource1ID.String())
		require.ErrorIs(t, err, &database.ErrNotFound{ID: Resource1ID.String()})
	})

	t.Run("save_and_get_arm", func(t *testing.T) {
		clear(t)
		// Testing that we can work with both UCP and ARM IDs.

		obj1 := createObject(ARMResourceID, Data1)
		err := client.Save(ctx, &obj1)
		require.NoError(t, err)
		require.NotEmpty(t, obj1.ETag)

		obj1Get, err := client.Get(ctx, ARMResourceID.String())
		require.NoError(t, err)
		compareObjects(t, &obj1, obj1Get)
		require.Equal(t, obj1Get.ETag, obj1.ETag)
	})

	t.Run("save_and_get_ucp", func(t *testing.T) {
		clear(t)
		// Testing that we can work with both UCP and ARM IDs.

		obj1 := createObject(Resource1ID, Data1)
		err := client.Save(ctx, &obj1)
		require.NoError(t, err)

		obj1Get, err := client.Get(ctx, Resource1ID.String())
		require.NoError(t, err)
		compareObjects(t, &obj1, obj1Get)
	})

	t.Run("save_and_get_scope", func(t *testing.T) {
		clear(t)
		// Testing that we can work with a scope like any other resource

		obj1 := createObject(ResourceGroup1ID, ResourceGroup1Data)
		err := client.Save(ctx, &obj1)
		require.NoError(t, err)

		obj1Get, err := client.Get(ctx, ResourceGroup1ID.String())
		require.NoError(t, err)
		compareObjects(t, &obj1, obj1Get)
	})

	t.Run("save_can_update", func(t *testing.T) {
		clear(t)

		obj1 := createObject(Resource1ID, Data1)
		err := client.Save(ctx, &obj1)
		require.NoError(t, err)

		obj1.Data = Data2
		err = client.Save(ctx, &obj1)
		require.NoError(t, err)

		obj1Get, err := client.Get(ctx, Resource1ID.String())
		require.NoError(t, err)
		compareObjects(t, &obj1, obj1Get)
	})

	t.Run("save_can_update_matching_etag", func(t *testing.T) {
		clear(t)

		obj1 := createObject(Resource1ID, Data1)
		err := client.Save(ctx, &obj1)
		require.NoError(t, err)
		require.NotEmpty(t, obj1.ETag)

		obj1.Data = Data2
		err = client.Save(ctx, &obj1, database.WithETag(obj1.ETag))
		require.NoError(t, err)

		obj1Get, err := client.Get(ctx, Resource1ID.String())
		require.NoError(t, err)
		compareObjects(t, &obj1, obj1Get)
	})

	t.Run("save_cannot_update_not_matching_etag", func(t *testing.T) {
		clear(t)

		obj1 := createObject(Resource1ID, Data1)
		err := client.Save(ctx, &obj1)
		require.NoError(t, err)

		obj1.Data = Data2
		err = client.Save(ctx, &obj1, database.WithETag(etag.New(MarshalOrPanic(Data2))))
		require.ErrorIs(t, err, &database.ErrConcurrency{})

		obj1.Data = Data1
		obj1Get, err := client.Get(ctx, Resource1ID.String())
		require.NoError(t, err)
		compareObjects(t, &obj1, obj1Get)
	})

	t.Run("save_cannot_update_missing_resource_with_not_matching_etag", func(t *testing.T) {
		clear(t)

		obj1 := createObject(Resource1ID, Data1)

		err := client.Save(ctx, &obj1, database.WithETag(etag.New(MarshalOrPanic(Data1))))
		require.ErrorIs(t, err, &database.ErrConcurrency{})

		obj1Get, err := client.Get(ctx, Resource1ID.String())
		require.ErrorIs(t, err, &database.ErrNotFound{ID: Resource1ID.String()})
		require.Nil(t, obj1Get)
	})

	t.Run("save_and_get_scope_only", func(t *testing.T) {
		clear(t)

		obj1 := createObject(parseOrPanic(ResourceGroup1Scope), Data1)
		err := client.Save(ctx, &obj1)
		require.NoError(t, err)

		obj1Get, err := client.Get(ctx, ResourceGroup1Scope)
		require.NoError(t, err)
		compareObjects(t, &obj1, obj1Get)
	})

	t.Run("save_and_delete", func(t *testing.T) {
		clear(t)

		obj1 := createObject(Resource1ID, Data1)
		err := client.Save(ctx, &obj1)
		require.NoError(t, err)

		err = client.Delete(ctx, Resource1ID.String())
		require.NoError(t, err)

		obj1Get, err := client.Get(ctx, Resource1ID.String())
		require.ErrorIs(t, err, &database.ErrNotFound{ID: Resource1ID.String()})
		require.Nil(t, obj1Get)
	})

	t.Run("save_and_delete_can_delete_with_matching_etag", func(t *testing.T) {
		clear(t)

		obj1 := createObject(Resource1ID, Data1)
		err := client.Save(ctx, &obj1)
		require.NoError(t, err)

		err = client.Delete(ctx, Resource1ID.String(), database.WithETag(obj1.ETag))
		require.NoError(t, err)

		obj1Get, err := client.Get(ctx, Resource1ID.String())
		require.ErrorIs(t, err, &database.ErrNotFound{ID: Resource1ID.String()})
		require.Nil(t, obj1Get)
	})

	t.Run("save_and_delete_cannot_delete_with_non_matching_etag", func(t *testing.T) {
		clear(t)

		obj1 := createObject(Resource1ID, Data1)
		err := client.Save(ctx, &obj1)
		require.NoError(t, err)

		err = client.Delete(ctx, Resource1ID.String(), database.WithETag(etag.New(MarshalOrPanic(Data2))))
		require.ErrorIs(t, err, &database.ErrConcurrency{})

		obj1Get, err := client.Get(ctx, Resource1ID.String())
		require.NoError(t, err)
		require.NotNil(t, obj1Get)
	})

	t.Run("delete_cannot_delete_missing_resource_with_not_matching_etag", func(t *testing.T) {
		clear(t)

		err := client.Delete(ctx, Resource1ID.String(), database.WithETag(etag.New(MarshalOrPanic(Data1))))
		require.ErrorIs(t, err, &database.ErrConcurrency{})
	})

	t.Run("list_can_be_empty", func(t *testing.T) {
		clear(t)

		objs, err := client.Query(ctx, database.Query{RootScope: RadiusScope, ResourceType: "asdf"})
		require.NoError(t, err)
		require.Empty(t, objs)
	})

	t.Run("query", func(t *testing.T) {
		clear(t)

		group1 := createObject(ResourceGroup1ID, ResourceGroup1Data)
		err := client.Save(ctx, &group1)
		require.NoError(t, err)

		group2 := createObject(ResourceGroup2ID, ResourceGroup2Data)
		err = client.Save(ctx, &group2)
		require.NoError(t, err)

		obj1 := createObject(Resource1ID, Data1)
		err = client.Save(ctx, &obj1)
		require.NoError(t, err)

		nested1 := createObject(NestedResource1ID, NestedData1)
		err = client.Save(ctx, &nested1)
		require.NoError(t, err)

		nested2 := createObject(NestedResource2ID, NestedData2)
		err = client.Save(ctx, &nested2)
		require.NoError(t, err)

		nested3 := createObject(NestedResource3ID, NestedData3)
		err = client.Save(ctx, &nested3)
		require.NoError(t, err)

		nested4 := createObject(NestedResource4ID, NestedData4)
		err = client.Save(ctx, &nested4)
		require.NoError(t, err)

		obj2 := createObject(Resource2ID, Data2)
		err = client.Save(ctx, &obj2)
		require.NoError(t, err)

		plane1 := createObject(RadiusPlaneID, RadiusPlaneData)
		err = client.Save(ctx, &plane1)
		require.NoError(t, err)

		t.Run("query_resources_at_resource_group_scope", func(t *testing.T) {
			objs, err := client.Query(ctx, database.Query{RootScope: ResourceGroup1Scope, ResourceType: NestedResourceType1})
			require.NoError(t, err)
			expected := []database.Object{
				nested1,
				nested2,
				nested3,
				nested4,
			}
			CompareObjectLists(t, expected, objs.Items)
		})

		t.Run("query_resources_at_resource_group_scope_with_field_filter", func(t *testing.T) {
			filters := []database.QueryFilter{{Field: "value", Value: "n1"}}
			objs, err := client.Query(ctx, database.Query{RootScope: ResourceGroup1Scope, ResourceType: NestedResourceType1, Filters: filters})
			require.NoError(t, err)
			expected := []database.Object{
				nested1,
			}
			CompareObjectLists(t, expected, objs.Items)
		})

		t.Run("query_resources_at_resource_group_scope_with_prefix", func(t *testing.T) {
			objs, err := client.Query(ctx, database.Query{RootScope: ResourceGroup1Scope, ResourceType: NestedResourceType1, RoutingScopePrefix: ResourcePath1})
			require.NoError(t, err)
			expected := []database.Object{
				nested1,
				nested2,
			}
			CompareObjectLists(t, expected, objs.Items)
		})

		t.Run("query_scopes_at_resource_group_scope", func(t *testing.T) {
			objs, err := client.Query(ctx, database.Query{RootScope: ResourceGroup1Scope, IsScopeQuery: true, ResourceType: "resourceGroups"})
			require.NoError(t, err)
			expected := []database.Object{}
			CompareObjectLists(t, expected, objs.Items)
		})

		t.Run("query_resources_at_plane_scope", func(t *testing.T) {
			objs, err := client.Query(ctx, database.Query{RootScope: RadiusScope, ResourceType: ResourceType1})
			require.NoError(t, err)
			require.Empty(t, objs)
		})

		t.Run("query_resources_at_plane_scope_recursive", func(t *testing.T) {
			objs, err := client.Query(ctx, database.Query{RootScope: RadiusScope, ScopeRecursive: true, ResourceType: NestedResourceType1})
			require.NoError(t, err)
			expected := []database.Object{
				nested1,
				nested2,
				nested3,
				nested4,
			}
			CompareObjectLists(t, expected, objs.Items)
		})

		t.Run("query_resources_at_plane_scope_recursive_with_field_filter", func(t *testing.T) {
			filters := []database.QueryFilter{{Field: "value", Value: "1"}}
			objs, err := client.Query(ctx, database.Query{RootScope: RadiusScope, ScopeRecursive: true, ResourceType: ResourceType1, Filters: filters})
			require.NoError(t, err)
			expected := []database.Object{
				obj1,
			}
			CompareObjectLists(t, expected, objs.Items)
		})

		t.Run("query_resources_at_plane_scope_recursive_with_prefix", func(t *testing.T) {
			objs, err := client.Query(ctx, database.Query{RootScope: RadiusScope, ScopeRecursive: true, ResourceType: NestedResourceType1, RoutingScopePrefix: ResourcePath1})
			require.NoError(t, err)
			expected := []database.Object{
				nested1,
				nested2,
			}
			CompareObjectLists(t, expected, objs.Items)
		})

		t.Run("query_scopes_at_plane_scope", func(t *testing.T) {
			objs, err := client.Query(ctx, database.Query{RootScope: PlaneScope, IsScopeQuery: true, ResourceType: "radius"})
			require.NoError(t, err)
			expected := []database.Object{
				plane1,
			}
			CompareObjectLists(t, expected, objs.Items)
		})

		t.Run("query_scopes_at_plane_scope_recursive", func(t *testing.T) {
			objs, err := client.Query(ctx, database.Query{RootScope: RadiusScope, ScopeRecursive: true, IsScopeQuery: true, ResourceType: "resourcegroups"})
			require.NoError(t, err)
			expected := []database.Object{
				group1,
				group2,
			}
			CompareObjectLists(t, expected, objs.Items)
		})

		t.Run("query_scopes_with_filter_non_matching", func(t *testing.T) {
			objs, err := client.Query(ctx, database.Query{RootScope: RadiusScope, IsScopeQuery: true, ResourceType: "resourcegroups", Filters: []database.QueryFilter{
				{Field: "value", Value: "asdf"},
			}})
			require.NoError(t, err)
			expected := []database.Object{}
			CompareObjectLists(t, expected, objs.Items)
		})

		t.Run("query_scopes_at_plane_scope_recursive", func(t *testing.T) {
			objs, err := client.Query(ctx, database.Query{RootScope: RadiusScope, ScopeRecursive: true, IsScopeQuery: true, ResourceType: "resourcegroups"})
			require.NoError(t, err)
			expected := []database.Object{
				group1,
				group2,
			}
			CompareObjectLists(t, expected, objs.Items)
		})

		t.Run("query_scopes_at_plane_scope_recursive_with_field_filter", func(t *testing.T) {
			filters := []database.QueryFilter{{Field: "value", Value: "1"}}
			objs, err := client.Query(ctx, database.Query{RootScope: RadiusScope, ScopeRecursive: true, IsScopeQuery: true, Filters: filters, ResourceType: "resourcegroups"})
			require.NoError(t, err)
			expected := []database.Object{
				group1,
			}
			CompareObjectLists(t, expected, objs.Items)
		})
	})
}
