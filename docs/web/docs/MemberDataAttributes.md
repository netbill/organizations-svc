# MemberDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Head** | **bool** | Indicates if the member is the head of the organization | 
**Position** | Pointer to **string** | The position or role of the member within the organization | [optional] 
**Label** | Pointer to **string** | A label or title associated with the member | [optional] 
**Version** | **int32** | The version number of the member record, used for optimistic locking | 
**CreatedAt** | **time.Time** | The date and time when the member was created | 
**UpdatedAt** | **time.Time** | The date and time when the member was last updated | 

## Methods

### NewMemberDataAttributes

`func NewMemberDataAttributes(head bool, version int32, createdAt time.Time, updatedAt time.Time, ) *MemberDataAttributes`

NewMemberDataAttributes instantiates a new MemberDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMemberDataAttributesWithDefaults

`func NewMemberDataAttributesWithDefaults() *MemberDataAttributes`

NewMemberDataAttributesWithDefaults instantiates a new MemberDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHead

`func (o *MemberDataAttributes) GetHead() bool`

GetHead returns the Head field if non-nil, zero value otherwise.

### GetHeadOk

`func (o *MemberDataAttributes) GetHeadOk() (*bool, bool)`

GetHeadOk returns a tuple with the Head field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHead

`func (o *MemberDataAttributes) SetHead(v bool)`

SetHead sets Head field to given value.


### GetPosition

`func (o *MemberDataAttributes) GetPosition() string`

GetPosition returns the Position field if non-nil, zero value otherwise.

### GetPositionOk

`func (o *MemberDataAttributes) GetPositionOk() (*string, bool)`

GetPositionOk returns a tuple with the Position field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPosition

`func (o *MemberDataAttributes) SetPosition(v string)`

SetPosition sets Position field to given value.

### HasPosition

`func (o *MemberDataAttributes) HasPosition() bool`

HasPosition returns a boolean if a field has been set.

### GetLabel

`func (o *MemberDataAttributes) GetLabel() string`

GetLabel returns the Label field if non-nil, zero value otherwise.

### GetLabelOk

`func (o *MemberDataAttributes) GetLabelOk() (*string, bool)`

GetLabelOk returns a tuple with the Label field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabel

`func (o *MemberDataAttributes) SetLabel(v string)`

SetLabel sets Label field to given value.

### HasLabel

`func (o *MemberDataAttributes) HasLabel() bool`

HasLabel returns a boolean if a field has been set.

### GetVersion

`func (o *MemberDataAttributes) GetVersion() int32`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *MemberDataAttributes) GetVersionOk() (*int32, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *MemberDataAttributes) SetVersion(v int32)`

SetVersion sets Version field to given value.


### GetCreatedAt

`func (o *MemberDataAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *MemberDataAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *MemberDataAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUpdatedAt

`func (o *MemberDataAttributes) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *MemberDataAttributes) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *MemberDataAttributes) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


