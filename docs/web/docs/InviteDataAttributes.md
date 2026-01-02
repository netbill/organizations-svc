# InviteDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**OrganizationId** | [**uuid.UUID**](uuid.UUID.md) | The ID of the organization to which the invite belongs | 
**AccountId** | [**uuid.UUID**](uuid.UUID.md) | The ID of the account that was invited | 
**Status** | **string** | The status of the invite | 
**ExpiresAt** | **time.Time** | The expiration date and time of the invite | 
**CreatedAt** | **time.Time** | The date and time when the invite was created | 

## Methods

### NewInviteDataAttributes

`func NewInviteDataAttributes(organizationId uuid.UUID, accountId uuid.UUID, status string, expiresAt time.Time, createdAt time.Time, ) *InviteDataAttributes`

NewInviteDataAttributes instantiates a new InviteDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInviteDataAttributesWithDefaults

`func NewInviteDataAttributesWithDefaults() *InviteDataAttributes`

NewInviteDataAttributesWithDefaults instantiates a new InviteDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOrganizationId

`func (o *InviteDataAttributes) GetOrganizationId() uuid.UUID`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *InviteDataAttributes) GetOrganizationIdOk() (*uuid.UUID, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *InviteDataAttributes) SetOrganizationId(v uuid.UUID)`

SetOrganizationId sets OrganizationId field to given value.


### GetAccountId

`func (o *InviteDataAttributes) GetAccountId() uuid.UUID`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *InviteDataAttributes) GetAccountIdOk() (*uuid.UUID, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *InviteDataAttributes) SetAccountId(v uuid.UUID)`

SetAccountId sets AccountId field to given value.


### GetStatus

`func (o *InviteDataAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *InviteDataAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *InviteDataAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetExpiresAt

`func (o *InviteDataAttributes) GetExpiresAt() time.Time`

GetExpiresAt returns the ExpiresAt field if non-nil, zero value otherwise.

### GetExpiresAtOk

`func (o *InviteDataAttributes) GetExpiresAtOk() (*time.Time, bool)`

GetExpiresAtOk returns a tuple with the ExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresAt

`func (o *InviteDataAttributes) SetExpiresAt(v time.Time)`

SetExpiresAt sets ExpiresAt field to given value.


### GetCreatedAt

`func (o *InviteDataAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *InviteDataAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *InviteDataAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


