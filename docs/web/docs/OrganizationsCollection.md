# OrganizationsCollection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]OrganizationData**](OrganizationData.md) |  | 
**Included** | Pointer to [**[]MemberData**](MemberData.md) |  | [optional] 
**Links** | [**PaginationData**](PaginationData.md) |  | 

## Methods

### NewOrganizationsCollection

`func NewOrganizationsCollection(data []OrganizationData, links PaginationData, ) *OrganizationsCollection`

NewOrganizationsCollection instantiates a new OrganizationsCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOrganizationsCollectionWithDefaults

`func NewOrganizationsCollectionWithDefaults() *OrganizationsCollection`

NewOrganizationsCollectionWithDefaults instantiates a new OrganizationsCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *OrganizationsCollection) GetData() []OrganizationData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *OrganizationsCollection) GetDataOk() (*[]OrganizationData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *OrganizationsCollection) SetData(v []OrganizationData)`

SetData sets Data field to given value.


### GetIncluded

`func (o *OrganizationsCollection) GetIncluded() []MemberData`

GetIncluded returns the Included field if non-nil, zero value otherwise.

### GetIncludedOk

`func (o *OrganizationsCollection) GetIncludedOk() (*[]MemberData, bool)`

GetIncludedOk returns a tuple with the Included field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIncluded

`func (o *OrganizationsCollection) SetIncluded(v []MemberData)`

SetIncluded sets Included field to given value.

### HasIncluded

`func (o *OrganizationsCollection) HasIncluded() bool`

HasIncluded returns a boolean if a field has been set.

### GetLinks

`func (o *OrganizationsCollection) GetLinks() PaginationData`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *OrganizationsCollection) GetLinksOk() (*PaginationData, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *OrganizationsCollection) SetLinks(v PaginationData)`

SetLinks sets Links field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


