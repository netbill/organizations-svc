# UpdateAgglomerationData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | agglomeration ID | 
**Type** | **string** |  | 
**Attributes** | [**UpdateAgglomerationDataAttributes**](UpdateAgglomerationDataAttributes.md) |  | 

## Methods

### NewUpdateAgglomerationData

`func NewUpdateAgglomerationData(id uuid.UUID, type_ string, attributes UpdateAgglomerationDataAttributes, ) *UpdateAgglomerationData`

NewUpdateAgglomerationData instantiates a new UpdateAgglomerationData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateAgglomerationDataWithDefaults

`func NewUpdateAgglomerationDataWithDefaults() *UpdateAgglomerationData`

NewUpdateAgglomerationDataWithDefaults instantiates a new UpdateAgglomerationData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UpdateAgglomerationData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UpdateAgglomerationData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UpdateAgglomerationData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *UpdateAgglomerationData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UpdateAgglomerationData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UpdateAgglomerationData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *UpdateAgglomerationData) GetAttributes() UpdateAgglomerationDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *UpdateAgglomerationData) GetAttributesOk() (*UpdateAgglomerationDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *UpdateAgglomerationData) SetAttributes(v UpdateAgglomerationDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


