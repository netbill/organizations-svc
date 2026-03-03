# OrganizationDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | **string** | The status of the organization | 
**Name** | **string** | The name of the organization | 
**IconKey** | Pointer to **string** | The media key for the organization&#39;s icon | [optional] 
**BannerKey** | Pointer to **string** | The media key for the organization&#39;s banner | [optional] 
**Version** | **int32** | The version number of the organization, used for optimistic concurrency control | 
**CreatedAt** | **time.Time** | The date and time when the organization was created | 
**UpdatedAt** | **time.Time** | The date and time when the organization was last updated | 

## Methods

### NewOrganizationDataAttributes

`func NewOrganizationDataAttributes(status string, name string, version int32, createdAt time.Time, updatedAt time.Time, ) *OrganizationDataAttributes`

NewOrganizationDataAttributes instantiates a new OrganizationDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOrganizationDataAttributesWithDefaults

`func NewOrganizationDataAttributesWithDefaults() *OrganizationDataAttributes`

NewOrganizationDataAttributesWithDefaults instantiates a new OrganizationDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *OrganizationDataAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *OrganizationDataAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *OrganizationDataAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetName

`func (o *OrganizationDataAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *OrganizationDataAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *OrganizationDataAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetIconKey

`func (o *OrganizationDataAttributes) GetIconKey() string`

GetIconKey returns the IconKey field if non-nil, zero value otherwise.

### GetIconKeyOk

`func (o *OrganizationDataAttributes) GetIconKeyOk() (*string, bool)`

GetIconKeyOk returns a tuple with the IconKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIconKey

`func (o *OrganizationDataAttributes) SetIconKey(v string)`

SetIconKey sets IconKey field to given value.

### HasIconKey

`func (o *OrganizationDataAttributes) HasIconKey() bool`

HasIconKey returns a boolean if a field has been set.

### GetBannerKey

`func (o *OrganizationDataAttributes) GetBannerKey() string`

GetBannerKey returns the BannerKey field if non-nil, zero value otherwise.

### GetBannerKeyOk

`func (o *OrganizationDataAttributes) GetBannerKeyOk() (*string, bool)`

GetBannerKeyOk returns a tuple with the BannerKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBannerKey

`func (o *OrganizationDataAttributes) SetBannerKey(v string)`

SetBannerKey sets BannerKey field to given value.

### HasBannerKey

`func (o *OrganizationDataAttributes) HasBannerKey() bool`

HasBannerKey returns a boolean if a field has been set.

### GetVersion

`func (o *OrganizationDataAttributes) GetVersion() int32`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *OrganizationDataAttributes) GetVersionOk() (*int32, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *OrganizationDataAttributes) SetVersion(v int32)`

SetVersion sets Version field to given value.


### GetCreatedAt

`func (o *OrganizationDataAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *OrganizationDataAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *OrganizationDataAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUpdatedAt

`func (o *OrganizationDataAttributes) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *OrganizationDataAttributes) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *OrganizationDataAttributes) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


