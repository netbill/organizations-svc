# \OrganizationsAPI

All URIs are relative to *http://localhost:8003*

Method | HTTP request | Description
------------- | ------------- | -------------
[**OrganizationsSvcV1OrganizationsGet**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsGet) | **Get** /organizations-svc/v1/organizations/ | Get organizations
[**OrganizationsSvcV1OrganizationsOrganizationIdActivatePost**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsOrganizationIdActivatePost) | **Post** /organizations-svc/v1/organizations/{organization_id}/activate | Activate organization
[**OrganizationsSvcV1OrganizationsOrganizationIdDeactivatePost**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsOrganizationIdDeactivatePost) | **Post** /organizations-svc/v1/organizations/{organization_id}/deactivate | Deactivate organization
[**OrganizationsSvcV1OrganizationsOrganizationIdDelete**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsOrganizationIdDelete) | **Delete** /organizations-svc/v1/organizations/{organization_id} | Delete organization
[**OrganizationsSvcV1OrganizationsOrganizationIdGet**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsOrganizationIdGet) | **Get** /organizations-svc/v1/organizations/{organization_id} | Get organization
[**OrganizationsSvcV1OrganizationsOrganizationIdMediaDelete**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsOrganizationIdMediaDelete) | **Delete** /organizations-svc/v1/organizations/{organization_id}/media | Delete organization upload media
[**OrganizationsSvcV1OrganizationsOrganizationIdMediaPost**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsOrganizationIdMediaPost) | **Post** /organizations-svc/v1/organizations/{organization_id}/media | Create organization upload media link
[**OrganizationsSvcV1OrganizationsOrganizationIdPatch**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsOrganizationIdPatch) | **Patch** /organizations-svc/v1/organizations/{organization_id} | Update organization
[**OrganizationsSvcV1OrganizationsOrganizationIdSuspendPost**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsOrganizationIdSuspendPost) | **Post** /organizations-svc/v1/organizations/{organization_id}/suspend | Suspend organization
[**OrganizationsSvcV1OrganizationsOrganizationIdUnsuspendPost**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsOrganizationIdUnsuspendPost) | **Post** /organizations-svc/v1/organizations/{organization_id}/unsuspend | Unsuspend organization
[**OrganizationsSvcV1OrganizationsPost**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsPost) | **Post** /organizations-svc/v1/organizations/ | Create organization



## OrganizationsSvcV1OrganizationsGet

> OrganizationsCollection OrganizationsSvcV1OrganizationsGet(ctx).Text(text).Status(status).PageLimit(pageLimit).PageOffset(pageOffset).Execute()

