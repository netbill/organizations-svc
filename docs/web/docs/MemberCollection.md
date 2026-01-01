# MemberCollection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]MemberData**](MemberData.md) |  | 
**Links** | [**PaginationData**](PaginationData.md) |  | 

## Methods

### NewMemberCollection

`func NewMemberCollection(data []MemberData, links PaginationData, ) *MemberCollection`

NewMemberCollection instantiates a new MemberCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMemberCollectionWithDefaults

`func NewMemberCollectionWithDefaults() *MemberCollection`

NewMemberCollectionWithDefaults instantiates a new MemberCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *MemberCollection) GetData() []MemberData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *MemberCollection) GetDataOk() (*[]MemberData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *MemberCollection) SetData(v []MemberData)`

SetData sets Data field to given value.


### GetLinks

`func (o *MemberCollection) GetLinks() PaginationData`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *MemberCollection) GetLinksOk() (*PaginationData, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *MemberCollection) SetLinks(v PaginationData)`

SetLinks sets Links field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


