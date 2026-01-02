# CreateRoleDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**OrganizationId** | [**uuid.UUID**](uuid.UUID.md) | The ID of the organization this role belongs to | 
**Rank** | **uint** | The rank of the role within the organization | 
**Name** | **string** | The name of the role | 
**Description** | **string** | A brief description of the role | 
**Color** | **string** | The color associated with the role in HEX format | 

## Methods

### NewCreateRoleDataAttributes

`func NewCreateRoleDataAttributes(organizationId uuid.UUID, rank uint, name string, description string, color string, ) *CreateRoleDataAttributes`

NewCreateRoleDataAttributes instantiates a new CreateRoleDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateRoleDataAttributesWithDefaults

`func NewCreateRoleDataAttributesWithDefaults() *CreateRoleDataAttributes`

NewCreateRoleDataAttributesWithDefaults instantiates a new CreateRoleDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOrganizationId

`func (o *CreateRoleDataAttributes) GetOrganizationId() uuid.UUID`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *CreateRoleDataAttributes) GetOrganizationIdOk() (*uuid.UUID, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *CreateRoleDataAttributes) SetOrganizationId(v uuid.UUID)`

SetOrganizationId sets OrganizationId field to given value.


### GetRank

`func (o *CreateRoleDataAttributes) GetRank() uint`

GetRank returns the Rank field if non-nil, zero value otherwise.

### GetRankOk

`func (o *CreateRoleDataAttributes) GetRankOk() (*uint, bool)`

GetRankOk returns a tuple with the Rank field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRank

`func (o *CreateRoleDataAttributes) SetRank(v uint)`

SetRank sets Rank field to given value.


### GetName

`func (o *CreateRoleDataAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateRoleDataAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateRoleDataAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *CreateRoleDataAttributes) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *CreateRoleDataAttributes) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *CreateRoleDataAttributes) SetDescription(v string)`

SetDescription sets Description field to given value.


### GetColor

`func (o *CreateRoleDataAttributes) GetColor() string`

GetColor returns the Color field if non-nil, zero value otherwise.

### GetColorOk

`func (o *CreateRoleDataAttributes) GetColorOk() (*string, bool)`

GetColorOk returns a tuple with the Color field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetColor

`func (o *CreateRoleDataAttributes) SetColor(v string)`

SetColor sets Color field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


