# CreateRoleData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | **string** |  | 
**Attributes** | [**CreateRoleDataAttributes**](CreateRoleDataAttributes.md) |  | 

## Methods

### NewCreateRoleData

`func NewCreateRoleData(type_ string, attributes CreateRoleDataAttributes, ) *CreateRoleData`

NewCreateRoleData instantiates a new CreateRoleData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateRoleDataWithDefaults

`func NewCreateRoleDataWithDefaults() *CreateRoleData`

NewCreateRoleDataWithDefaults instantiates a new CreateRoleData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *CreateRoleData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *CreateRoleData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *CreateRoleData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *CreateRoleData) GetAttributes() CreateRoleDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *CreateRoleData) GetAttributesOk() (*CreateRoleDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *CreateRoleData) SetAttributes(v CreateRoleDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


