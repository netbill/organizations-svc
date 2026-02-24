# Member

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**MemberData**](MemberData.md) |  | 
**Included** | Pointer to [**[]MemberIncludedInner**](MemberIncludedInner.md) | Included related resources (profile, organization) | [optional] 

## Methods

### NewMember

`func NewMember(data MemberData, ) *Member`

NewMember instantiates a new Member object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMemberWithDefaults

`func NewMemberWithDefaults() *Member`

NewMemberWithDefaults instantiates a new Member object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *Member) GetData() MemberData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *Member) GetDataOk() (*MemberData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *Member) SetData(v MemberData)`

SetData sets Data field to given value.


### GetIncluded

`func (o *Member) GetIncluded() []MemberIncludedInner`

GetIncluded returns the Included field if non-nil, zero value otherwise.

### GetIncludedOk

`func (o *Member) GetIncludedOk() (*[]MemberIncludedInner, bool)`

GetIncludedOk returns a tuple with the Included field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIncluded

`func (o *Member) SetIncluded(v []MemberIncludedInner)`

SetIncluded sets Included field to given value.

### HasIncluded

`func (o *Member) HasIncluded() bool`

HasIncluded returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


