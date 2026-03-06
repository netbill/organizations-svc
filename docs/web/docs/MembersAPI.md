# \MembersAPI

All URIs are relative to *http://localhost:8003*

Method | HTTP request | Description
------------- | ------------- | -------------
[**OrganizationsSvcV1MembersGet**](MembersAPI.md#OrganizationsSvcV1MembersGet) | **Get** /organizations-svc/v1/members/ | Get members
[**OrganizationsSvcV1MembersMemberIdDelete**](MembersAPI.md#OrganizationsSvcV1MembersMemberIdDelete) | **Delete** /organizations-svc/v1/members/{member_id} | Delete member
[**OrganizationsSvcV1MembersMemberIdGet**](MembersAPI.md#OrganizationsSvcV1MembersMemberIdGet) | **Get** /organizations-svc/v1/members/{member_id} | Get member
[**OrganizationsSvcV1MembersMemberIdPatch**](MembersAPI.md#OrganizationsSvcV1MembersMemberIdPatch) | **Patch** /organizations-svc/v1/members/{member_id} | Update member



## OrganizationsSvcV1MembersGet

> MembersCollection OrganizationsSvcV1MembersGet(ctx).OrganizationId(organizationId).AccountId(accountId).Text(text).Include(include).PageLimit(pageLimit).PageOffset(pageOffset).Execute()

Get members



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
	organizationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Filter members by organization ID. (optional)
	accountId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Filter members by account ID. (optional)
	text := "text_example" // string | Text filter for members. (optional)
	include := []string{"Include_example"} // []string | Optional list of related resources to include in the response. Supported values: `profile`, `organization`.  (optional)
	pageLimit := int32(56) // int32 | Max number of items to return (1-100). (optional)
	pageOffset := int32(56) // int32 | Number of items to skip. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MembersAPI.OrganizationsSvcV1MembersGet(context.Background()).OrganizationId(organizationId).AccountId(accountId).Text(text).Include(include).PageLimit(pageLimit).PageOffset(pageOffset).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MembersAPI.OrganizationsSvcV1MembersGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OrganizationsSvcV1MembersGet`: MembersCollection
	fmt.Fprintf(os.Stdout, "Response from `MembersAPI.OrganizationsSvcV1MembersGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1MembersGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **organizationId** | **uuid.UUID** | Filter members by organization ID. | 
 **accountId** | **uuid.UUID** | Filter members by account ID. | 
 **text** | **string** | Text filter for members. | 
 **include** | **[]string** | Optional list of related resources to include in the response. Supported values: &#x60;profile&#x60;, &#x60;organization&#x60;.  | 
 **pageLimit** | **int32** | Max number of items to return (1-100). | 
 **pageOffset** | **int32** | Number of items to skip. | 

### Return type

[**MembersCollection**](MembersCollection.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OrganizationsSvcV1MembersMemberIdDelete

> OrganizationsSvcV1MembersMemberIdDelete(ctx, memberId).Execute()

Delete member



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
	memberId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Member ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.MembersAPI.OrganizationsSvcV1MembersMemberIdDelete(context.Background(), memberId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MembersAPI.OrganizationsSvcV1MembersMemberIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**memberId** | **uuid.UUID** | Member ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1MembersMemberIdDeleteRequest struct via the builder pattern


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


## OrganizationsSvcV1MembersMemberIdGet

> Member OrganizationsSvcV1MembersMemberIdGet(ctx, memberId).Include(include).Execute()

Get member



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
	memberId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Member ID
	include := []string{"Include_example"} // []string | Optional list of related resources to include in the response. Supported values: `profile`, `organizations`.  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MembersAPI.OrganizationsSvcV1MembersMemberIdGet(context.Background(), memberId).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MembersAPI.OrganizationsSvcV1MembersMemberIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OrganizationsSvcV1MembersMemberIdGet`: Member
	fmt.Fprintf(os.Stdout, "Response from `MembersAPI.OrganizationsSvcV1MembersMemberIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**memberId** | **uuid.UUID** | Member ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1MembersMemberIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **include** | **[]string** | Optional list of related resources to include in the response. Supported values: &#x60;profile&#x60;, &#x60;organizations&#x60;.  | 

### Return type

[**Member**](Member.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OrganizationsSvcV1MembersMemberIdPatch

> Member OrganizationsSvcV1MembersMemberIdPatch(ctx, memberId).UpdateMember(updateMember).Execute()

Update member



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
	memberId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Member ID
	updateMember := *openapiclient.NewUpdateMember(*openapiclient.NewUpdateMemberData("TODO", "Type_example", *openapiclient.NewUpdateMemberDataAttributes())) // UpdateMember | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MembersAPI.OrganizationsSvcV1MembersMemberIdPatch(context.Background(), memberId).UpdateMember(updateMember).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MembersAPI.OrganizationsSvcV1MembersMemberIdPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OrganizationsSvcV1MembersMemberIdPatch`: Member
	fmt.Fprintf(os.Stdout, "Response from `MembersAPI.OrganizationsSvcV1MembersMemberIdPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**memberId** | **uuid.UUID** | Member ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiOrganizationsSvcV1MembersMemberIdPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateMember** | [**UpdateMember**](UpdateMember.md) |  | 

### Return type

[**Member**](Member.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

