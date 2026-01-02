# RoleDataRelationshipsPermissionsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | Unique identifier for the permission | 
**Code** | **string** | A short code representing the permission | 
**Description** | **string** | A detailed description of what the permission allows | 
**Enabled** | **bool** | Indicates if the role has this permission | 

## Methods

### NewRoleDataRelationshipsPermissionsInner

`func NewRoleDataRelationshipsPermissionsInner(id uuid.UUID, code string, description string, enabled bool, ) *RoleDataRelationshipsPermissionsInner`

NewRoleDataRelationshipsPermissionsInner instantiates a new RoleDataRelationshipsPermissionsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRoleDataRelationshipsPermissionsInnerWithDefaults

`func NewRoleDataRelationshipsPermissionsInnerWithDefaults() *RoleDataRelationshipsPermissionsInner`

NewRoleDataRelationshipsPermissionsInnerWithDefaults instantiates a new RoleDataRelationshipsPermissionsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *RoleDataRelationshipsPermissionsInner) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *RoleDataRelationshipsPermissionsInner) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *RoleDataRelationshipsPermissionsInner) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetCode

`func (o *RoleDataRelationshipsPermissionsInner) GetCode() string`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *RoleDataRelationshipsPermissionsInner) GetCodeOk() (*string, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *RoleDataRelationshipsPermissionsInner) SetCode(v string)`

SetCode sets Code field to given value.


### GetDescription

`func (o *RoleDataRelationshipsPermissionsInner) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *RoleDataRelationshipsPermissionsInner) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *RoleDataRelationshipsPermissionsInner) SetDescription(v string)`

SetDescription sets Description field to given value.


### GetEnabled

`func (o *RoleDataRelationshipsPermissionsInner) GetEnabled() bool`

GetEnabled returns the Enabled field if non-nil, zero value otherwise.

### GetEnabledOk

`func (o *RoleDataRelationshipsPermissionsInner) GetEnabledOk() (*bool, bool)`

GetEnabledOk returns a tuple with the Enabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnabled

`func (o *RoleDataRelationshipsPermissionsInner) SetEnabled(v bool)`

SetEnabled sets Enabled field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


