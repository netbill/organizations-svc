# Invite

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**InviteData**](InviteData.md) |  | 
**Included** | Pointer to [**[]InviteIncludedInner**](InviteIncludedInner.md) |  | [optional] 

## Methods

### NewInvite

`func NewInvite(data InviteData, ) *Invite`

NewInvite instantiates a new Invite object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInviteWithDefaults

`func NewInviteWithDefaults() *Invite`

NewInviteWithDefaults instantiates a new Invite object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *Invite) GetData() InviteData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *Invite) GetDataOk() (*InviteData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *Invite) SetData(v InviteData)`

SetData sets Data field to given value.


### GetIncluded

`func (o *Invite) GetIncluded() []InviteIncludedInner`

GetIncluded returns the Included field if non-nil, zero value otherwise.

### GetIncludedOk

`func (o *Invite) GetIncludedOk() (*[]InviteIncludedInner, bool)`

GetIncludedOk returns a tuple with the Included field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIncluded

`func (o *Invite) SetIncluded(v []InviteIncludedInner)`

SetIncluded sets Included field to given value.

### HasIncluded

`func (o *Invite) HasIncluded() bool`

HasIncluded returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