Get organizations



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	text := "text_example" // string | Text filter for organizations. (optional)
	status := "status_example" // string | Filter by organization status. (optional)
	pageLimit := int32(56) // int32 | Max number of items to return (1-100). (optional)
	pageOffset := int32(56) // int32 | Number of items to skip. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OrganizationsAPI.OrganizationsSvcV1OrganizationsGet(context.Background()).Text(text).Status(status).PageLimit(pageLimit).PageOffset(pageOffset).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.OrganizationsSvcV1OrganizationsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OrganizationsSvcV1OrganizationsGet`: OrganizationsCollection
	fmt.Fprintf(os.Stdout, "Response from `OrganizationsAPI.OrganizationsSvcV1OrganizationsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1OrganizationsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **text** | **string** | Text filter for organizations. | 
 **status** | **string** | Filter by organization status. | 
 **pageLimit** | **int32** | Max number of items to return (1-100). | 
 **pageOffset** | **int32** | Number of items to skip. | 

### Return type

[**OrganizationsCollection**](OrganizationsCollection.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OrganizationsSvcV1OrganizationsOrganizationIdActivatePost

> Organization OrganizationsSvcV1OrganizationsOrganizationIdActivatePost(ctx, organizationId).Execute()

Activate organization



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	organizationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Organization ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdActivatePost(context.Background(), organizationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdActivatePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OrganizationsSvcV1OrganizationsOrganizationIdActivatePost`: Organization
	fmt.Fprintf(os.Stdout, "Response from `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdActivatePost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **uuid.UUID** | Organization ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1OrganizationsOrganizationIdActivatePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Organization**](Organization.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OrganizationsSvcV1OrganizationsOrganizationIdDeactivatePost

> Organization OrganizationsSvcV1OrganizationsOrganizationIdDeactivatePost(ctx, organizationId).Execute()

Deactivate organization



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	organizationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Organization ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdDeactivatePost(context.Background(), organizationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdDeactivatePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OrganizationsSvcV1OrganizationsOrganizationIdDeactivatePost`: Organization
	fmt.Fprintf(os.Stdout, "Response from `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdDeactivatePost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **uuid.UUID** | Organization ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1OrganizationsOrganizationIdDeactivatePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Organization**](Organization.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OrganizationsSvcV1OrganizationsOrganizationIdDelete

> OrganizationsSvcV1OrganizationsOrganizationIdDelete(ctx, organizationId).Execute()

Delete organization



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	organizationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Organization ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdDelete(context.Background(), organizationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **uuid.UUID** | Organization ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1OrganizationsOrganizationIdDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OrganizationsSvcV1OrganizationsOrganizationIdGet

> Organization OrganizationsSvcV1OrganizationsOrganizationIdGet(ctx, organizationId).Execute()

Get organization



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	organizationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Organization ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdGet(context.Background(), organizationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OrganizationsSvcV1OrganizationsOrganizationIdGet`: Organization
	fmt.Fprintf(os.Stdout, "Response from `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **uuid.UUID** | Organization ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1OrganizationsOrganizationIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Organization**](Organization.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OrganizationsSvcV1OrganizationsOrganizationIdMediaDelete

> OrganizationsSvcV1OrganizationsOrganizationIdMediaDelete(ctx, organizationId).DeleteUploadOrgMedia(deleteUploadOrgMedia).Execute()

Delete organization upload media



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	organizationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Organization ID
	deleteUploadOrgMedia := *openapiclient.NewDeleteUploadOrgMedia(*openapiclient.NewDeleteUploadOrgMediaData("TODO", "Type_example", *openapiclient.NewDeleteUploadOrgMediaDataAttributes())) // DeleteUploadOrgMedia | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdMediaDelete(context.Background(), organizationId).DeleteUploadOrgMedia(deleteUploadOrgMedia).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdMediaDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **uuid.UUID** | Organization ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1OrganizationsOrganizationIdMediaDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **deleteUploadOrgMedia** | [**DeleteUploadOrgMedia**](DeleteUploadOrgMedia.md) |  | 

### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OrganizationsSvcV1OrganizationsOrganizationIdMediaPost

> UploadOrgMediaLinks OrganizationsSvcV1OrganizationsOrganizationIdMediaPost(ctx, organizationId).Execute()

Create organization upload media link



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	organizationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Organization ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdMediaPost(context.Background(), organizationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdMediaPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OrganizationsSvcV1OrganizationsOrganizationIdMediaPost`: UploadOrgMediaLinks
	fmt.Fprintf(os.Stdout, "Response from `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdMediaPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **uuid.UUID** | Organization ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1OrganizationsOrganizationIdMediaPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**UploadOrgMediaLinks**](UploadOrgMediaLinks.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OrganizationsSvcV1OrganizationsOrganizationIdPatch

> Organization OrganizationsSvcV1OrganizationsOrganizationIdPatch(ctx, organizationId).UpdateOrganization(updateOrganization).Execute()

Update organization



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	organizationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Organization ID
	updateOrganization := *openapiclient.NewUpdateOrganization(*openapiclient.NewUpdateOrganizationData("TODO", "Type_example", *openapiclient.NewUpdateOrganizationDataAttributes())) // UpdateOrganization | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdPatch(context.Background(), organizationId).UpdateOrganization(updateOrganization).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OrganizationsSvcV1OrganizationsOrganizationIdPatch`: Organization
	fmt.Fprintf(os.Stdout, "Response from `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **uuid.UUID** | Organization ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1OrganizationsOrganizationIdPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateOrganization** | [**UpdateOrganization**](UpdateOrganization.md) |  | 

### Return type

[**Organization**](Organization.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OrganizationsSvcV1OrganizationsOrganizationIdSuspendPost

> OrganizationsSvcV1OrganizationsOrganizationIdSuspendPost(ctx, organizationId).Execute()

Suspend organization



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	organizationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Organization ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdSuspendPost(context.Background(), organizationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdSuspendPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **uuid.UUID** | Organization ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1OrganizationsOrganizationIdSuspendPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OrganizationsSvcV1OrganizationsOrganizationIdUnsuspendPost

> OrganizationsSvcV1OrganizationsOrganizationIdUnsuspendPost(ctx, organizationId).Execute()

Unsuspend organization



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	organizationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Organization ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdUnsuspendPost(context.Background(), organizationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdUnsuspendPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **uuid.UUID** | Organization ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1OrganizationsOrganizationIdUnsuspendPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OrganizationsSvcV1OrganizationsPost

> Organization OrganizationsSvcV1OrganizationsPost(ctx).CreateOrganization(createOrganization).Execute()

Create organization



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	createOrganization := *openapiclient.NewCreateOrganization(*openapiclient.NewCreateOrganizationData("Type_example", *openapiclient.NewCreateOrganizationDataAttributes("Name_example"))) // CreateOrganization | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OrganizationsAPI.OrganizationsSvcV1OrganizationsPost(context.Background()).CreateOrganization(createOrganization).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.OrganizationsSvcV1OrganizationsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OrganizationsSvcV1OrganizationsPost`: Organization
	fmt.Fprintf(os.Stdout, "Response from `OrganizationsAPI.OrganizationsSvcV1OrganizationsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1OrganizationsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createOrganization** | [**CreateOrganization**](CreateOrganization.md) |  | 

### Return type

[**Organization**](Organization.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

