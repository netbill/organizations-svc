# RolePermissionsDataInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | Unique identifier for the permission | 
**Code** | **string** | A short code representing the permission | 
**Description** | **string** | A detailed description of what the permission allows | 

## Methods

### NewRolePermissionsDataInner

`func NewRolePermissionsDataInner(id uuid.UUID, code string, description string, ) *RolePermissionsDataInner`

NewRolePermissionsDataInner instantiates a new RolePermissionsDataInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRolePermissionsDataInnerWithDefaults

`func NewRolePermissionsDataInnerWithDefaults() *RolePermissionsDataInner`

NewRolePermissionsDataInnerWithDefaults instantiates a new RolePermissionsDataInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *RolePermissionsDataInner) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *RolePermissionsDataInner) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *RolePermissionsDataInner) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetCode

`func (o *RolePermissionsDataInner) GetCode() string`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *RolePermissionsDataInner) GetCodeOk() (*string, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *RolePermissionsDataInner) SetCode(v string)`

SetCode sets Code field to given value.


### GetDescription

`func (o *RolePermissionsDataInner) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *RolePermissionsDataInner) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *RolePermissionsDataInner) SetDescription(v string)`

SetDescription sets Description field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


