# AgglomerationDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | **string** | The status of the agglomeration | 
**Name** | **string** | The name of the agglomeration | 
**Icon** | **string** | The icon representing the agglomeration | 
**CreatedAt** | **time.Time** | The date and time when the agglomeration was created | 
**UpdatedAt** | **time.Time** | The date and time when the agglomeration was last updated | 

## Methods

### NewAgglomerationDataAttributes

`func NewAgglomerationDataAttributes(status string, name string, icon string, createdAt time.Time, updatedAt time.Time, ) *AgglomerationDataAttributes`

NewAgglomerationDataAttributes instantiates a new AgglomerationDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAgglomerationDataAttributesWithDefaults

`func NewAgglomerationDataAttributesWithDefaults() *AgglomerationDataAttributes`

NewAgglomerationDataAttributesWithDefaults instantiates a new AgglomerationDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *AgglomerationDataAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *AgglomerationDataAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *AgglomerationDataAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetName

`func (o *AgglomerationDataAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *AgglomerationDataAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *AgglomerationDataAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetIcon

`func (o *AgglomerationDataAttributes) GetIcon() string`

GetIcon returns the Icon field if non-nil, zero value otherwise.

### GetIconOk

`func (o *AgglomerationDataAttributes) GetIconOk() (*string, bool)`

GetIconOk returns a tuple with the Icon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIcon

`func (o *AgglomerationDataAttributes) SetIcon(v string)`

SetIcon sets Icon field to given value.


### GetCreatedAt

`func (o *AgglomerationDataAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *AgglomerationDataAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *AgglomerationDataAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUpdatedAt

`func (o *AgglomerationDataAttributes) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *AgglomerationDataAttributes) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *AgglomerationDataAttributes) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


