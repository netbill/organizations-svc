# ProfileAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Username** | **string** | username for this account and profile | 
**Pseudonym** | Pointer to **string** | profile pseudonym | [optional] 
**AvatarKey** | Pointer to **string** | avatar key | [optional] 
**Version** | **int32** | profile version | 
**UpdatedAt** | **time.Time** | updated at | 
**CreatedAt** | **time.Time** | created at | 

## Methods

### NewProfileAttributes

`func NewProfileAttributes(username string, version int32, updatedAt time.Time, createdAt time.Time, ) *ProfileAttributes`

NewProfileAttributes instantiates a new ProfileAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProfileAttributesWithDefaults

`func NewProfileAttributesWithDefaults() *ProfileAttributes`

NewProfileAttributesWithDefaults instantiates a new ProfileAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUsername

`func (o *ProfileAttributes) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *ProfileAttributes) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *ProfileAttributes) SetUsername(v string)`

SetUsername sets Username field to given value.


### GetPseudonym

`func (o *ProfileAttributes) GetPseudonym() string`

GetPseudonym returns the Pseudonym field if non-nil, zero value otherwise.

### GetPseudonymOk

`func (o *ProfileAttributes) GetPseudonymOk() (*string, bool)`

GetPseudonymOk returns a tuple with the Pseudonym field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPseudonym

`func (o *ProfileAttributes) SetPseudonym(v string)`

SetPseudonym sets Pseudonym field to given value.

### HasPseudonym

`func (o *ProfileAttributes) HasPseudonym() bool`

HasPseudonym returns a boolean if a field has been set.

### GetAvatarKey

`func (o *ProfileAttributes) GetAvatarKey() string`

GetAvatarKey returns the AvatarKey field if non-nil, zero value otherwise.

### GetAvatarKeyOk

`func (o *ProfileAttributes) GetAvatarKeyOk() (*string, bool)`

GetAvatarKeyOk returns a tuple with the AvatarKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvatarKey

`func (o *ProfileAttributes) SetAvatarKey(v string)`

SetAvatarKey sets AvatarKey field to given value.

### HasAvatarKey

`func (o *ProfileAttributes) HasAvatarKey() bool`

HasAvatarKey returns a boolean if a field has been set.

### GetVersion

`func (o *ProfileAttributes) GetVersion() int32`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *ProfileAttributes) GetVersionOk() (*int32, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *ProfileAttributes) SetVersion(v int32)`

SetVersion sets Version field to given value.


### GetUpdatedAt

`func (o *ProfileAttributes) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *ProfileAttributes) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *ProfileAttributes) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.


### GetCreatedAt

`func (o *ProfileAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *ProfileAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *ProfileAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


