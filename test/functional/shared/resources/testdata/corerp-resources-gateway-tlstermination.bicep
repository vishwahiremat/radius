import radius as radius

@description('Specifies the environment for resources.')
param environment string

resource app 'Applications.Core/applications@2023-10-01-preview' = {
  name: 'corerp-resources-gateway-tlstermination'
  properties: {
    environment: environment
  }
}



resource certificate 'Applications.Core/secretStores@2023-10-01-preview' = {
  name: 'tls-gtwy-cert'
  properties: {
    type: 'generic'
    resource: 'radius-testing/tls-gtwy-cert'
    data: {
      'pat': {
        value: 'test'
      }
      'username': {
        value: 'vish'
      }
    }
  }
}


