# RolesCollection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]RoleData**](RoleData.md) |  | 
**Links** | [**PaginationData**](PaginationData.md) |  | 

## Methods

### NewRolesCollection

`func NewRolesCollection(data []RoleData, links PaginationData, ) *RolesCollection`

NewRolesCollection instantiates a new RolesCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRolesCollectionWithDefaults

`func NewRolesCollectionWithDefaults() *RolesCollection`

NewRolesCollectionWithDefaults instantiates a new RolesCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *RolesCollection) GetData() []RoleData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *RolesCollection) GetDataOk() (*[]RoleData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *RolesCollection) SetData(v []RoleData)`

SetData sets Data field to given value.


### GetLinks

`func (o *RolesCollection) GetLinks() PaginationData`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *RolesCollection) GetLinksOk() (*PaginationData, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *RolesCollection) SetLinks(v PaginationData)`

SetLinks sets Links field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


