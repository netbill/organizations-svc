# CreateInviteDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**OrganizationId** | [**uuid.UUID**](uuid.UUID.md) | The ID of the organization to which the invite belongs | 
**AccountId** | [**uuid.UUID**](uuid.UUID.md) | The ID of the account that was invited | 

## Methods

### NewCreateInviteDataAttributes

`func NewCreateInviteDataAttributes(organizationId uuid.UUID, accountId uuid.UUID, ) *CreateInviteDataAttributes`

NewCreateInviteDataAttributes instantiates a new CreateInviteDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateInviteDataAttributesWithDefaults

`func NewCreateInviteDataAttributesWithDefaults() *CreateInviteDataAttributes`

NewCreateInviteDataAttributesWithDefaults instantiates a new CreateInviteDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOrganizationId

`func (o *CreateInviteDataAttributes) GetOrganizationId() uuid.UUID`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *CreateInviteDataAttributes) GetOrganizationIdOk() (*uuid.UUID, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *CreateInviteDataAttributes) SetOrganizationId(v uuid.UUID)`

SetOrganizationId sets OrganizationId field to given value.


### GetAccountId

`func (o *CreateInviteDataAttributes) GetAccountId() uuid.UUID`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *CreateInviteDataAttributes) GetAccountIdOk() (*uuid.UUID, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *CreateInviteDataAttributes) SetAccountId(v uuid.UUID)`

SetAccountId sets AccountId field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


