# SentInviteDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**OrganizationId** | [**uuid.UUID**](uuid.UUID.md) | The ID of the organization to which the invite belongs | 
**AccountId** | [**uuid.UUID**](uuid.UUID.md) | The ID of the account that was invited | 

## Methods

### NewSentInviteDataAttributes

`func NewSentInviteDataAttributes(organizationId uuid.UUID, accountId uuid.UUID, ) *SentInviteDataAttributes`

NewSentInviteDataAttributes instantiates a new SentInviteDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSentInviteDataAttributesWithDefaults

`func NewSentInviteDataAttributesWithDefaults() *SentInviteDataAttributes`

NewSentInviteDataAttributesWithDefaults instantiates a new SentInviteDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOrganizationId

`func (o *SentInviteDataAttributes) GetOrganizationId() uuid.UUID`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *SentInviteDataAttributes) GetOrganizationIdOk() (*uuid.UUID, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *SentInviteDataAttributes) SetOrganizationId(v uuid.UUID)`

SetOrganizationId sets OrganizationId field to given value.


### GetAccountId

`func (o *SentInviteDataAttributes) GetAccountId() uuid.UUID`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *SentInviteDataAttributes) GetAccountIdOk() (*uuid.UUID, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *SentInviteDataAttributes) SetAccountId(v uuid.UUID)`

SetAccountId sets AccountId field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


