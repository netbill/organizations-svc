# UpdateRolePermissionsData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | role ID | 
**Type** | **string** |  | 
**Attributes** | [**UpdateRolePermissionsDataAttributes**](UpdateRolePermissionsDataAttributes.md) |  | 

## Methods

### NewUpdateRolePermissionsData

`func NewUpdateRolePermissionsData(id uuid.UUID, type_ string, attributes UpdateRolePermissionsDataAttributes, ) *UpdateRolePermissionsData`

NewUpdateRolePermissionsData instantiates a new UpdateRolePermissionsData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateRolePermissionsDataWithDefaults

`func NewUpdateRolePermissionsDataWithDefaults() *UpdateRolePermissionsData`

NewUpdateRolePermissionsDataWithDefaults instantiates a new UpdateRolePermissionsData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UpdateRolePermissionsData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UpdateRolePermissionsData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UpdateRolePermissionsData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *UpdateRolePermissionsData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UpdateRolePermissionsData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UpdateRolePermissionsData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *UpdateRolePermissionsData) GetAttributes() UpdateRolePermissionsDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *UpdateRolePermissionsData) GetAttributesOk() (*UpdateRolePermissionsDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *UpdateRolePermissionsData) SetAttributes(v UpdateRolePermissionsDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


