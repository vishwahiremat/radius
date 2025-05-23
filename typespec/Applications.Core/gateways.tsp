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

import "@typespec/rest";
import "@typespec/versioning";
import "@typespec/openapi";
import "@azure-tools/typespec-autorest";
import "@azure-tools/typespec-azure-core";
import "@azure-tools/typespec-azure-resource-manager";

import "../radius/v1/ucprootscope.tsp";
import "../radius/v1/resources.tsp";
import "./common.tsp";
import "../radius/v1/trackedresource.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;
using TypeSpec.Versioning;
using Autorest;
using Azure.Core;
using Azure.ResourceManager;
using OpenAPI;

namespace Applications.Core;

model GatewayResource
  is TrackedResourceRequired<GatewayProperties, "gateways"> {
  @doc("Gateway name")
  @key("gatewayName")
  @path
  @segment("gateways")
  name: ResourceNameString;
}

@doc("Gateway properties")
model GatewayProperties {
  ...ApplicationScopedResource;

  @doc("Sets Gateway to not be exposed externally (no public IP address associated). Defaults to false (exposed to internet).")
  internal?: boolean;

  @doc("Declare hostname information for the Gateway. Leaving the hostname empty auto-assigns one: mygateway.myapp.PUBLICHOSTNAMEORIP.nip.io.")
  hostname?: GatewayHostname;

  @doc("Routes attached to this Gateway")
  @extension("x-ms-identifiers", [])
  routes: GatewayRoute[];

  @doc("TLS configuration for the Gateway.")
  tls?: GatewayTls;

  @doc("URL of the gateway resource. Readonly")
  @visibility("read")
  url?: string;
}

@doc("Tls Minimum versions for Gateway resource.")
enum TlsMinVersion {
  @doc("TLS Version 1.2")
  tls12: "1.2",

  @doc("TLS Version 1.3")
  tls13: "1.3",
}

@doc("TLS configuration definition for Gateway resource.")
model GatewayTls {
  @doc("If true, gateway lets the https traffic sslPassthrough to the backend servers for decryption.")
  sslPassthrough?: boolean;

  @doc("TLS minimum protocol version (defaults to 1.2).")
  minimumProtocolVersion?: TlsMinVersion = TlsMinVersion.tls12;

  @doc("The resource id for the secret containing the TLS certificate and key for the gateway.")
  certificateFrom?: string;
}

@doc("Declare hostname information for the Gateway. Leaving the hostname empty auto-assigns one: mygateway.myapp.PUBLICHOSTNAMEORIP.nip.io.")
model GatewayHostname {
  @doc("Specify a prefix for the hostname: myhostname.myapp.PUBLICHOSTNAMEORIP.nip.io. Mutually exclusive with 'fullyQualifiedHostname' and will be overridden if both are defined.")
  prefix?: string;

  @doc("Specify a fully-qualified domain name: myapp.mydomain.com. Mutually exclusive with 'prefix' and will take priority if both are defined.")
  fullyQualifiedHostname?: string;
}

@doc("Route attached to Gateway")
model GatewayRoute {
  @doc("The path to match the incoming request path on. Ex - /myservice.")
  path?: string;

  @doc("The URL or id of the service to route to. Ex - 'http://myservice'.")
  destination?: string;

  @doc("Optionally update the prefix when sending the request to the service. Ex - replacePrefix: '/' and path: '/myservice' will transform '/myservice/myroute' to '/myroute'")
  replacePrefix?: string;

  @doc("Enables websocket support for the route. Defaults to false.")
  enableWebsockets?: boolean;

  @doc("The timeout policy for the route.")
  timeoutPolicy?: GatewayRouteTimeoutPolicy;
}

@doc("Gateway route timeout policy")
model GatewayRouteTimeoutPolicy {
  @doc("The request timeout in duration for the route. Defaults to 15 seconds.")
  @pattern("^(\\d+(\\.\\d+)?(ns|us|µs|ms|s|m|h))+$")
  request?: string;

  @doc("The backend request timeout in duration for the route. Cannot be greater than the request timeout.")
  @pattern("^(\\d+(\\.\\d+)?(ns|us|µs|ms|s|m|h))+$")
  backendRequest?: string;
}

@armResourceOperations
interface Gateways {
  get is ArmResourceRead<GatewayResource, UCPBaseParameters<GatewayResource>>;

  create is ArmResourceCreateOrReplaceAsync<
    GatewayResource,
    UCPBaseParameters<GatewayResource>
  >;

  createOrUpdate is ArmResourcePatchAsync<
    GatewayResource,
    GatewayProperties,
    UCPBaseParameters<GatewayResource>
  >;

  delete is ArmResourceDeleteWithoutOkAsync<
    GatewayResource,
    UCPBaseParameters<GatewayResource>
  >;

  listByScope is ArmResourceListByParent<
    GatewayResource,
    UCPBaseParameters<GatewayResource>,
    "Scope",
    "Scope"
  >;
}
