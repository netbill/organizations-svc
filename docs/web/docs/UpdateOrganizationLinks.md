# UpdateOrganizationLinks

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**UpdateOrganizationLinksData**](UpdateOrganizationLinksData.md) |  | 
**Included** | [**[]OrganizationData**](OrganizationData.md) |  | 

## Methods

### NewUpdateOrganizationLinks

`func NewUpdateOrganizationLinks(data UpdateOrganizationLinksData, included []OrganizationData, ) *UpdateOrganizationLinks`

NewUpdateOrganizationLinks instantiates a new UpdateOrganizationLinks object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateOrganizationLinksWithDefaults

`func NewUpdateOrganizationLinksWithDefaults() *UpdateOrganizationLinks`

NewUpdateOrganizationLinksWithDefaults instantiates a new UpdateOrganizationLinks object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *UpdateOrganizationLinks) GetData() UpdateOrganizationLinksData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *UpdateOrganizationLinks) GetDataOk() (*UpdateOrganizationLinksData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *UpdateOrganizationLinks) SetData(v UpdateOrganizationLinksData)`

SetData sets Data field to given value.


### GetIncluded

`func (o *UpdateOrganizationLinks) GetIncluded() []OrganizationData`

GetIncluded returns the Included field if non-nil, zero value otherwise.

### GetIncludedOk

`func (o *UpdateOrganizationLinks) GetIncludedOk() (*[]OrganizationData, bool)`

GetIncludedOk returns a tuple with the Included field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIncluded

`func (o *UpdateOrganizationLinks) SetIncluded(v []OrganizationData)`

SetIncluded sets Included field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


