# \OrganizationsAPI

All URIs are relative to *http://localhost:8003*

Method | HTTP request | Description
------------- | ------------- | -------------
[**OrganizationsSvcV1OrganizationsGet**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsGet) | **Get** /organizations-svc/v1/organizations/ | Get organizations
[**OrganizationsSvcV1OrganizationsOrganizationIdActivatePost**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsOrganizationIdActivatePost) | **Post** /organizations-svc/v1/organizations/{organization_id}/activate | Activate organization
[**OrganizationsSvcV1OrganizationsOrganizationIdDeactivatePost**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsOrganizationIdDeactivatePost) | **Post** /organizations-svc/v1/organizations/{organization_id}/deactivate | Deactivate organization
[**OrganizationsSvcV1OrganizationsOrganizationIdDelete**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsOrganizationIdDelete) | **Delete** /organizations-svc/v1/organizations/{organization_id} | Delete organization
[**OrganizationsSvcV1OrganizationsOrganizationIdGet**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsOrganizationIdGet) | **Get** /organizations-svc/v1/organizations/{organization_id} | Get organization
[**OrganizationsSvcV1OrganizationsOrganizationIdInvitesGet**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsOrganizationIdInvitesGet) | **Get** /organizations-svc/v1/organizations/{organization_id}/invites | Get organization invites
[**OrganizationsSvcV1OrganizationsOrganizationIdLeaveDelete**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsOrganizationIdLeaveDelete) | **Delete** /organizations-svc/v1/organizations/{organization_id}/leave | Leave organization
[**OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadBannerDelete**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadBannerDelete) | **Delete** /organizations-svc/v1/organizations/{organization_id}/media/upload/banner | Delete organization upload banner
[**OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadIconDelete**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadIconDelete) | **Delete** /organizations-svc/v1/organizations/{organization_id}/media/upload/icon | Delete organization upload icon
[**OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadUrlPost**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadUrlPost) | **Post** /organizations-svc/v1/organizations/{organization_id}/media/upload/url | Create organization upload media link
[**OrganizationsSvcV1OrganizationsOrganizationIdPut**](OrganizationsAPI.md#OrganizationsSvcV1OrganizationsOrganizationIdPut) | **Put** /organizations-svc/v1/organizations/{organization_id} | Update organization
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


## OrganizationsSvcV1OrganizationsOrganizationIdInvitesGet

> InvitesCollection OrganizationsSvcV1OrganizationsOrganizationIdInvitesGet(ctx, organizationId).Include(include).PageLimit(pageLimit).PageOffset(pageOffset).Execute()

Get organization invites



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
	include := []string{"Include_example"} // []string | Optional list of related resources to include in the response. Supported values: `profile`, `organization`.  (optional)
	pageLimit := int32(56) // int32 | Max number of items to return (1-100). (optional)
	pageOffset := int32(56) // int32 | Number of items to skip. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdInvitesGet(context.Background(), organizationId).Include(include).PageLimit(pageLimit).PageOffset(pageOffset).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdInvitesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OrganizationsSvcV1OrganizationsOrganizationIdInvitesGet`: InvitesCollection
	fmt.Fprintf(os.Stdout, "Response from `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdInvitesGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **uuid.UUID** | Organization ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1OrganizationsOrganizationIdInvitesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **include** | **[]string** | Optional list of related resources to include in the response. Supported values: &#x60;profile&#x60;, &#x60;organization&#x60;.  | 
 **pageLimit** | **int32** | Max number of items to return (1-100). | 
 **pageOffset** | **int32** | Number of items to skip. | 

### Return type

[**InvitesCollection**](InvitesCollection.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OrganizationsSvcV1OrganizationsOrganizationIdLeaveDelete

> OrganizationsSvcV1OrganizationsOrganizationIdLeaveDelete(ctx, organizationId).Execute()

Leave organization



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
	r, err := apiClient.OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdLeaveDelete(context.Background(), organizationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdLeaveDelete``: %v\n", err)
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

Other parameters are passed through a pointer to a apiOrganizationsSvcV1OrganizationsOrganizationIdLeaveDeleteRequest struct via the builder pattern


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


## OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadBannerDelete

> OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadBannerDelete(ctx, organizationId).DeleteUploadOrgBanner(deleteUploadOrgBanner).Execute()

Delete organization upload banner



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
	deleteUploadOrgBanner := *openapiclient.NewDeleteUploadOrgBanner(*openapiclient.NewDeleteUploadOrgBannerData("TODO", "Type_example", *openapiclient.NewDeleteUploadOrgBannerDataAttributes("BannerKey_example"))) // DeleteUploadOrgBanner | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadBannerDelete(context.Background(), organizationId).DeleteUploadOrgBanner(deleteUploadOrgBanner).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadBannerDelete``: %v\n", err)
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

Other parameters are passed through a pointer to a apiOrganizationsSvcV1OrganizationsOrganizationIdMediaUploadBannerDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **deleteUploadOrgBanner** | [**DeleteUploadOrgBanner**](DeleteUploadOrgBanner.md) |  | 

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


## OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadIconDelete

> OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadIconDelete(ctx, organizationId).DeleteUploadOrgIcon(deleteUploadOrgIcon).Execute()

Delete organization upload icon



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
	deleteUploadOrgIcon := *openapiclient.NewDeleteUploadOrgIcon(*openapiclient.NewDeleteUploadOrgIconData("TODO", "Type_example", *openapiclient.NewDeleteUploadOrgIconDataAttributes("IconKey_example"))) // DeleteUploadOrgIcon | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadIconDelete(context.Background(), organizationId).DeleteUploadOrgIcon(deleteUploadOrgIcon).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadIconDelete``: %v\n", err)
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

Other parameters are passed through a pointer to a apiOrganizationsSvcV1OrganizationsOrganizationIdMediaUploadIconDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **deleteUploadOrgIcon** | [**DeleteUploadOrgIcon**](DeleteUploadOrgIcon.md) |  | 

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


## OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadUrlPost

> UploadOrgMediaLinks OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadUrlPost(ctx, organizationId).Execute()

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
	resp, r, err := apiClient.OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadUrlPost(context.Background(), organizationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadUrlPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadUrlPost`: UploadOrgMediaLinks
	fmt.Fprintf(os.Stdout, "Response from `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdMediaUploadUrlPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **uuid.UUID** | Organization ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1OrganizationsOrganizationIdMediaUploadUrlPostRequest struct via the builder pattern


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


## OrganizationsSvcV1OrganizationsOrganizationIdPut

> Organization OrganizationsSvcV1OrganizationsOrganizationIdPut(ctx, organizationId).UpdateOrganization(updateOrganization).Execute()

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
	updateOrganization := *openapiclient.NewUpdateOrganization(*openapiclient.NewUpdateOrganizationData("TODO", "Type_example", *openapiclient.NewUpdateOrganizationDataAttributes("Name_example"))) // UpdateOrganization | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdPut(context.Background(), organizationId).UpdateOrganization(updateOrganization).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OrganizationsSvcV1OrganizationsOrganizationIdPut`: Organization
	fmt.Fprintf(os.Stdout, "Response from `OrganizationsAPI.OrganizationsSvcV1OrganizationsOrganizationIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organizationId** | **uuid.UUID** | Organization ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1OrganizationsOrganizationIdPutRequest struct via the builder pattern


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

