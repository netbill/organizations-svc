# RoleData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | member ID | 
**Type** | **string** |  | 
**Attributes** | [**RoleDataAttributes**](RoleDataAttributes.md) |  | 
**Relationships** | Pointer to [**RoleDataRelationships**](RoleDataRelationships.md) |  | [optional] 

## Methods

### NewRoleData

`func NewRoleData(id uuid.UUID, type_ string, attributes RoleDataAttributes, ) *RoleData`

NewRoleData instantiates a new RoleData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRoleDataWithDefaults

`func NewRoleDataWithDefaults() *RoleData`

NewRoleDataWithDefaults instantiates a new RoleData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *RoleData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *RoleData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *RoleData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *RoleData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *RoleData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *RoleData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *RoleData) GetAttributes() RoleDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *RoleData) GetAttributesOk() (*RoleDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *RoleData) SetAttributes(v RoleDataAttributes)`

SetAttributes sets Attributes field to given value.


### GetRelationships

`func (o *RoleData) GetRelationships() RoleDataRelationships`

GetRelationships returns the Relationships field if non-nil, zero value otherwise.

### GetRelationshipsOk

`func (o *RoleData) GetRelationshipsOk() (*RoleDataRelationships, bool)`

GetRelationshipsOk returns a tuple with the Relationships field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRelationships

`func (o *RoleData) SetRelationships(v RoleDataRelationships)`

SetRelationships sets Relationships field to given value.

### HasRelationships

`func (o *RoleData) HasRelationships() bool`

HasRelationships returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


