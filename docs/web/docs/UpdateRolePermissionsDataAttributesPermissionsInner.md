# UpdateRolePermissionsDataAttributesPermissionsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | permission id | 
**Status** | **bool** | The new status of the permission for the role | 

## Methods

### NewUpdateRolePermissionsDataAttributesPermissionsInner

`func NewUpdateRolePermissionsDataAttributesPermissionsInner(id uuid.UUID, status bool, ) *UpdateRolePermissionsDataAttributesPermissionsInner`

NewUpdateRolePermissionsDataAttributesPermissionsInner instantiates a new UpdateRolePermissionsDataAttributesPermissionsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateRolePermissionsDataAttributesPermissionsInnerWithDefaults

`func NewUpdateRolePermissionsDataAttributesPermissionsInnerWithDefaults() *UpdateRolePermissionsDataAttributesPermissionsInner`

NewUpdateRolePermissionsDataAttributesPermissionsInnerWithDefaults instantiates a new UpdateRolePermissionsDataAttributesPermissionsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UpdateRolePermissionsDataAttributesPermissionsInner) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UpdateRolePermissionsDataAttributesPermissionsInner) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UpdateRolePermissionsDataAttributesPermissionsInner) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetStatus

`func (o *UpdateRolePermissionsDataAttributesPermissionsInner) GetStatus() bool`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *UpdateRolePermissionsDataAttributesPermissionsInner) GetStatusOk() (*bool, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *UpdateRolePermissionsDataAttributesPermissionsInner) SetStatus(v bool)`

SetStatus sets Status field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


