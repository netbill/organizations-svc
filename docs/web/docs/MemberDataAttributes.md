# MemberDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccountId** | [**uuid.UUID**](uuid.UUID.md) | The ID of the account associated with the member | 
**AgglomerationId** | [**uuid.UUID**](uuid.UUID.md) | The ID of the agglomeration the member belongs to | 
**Position** | Pointer to **string** | The position or role of the member within the agglomeration | [optional] 
**Label** | Pointer to **string** | A label or title associated with the member | [optional] 
**Username** | **string** | The username of the member | 
**Official** | **bool** | Indicates if the member is an official representative of the agglomeration | 
**CreatedAt** | **time.Time** | The date and time when the member was created | 
**UpdatedAt** | **time.Time** | The date and time when the member was last updated | 

## Methods

### NewMemberDataAttributes

`func NewMemberDataAttributes(accountId uuid.UUID, agglomerationId uuid.UUID, username string, official bool, createdAt time.Time, updatedAt time.Time, ) *MemberDataAttributes`

NewMemberDataAttributes instantiates a new MemberDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMemberDataAttributesWithDefaults

`func NewMemberDataAttributesWithDefaults() *MemberDataAttributes`

NewMemberDataAttributesWithDefaults instantiates a new MemberDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccountId

`func (o *MemberDataAttributes) GetAccountId() uuid.UUID`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *MemberDataAttributes) GetAccountIdOk() (*uuid.UUID, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *MemberDataAttributes) SetAccountId(v uuid.UUID)`

SetAccountId sets AccountId field to given value.


### GetAgglomerationId

`func (o *MemberDataAttributes) GetAgglomerationId() uuid.UUID`

GetAgglomerationId returns the AgglomerationId field if non-nil, zero value otherwise.

### GetAgglomerationIdOk

`func (o *MemberDataAttributes) GetAgglomerationIdOk() (*uuid.UUID, bool)`

GetAgglomerationIdOk returns a tuple with the AgglomerationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAgglomerationId

`func (o *MemberDataAttributes) SetAgglomerationId(v uuid.UUID)`

SetAgglomerationId sets AgglomerationId field to given value.


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

### GetUsername

`func (o *MemberDataAttributes) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *MemberDataAttributes) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *MemberDataAttributes) SetUsername(v string)`

SetUsername sets Username field to given value.


### GetOfficial

`func (o *MemberDataAttributes) GetOfficial() bool`

GetOfficial returns the Official field if non-nil, zero value otherwise.

### GetOfficialOk

`func (o *MemberDataAttributes) GetOfficialOk() (*bool, bool)`

GetOfficialOk returns a tuple with the Official field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOfficial

`func (o *MemberDataAttributes) SetOfficial(v bool)`

SetOfficial sets Official field to given value.


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


