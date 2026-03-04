# MembersCollection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]MemberData**](MemberData.md) |  | 
**Included** | Pointer to [**[]MemberIncludedInner**](MemberIncludedInner.md) |  | [optional] 
**Links** | [**PaginationData**](PaginationData.md) |  | 

## Methods

### NewMembersCollection

`func NewMembersCollection(data []MemberData, links PaginationData, ) *MembersCollection`

NewMembersCollection instantiates a new MembersCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMembersCollectionWithDefaults

`func NewMembersCollectionWithDefaults() *MembersCollection`

NewMembersCollectionWithDefaults instantiates a new MembersCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *MembersCollection) GetData() []MemberData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *MembersCollection) GetDataOk() (*[]MemberData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *MembersCollection) SetData(v []MemberData)`

SetData sets Data field to given value.


### GetIncluded

`func (o *MembersCollection) GetIncluded() []MemberIncludedInner`

GetIncluded returns the Included field if non-nil, zero value otherwise.

### GetIncludedOk

`func (o *MembersCollection) GetIncludedOk() (*[]MemberIncludedInner, bool)`

GetIncludedOk returns a tuple with the Included field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIncluded

`func (o *MembersCollection) SetIncluded(v []MemberIncludedInner)`

SetIncluded sets Included field to given value.

### HasIncluded

`func (o *MembersCollection) HasIncluded() bool`

HasIncluded returns a boolean if a field has been set.

### GetLinks

`func (o *MembersCollection) GetLinks() PaginationData`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *MembersCollection) GetLinksOk() (*PaginationData, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *MembersCollection) SetLinks(v PaginationData)`

SetLinks sets Links field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


