# UploadResourcesLink

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Key** | **string** | Key for the avatar in the storage service | 
**UploadUrl** | **string** | Pre-signed URL for uploading the avatar | 
**PreloadUrl** | **string** | URL for preloading the avatar after upload | 

## Methods

### NewUploadResourcesLink

`func NewUploadResourcesLink(key string, uploadUrl string, preloadUrl string, ) *UploadResourcesLink`

NewUploadResourcesLink instantiates a new UploadResourcesLink object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUploadResourcesLinkWithDefaults

`func NewUploadResourcesLinkWithDefaults() *UploadResourcesLink`

NewUploadResourcesLinkWithDefaults instantiates a new UploadResourcesLink object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKey

`func (o *UploadResourcesLink) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *UploadResourcesLink) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *UploadResourcesLink) SetKey(v string)`

SetKey sets Key field to given value.


### GetUploadUrl

`func (o *UploadResourcesLink) GetUploadUrl() string`

GetUploadUrl returns the UploadUrl field if non-nil, zero value otherwise.

### GetUploadUrlOk

`func (o *UploadResourcesLink) GetUploadUrlOk() (*string, bool)`

GetUploadUrlOk returns a tuple with the UploadUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUploadUrl

`func (o *UploadResourcesLink) SetUploadUrl(v string)`

SetUploadUrl sets UploadUrl field to given value.


### GetPreloadUrl

`func (o *UploadResourcesLink) GetPreloadUrl() string`

GetPreloadUrl returns the PreloadUrl field if non-nil, zero value otherwise.

### GetPreloadUrlOk

`func (o *UploadResourcesLink) GetPreloadUrlOk() (*string, bool)`

GetPreloadUrlOk returns a tuple with the PreloadUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreloadUrl

`func (o *UploadResourcesLink) SetPreloadUrl(v string)`

SetPreloadUrl sets PreloadUrl field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


