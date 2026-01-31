# UpdateOrganizationDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The name of the organization | 
**DeleteIcon** | **bool** | Flag to indicate if the organization&#39;s icon should be deleted | 
**DeleteBanner** | **bool** | Flag to indicate if the organization&#39;s banner should be deleted | 

## Methods

### NewUpdateOrganizationDataAttributes

`func NewUpdateOrganizationDataAttributes(name string, deleteIcon bool, deleteBanner bool, ) *UpdateOrganizationDataAttributes`

NewUpdateOrganizationDataAttributes instantiates a new UpdateOrganizationDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateOrganizationDataAttributesWithDefaults

`func NewUpdateOrganizationDataAttributesWithDefaults() *UpdateOrganizationDataAttributes`

NewUpdateOrganizationDataAttributesWithDefaults instantiates a new UpdateOrganizationDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *UpdateOrganizationDataAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *UpdateOrganizationDataAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *UpdateOrganizationDataAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetDeleteIcon

`func (o *UpdateOrganizationDataAttributes) GetDeleteIcon() bool`

GetDeleteIcon returns the DeleteIcon field if non-nil, zero value otherwise.

### GetDeleteIconOk

`func (o *UpdateOrganizationDataAttributes) GetDeleteIconOk() (*bool, bool)`

GetDeleteIconOk returns a tuple with the DeleteIcon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeleteIcon

`func (o *UpdateOrganizationDataAttributes) SetDeleteIcon(v bool)`

SetDeleteIcon sets DeleteIcon field to given value.


### GetDeleteBanner

`func (o *UpdateOrganizationDataAttributes) GetDeleteBanner() bool`

GetDeleteBanner returns the DeleteBanner field if non-nil, zero value otherwise.

### GetDeleteBannerOk

`func (o *UpdateOrganizationDataAttributes) GetDeleteBannerOk() (*bool, bool)`

GetDeleteBannerOk returns a tuple with the DeleteBanner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeleteBanner

`func (o *UpdateOrganizationDataAttributes) SetDeleteBanner(v bool)`

SetDeleteBanner sets DeleteBanner field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


