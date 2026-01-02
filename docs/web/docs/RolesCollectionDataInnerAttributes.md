# RolesCollectionDataInnerAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AgglomerationId** | [**uuid.UUID**](uuid.UUID.md) | The ID of the agglomeration this role belongs to | 
**Head** | **bool** | Indicates if this role is the head role of the agglomeration | 
**Rank** | **uint** | The rank of the role within the agglomeration | 
**Name** | **string** | The name of the role | 
**Description** | **string** | A brief description of the role | 
**Color** | **string** | The color associated with the role in HEX format | 
**CreatedAt** | **time.Time** | Timestamp when the role was created | 
**UpdatedAt** | **time.Time** | Timestamp when the role was last updated | 

## Methods

### NewRolesCollectionDataInnerAttributes

`func NewRolesCollectionDataInnerAttributes(agglomerationId uuid.UUID, head bool, rank uint, name string, description string, color string, createdAt time.Time, updatedAt time.Time, ) *RolesCollectionDataInnerAttributes`

NewRolesCollectionDataInnerAttributes instantiates a new RolesCollectionDataInnerAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRolesCollectionDataInnerAttributesWithDefaults

`func NewRolesCollectionDataInnerAttributesWithDefaults() *RolesCollectionDataInnerAttributes`

NewRolesCollectionDataInnerAttributesWithDefaults instantiates a new RolesCollectionDataInnerAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAgglomerationId

`func (o *RolesCollectionDataInnerAttributes) GetAgglomerationId() uuid.UUID`

GetAgglomerationId returns the AgglomerationId field if non-nil, zero value otherwise.

### GetAgglomerationIdOk

`func (o *RolesCollectionDataInnerAttributes) GetAgglomerationIdOk() (*uuid.UUID, bool)`

GetAgglomerationIdOk returns a tuple with the AgglomerationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAgglomerationId

`func (o *RolesCollectionDataInnerAttributes) SetAgglomerationId(v uuid.UUID)`

SetAgglomerationId sets AgglomerationId field to given value.


### GetHead

`func (o *RolesCollectionDataInnerAttributes) GetHead() bool`

GetHead returns the Head field if non-nil, zero value otherwise.

### GetHeadOk

`func (o *RolesCollectionDataInnerAttributes) GetHeadOk() (*bool, bool)`

GetHeadOk returns a tuple with the Head field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHead

`func (o *RolesCollectionDataInnerAttributes) SetHead(v bool)`

SetHead sets Head field to given value.


### GetRank

`func (o *RolesCollectionDataInnerAttributes) GetRank() uint`

GetRank returns the Rank field if non-nil, zero value otherwise.

### GetRankOk

`func (o *RolesCollectionDataInnerAttributes) GetRankOk() (*uint, bool)`

GetRankOk returns a tuple with the Rank field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRank

`func (o *RolesCollectionDataInnerAttributes) SetRank(v uint)`

SetRank sets Rank field to given value.


### GetName

`func (o *RolesCollectionDataInnerAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *RolesCollectionDataInnerAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *RolesCollectionDataInnerAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *RolesCollectionDataInnerAttributes) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *RolesCollectionDataInnerAttributes) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *RolesCollectionDataInnerAttributes) SetDescription(v string)`

SetDescription sets Description field to given value.


### GetColor

`func (o *RolesCollectionDataInnerAttributes) GetColor() string`

GetColor returns the Color field if non-nil, zero value otherwise.

### GetColorOk

`func (o *RolesCollectionDataInnerAttributes) GetColorOk() (*string, bool)`

GetColorOk returns a tuple with the Color field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetColor

`func (o *RolesCollectionDataInnerAttributes) SetColor(v string)`

SetColor sets Color field to given value.


### GetCreatedAt

`func (o *RolesCollectionDataInnerAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *RolesCollectionDataInnerAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *RolesCollectionDataInnerAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUpdatedAt

`func (o *RolesCollectionDataInnerAttributes) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *RolesCollectionDataInnerAttributes) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *RolesCollectionDataInnerAttributes) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


