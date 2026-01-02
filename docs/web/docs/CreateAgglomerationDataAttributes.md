# CreateOrganizationDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Head** | [**uuid.UUID**](uuid.UUID.md) | Account ID of the head of the organization | 
**Name** | **string** | The name of the organization | 
**Icon** | Pointer to **string** | The icon representing the organization | [optional] 

## Methods

### NewCreateOrganizationDataAttributes

`func NewCreateOrganizationDataAttributes(head uuid.UUID, name string, ) *CreateOrganizationDataAttributes`

NewCreateOrganizationDataAttributes instantiates a new CreateOrganizationDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateOrganizationDataAttributesWithDefaults

`func NewCreateOrganizationDataAttributesWithDefaults() *CreateOrganizationDataAttributes`

NewCreateOrganizationDataAttributesWithDefaults instantiates a new CreateOrganizationDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHead

`func (o *CreateOrganizationDataAttributes) GetHead() uuid.UUID`

GetHead returns the Head field if non-nil, zero value otherwise.

### GetHeadOk

`func (o *CreateOrganizationDataAttributes) GetHeadOk() (*uuid.UUID, bool)`

GetHeadOk returns a tuple with the Head field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHead

`func (o *CreateOrganizationDataAttributes) SetHead(v uuid.UUID)`

SetHead sets Head field to given value.


### GetName

`func (o *CreateOrganizationDataAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateOrganizationDataAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateOrganizationDataAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetIcon

`func (o *CreateOrganizationDataAttributes) GetIcon() string`

GetIcon returns the Icon field if non-nil, zero value otherwise.

### GetIconOk

`func (o *CreateOrganizationDataAttributes) GetIconOk() (*string, bool)`

GetIconOk returns a tuple with the Icon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIcon

`func (o *CreateOrganizationDataAttributes) SetIcon(v string)`

SetIcon sets Icon field to given value.

### HasIcon

`func (o *CreateOrganizationDataAttributes) HasIcon() bool`

HasIcon returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


