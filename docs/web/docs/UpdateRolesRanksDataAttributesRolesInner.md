# UpdateRolesRanksDataAttributesRolesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | role ID | 
**Rank** | **uint** | The new rank of the role in the hierarchy | 

## Methods

### NewUpdateRolesRanksDataAttributesRolesInner

`func NewUpdateRolesRanksDataAttributesRolesInner(id uuid.UUID, rank uint, ) *UpdateRolesRanksDataAttributesRolesInner`

NewUpdateRolesRanksDataAttributesRolesInner instantiates a new UpdateRolesRanksDataAttributesRolesInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateRolesRanksDataAttributesRolesInnerWithDefaults

`func NewUpdateRolesRanksDataAttributesRolesInnerWithDefaults() *UpdateRolesRanksDataAttributesRolesInner`

NewUpdateRolesRanksDataAttributesRolesInnerWithDefaults instantiates a new UpdateRolesRanksDataAttributesRolesInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UpdateRolesRanksDataAttributesRolesInner) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UpdateRolesRanksDataAttributesRolesInner) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UpdateRolesRanksDataAttributesRolesInner) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetRank

`func (o *UpdateRolesRanksDataAttributesRolesInner) GetRank() uint`

GetRank returns the Rank field if non-nil, zero value otherwise.

### GetRankOk

`func (o *UpdateRolesRanksDataAttributesRolesInner) GetRankOk() (*uint, bool)`

GetRankOk returns a tuple with the Rank field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRank

`func (o *UpdateRolesRanksDataAttributesRolesInner) SetRank(v uint)`

SetRank sets Rank field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


