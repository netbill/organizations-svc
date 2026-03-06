# \InvitesAPI

All URIs are relative to *http://localhost:8003*

Method | HTTP request | Description
------------- | ------------- | -------------
[**OrganizationsSvcV1InvitesGet**](InvitesAPI.md#OrganizationsSvcV1InvitesGet) | **Get** /organizations-svc/v1/invites/ | Get invites
[**OrganizationsSvcV1InvitesInviteIdAcceptPatch**](InvitesAPI.md#OrganizationsSvcV1InvitesInviteIdAcceptPatch) | **Patch** /organizations-svc/v1/invites/{invite_id}/accept | Accept invite
[**OrganizationsSvcV1InvitesInviteIdCancelledPatch**](InvitesAPI.md#OrganizationsSvcV1InvitesInviteIdCancelledPatch) | **Patch** /organizations-svc/v1/invites/{invite_id}/cancelled | Cancel invite
[**OrganizationsSvcV1InvitesInviteIdDeclinePatch**](InvitesAPI.md#OrganizationsSvcV1InvitesInviteIdDeclinePatch) | **Patch** /organizations-svc/v1/invites/{invite_id}/decline | Decline invite
[**OrganizationsSvcV1InvitesInviteIdGet**](InvitesAPI.md#OrganizationsSvcV1InvitesInviteIdGet) | **Get** /organizations-svc/v1/invites/{invite_id} | Get invite
[**OrganizationsSvcV1InvitesPost**](InvitesAPI.md#OrganizationsSvcV1InvitesPost) | **Post** /organizations-svc/v1/invites/ | Create invite



## OrganizationsSvcV1InvitesGet

> InvitesCollection OrganizationsSvcV1InvitesGet(ctx).AccountId(accountId).OrganizationId(organizationId).Include(include).PageLimit(pageLimit).PageOffset(pageOffset).Execute()

Get invites



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
	accountId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Filter invites by account ID. Must match the authenticated account.  (optional)
	organizationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Filter invites by organization ID. The initiator must be a member of the organization.  (optional)
	include := []string{"Include_example"} // []string | Optional list of related resources to include in the response. Supported values: `profile`, `organization`.  (optional)
	pageLimit := int32(56) // int32 | Max number of items to return (1-100). (optional)
	pageOffset := int32(56) // int32 | Number of items to skip. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InvitesAPI.OrganizationsSvcV1InvitesGet(context.Background()).AccountId(accountId).OrganizationId(organizationId).Include(include).PageLimit(pageLimit).PageOffset(pageOffset).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InvitesAPI.OrganizationsSvcV1InvitesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OrganizationsSvcV1InvitesGet`: InvitesCollection
	fmt.Fprintf(os.Stdout, "Response from `InvitesAPI.OrganizationsSvcV1InvitesGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1InvitesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **accountId** | **uuid.UUID** | Filter invites by account ID. Must match the authenticated account.  | 
 **organizationId** | **uuid.UUID** | Filter invites by organization ID. The initiator must be a member of the organization.  | 
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


## OrganizationsSvcV1InvitesInviteIdAcceptPatch

> Invite OrganizationsSvcV1InvitesInviteIdAcceptPatch(ctx, inviteId).Execute()

Accept invite



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
	inviteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Invite ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InvitesAPI.OrganizationsSvcV1InvitesInviteIdAcceptPatch(context.Background(), inviteId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InvitesAPI.OrganizationsSvcV1InvitesInviteIdAcceptPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OrganizationsSvcV1InvitesInviteIdAcceptPatch`: Invite
	fmt.Fprintf(os.Stdout, "Response from `InvitesAPI.OrganizationsSvcV1InvitesInviteIdAcceptPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**inviteId** | **uuid.UUID** | Invite ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1InvitesInviteIdAcceptPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Invite**](Invite.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OrganizationsSvcV1InvitesInviteIdCancelledPatch

> Invite OrganizationsSvcV1InvitesInviteIdCancelledPatch(ctx, inviteId).Execute()

Cancel invite



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
	inviteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Invite ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InvitesAPI.OrganizationsSvcV1InvitesInviteIdCancelledPatch(context.Background(), inviteId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InvitesAPI.OrganizationsSvcV1InvitesInviteIdCancelledPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OrganizationsSvcV1InvitesInviteIdCancelledPatch`: Invite
	fmt.Fprintf(os.Stdout, "Response from `InvitesAPI.OrganizationsSvcV1InvitesInviteIdCancelledPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**inviteId** | **uuid.UUID** | Invite ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1InvitesInviteIdCancelledPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Invite**](Invite.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OrganizationsSvcV1InvitesInviteIdDeclinePatch

> Invite OrganizationsSvcV1InvitesInviteIdDeclinePatch(ctx, inviteId).Execute()

Decline invite



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
	inviteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Invite ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InvitesAPI.OrganizationsSvcV1InvitesInviteIdDeclinePatch(context.Background(), inviteId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InvitesAPI.OrganizationsSvcV1InvitesInviteIdDeclinePatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OrganizationsSvcV1InvitesInviteIdDeclinePatch`: Invite
	fmt.Fprintf(os.Stdout, "Response from `InvitesAPI.OrganizationsSvcV1InvitesInviteIdDeclinePatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**inviteId** | **uuid.UUID** | Invite ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1InvitesInviteIdDeclinePatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Invite**](Invite.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OrganizationsSvcV1InvitesInviteIdGet

> Invite OrganizationsSvcV1InvitesInviteIdGet(ctx, inviteId).Include(include).Execute()

Get invite



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
	inviteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Invite ID
	include := []string{"Include_example"} // []string | Optional list of related resources to include in the response. Supported values: `profile`, `organizations`.  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InvitesAPI.OrganizationsSvcV1InvitesInviteIdGet(context.Background(), inviteId).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InvitesAPI.OrganizationsSvcV1InvitesInviteIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OrganizationsSvcV1InvitesInviteIdGet`: Invite
	fmt.Fprintf(os.Stdout, "Response from `InvitesAPI.OrganizationsSvcV1InvitesInviteIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**inviteId** | **uuid.UUID** | Invite ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1InvitesInviteIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **include** | **[]string** | Optional list of related resources to include in the response. Supported values: &#x60;profile&#x60;, &#x60;organizations&#x60;.  | 

### Return type

[**Invite**](Invite.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OrganizationsSvcV1InvitesPost

> Invite OrganizationsSvcV1InvitesPost(ctx).CreateInvite(createInvite).Execute()

Create invite



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
	createInvite := *openapiclient.NewCreateInvite(*openapiclient.NewCreateInviteData("Type_example", *openapiclient.NewCreateInviteDataAttributes("TODO", "TODO"))) // CreateInvite | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InvitesAPI.OrganizationsSvcV1InvitesPost(context.Background()).CreateInvite(createInvite).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InvitesAPI.OrganizationsSvcV1InvitesPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OrganizationsSvcV1InvitesPost`: Invite
	fmt.Fprintf(os.Stdout, "Response from `InvitesAPI.OrganizationsSvcV1InvitesPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1InvitesPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createInvite** | [**CreateInvite**](CreateInvite.md) |  | 

### Return type

[**Invite**](Invite.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

