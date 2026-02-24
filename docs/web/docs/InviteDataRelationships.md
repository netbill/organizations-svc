# InviteDataRelationships

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Organization** | [**MemberDataRelationshipsOrganization**](MemberDataRelationshipsOrganization.md) |  | 
**Invitee** | [**MemberDataRelationshipsProfile**](MemberDataRelationshipsProfile.md) |  | 

## Methods

### NewInviteDataRelationships

`func NewInviteDataRelationships(organization MemberDataRelationshipsOrganization, invitee MemberDataRelationshipsProfile, ) *InviteDataRelationships`

NewInviteDataRelationships instantiates a new InviteDataRelationships object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInviteDataRelationshipsWithDefaults

`func NewInviteDataRelationshipsWithDefaults() *InviteDataRelationships`

NewInviteDataRelationshipsWithDefaults instantiates a new InviteDataRelationships object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOrganization

`func (o *InviteDataRelationships) GetOrganization() MemberDataRelationshipsOrganization`

GetOrganization returns the Organization field if non-nil, zero value otherwise.

### GetOrganizationOk

`func (o *InviteDataRelationships) GetOrganizationOk() (*MemberDataRelationshipsOrganization, bool)`

GetOrganizationOk returns a tuple with the Organization field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganization

`func (o *InviteDataRelationships) SetOrganization(v MemberDataRelationshipsOrganization)`

SetOrganization sets Organization field to given value.


### GetInvitee

`func (o *InviteDataRelationships) GetInvitee() MemberDataRelationshipsProfile`

GetInvitee returns the Invitee field if non-nil, zero value otherwise.

### GetInviteeOk

`func (o *InviteDataRelationships) GetInviteeOk() (*MemberDataRelationshipsProfile, bool)`

GetInviteeOk returns a tuple with the Invitee field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInvitee

`func (o *InviteDataRelationships) SetInvitee(v MemberDataRelationshipsProfile)`

SetInvitee sets Invitee field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


