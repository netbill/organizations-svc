# CreateRoleDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AgglomerationId** | [**uuid.UUID**](uuid.UUID.md) | The ID of the agglomeration this role belongs to | 
**Rank** | **uint** | The rank of the role within the agglomeration | 
**Name** | **string** | The name of the role | 
**Description** | **string** | A brief description of the role | 
**Color** | **string** | The color associated with the role in HEX format | 

## Methods

### NewCreateRoleDataAttributes

`func NewCreateRoleDataAttributes(agglomerationId uuid.UUID, rank uint, name string, description string, color string, ) *CreateRoleDataAttributes`

NewCreateRoleDataAttributes instantiates a new CreateRoleDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateRoleDataAttributesWithDefaults

`func NewCreateRoleDataAttributesWithDefaults() *CreateRoleDataAttributes`

NewCreateRoleDataAttributesWithDefaults instantiates a new CreateRoleDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAgglomerationId

`func (o *CreateRoleDataAttributes) GetAgglomerationId() uuid.UUID`

GetAgglomerationId returns the AgglomerationId field if non-nil, zero value otherwise.

### GetAgglomerationIdOk

`func (o *CreateRoleDataAttributes) GetAgglomerationIdOk() (*uuid.UUID, bool)`

GetAgglomerationIdOk returns a tuple with the AgglomerationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAgglomerationId

`func (o *CreateRoleDataAttributes) SetAgglomerationId(v uuid.UUID)`

SetAgglomerationId sets AgglomerationId field to given value.


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


