# UpdateOrganizationLinksData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | Upload session id | 
**Type** | **string** |  | 
**Attributes** | [**UpdateOrganizationLinksDataAttributes**](UpdateOrganizationLinksDataAttributes.md) |  | 
**Relationships** | [**UpdateOrganizationLinksDataRelationships**](UpdateOrganizationLinksDataRelationships.md) |  | 

## Methods

### NewUpdateOrganizationLinksData

`func NewUpdateOrganizationLinksData(id uuid.UUID, type_ string, attributes UpdateOrganizationLinksDataAttributes, relationships UpdateOrganizationLinksDataRelationships, ) *UpdateOrganizationLinksData`

NewUpdateOrganizationLinksData instantiates a new UpdateOrganizationLinksData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateOrganizationLinksDataWithDefaults

`func NewUpdateOrganizationLinksDataWithDefaults() *UpdateOrganizationLinksData`

NewUpdateOrganizationLinksDataWithDefaults instantiates a new UpdateOrganizationLinksData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UpdateOrganizationLinksData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UpdateOrganizationLinksData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UpdateOrganizationLinksData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *UpdateOrganizationLinksData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UpdateOrganizationLinksData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UpdateOrganizationLinksData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *UpdateOrganizationLinksData) GetAttributes() UpdateOrganizationLinksDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *UpdateOrganizationLinksData) GetAttributesOk() (*UpdateOrganizationLinksDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *UpdateOrganizationLinksData) SetAttributes(v UpdateOrganizationLinksDataAttributes)`

SetAttributes sets Attributes field to given value.


### GetRelationships

`func (o *UpdateOrganizationLinksData) GetRelationships() UpdateOrganizationLinksDataRelationships`

GetRelationships returns the Relationships field if non-nil, zero value otherwise.

### GetRelationshipsOk

`func (o *UpdateOrganizationLinksData) GetRelationshipsOk() (*UpdateOrganizationLinksDataRelationships, bool)`

GetRelationshipsOk returns a tuple with the Relationships field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRelationships

`func (o *UpdateOrganizationLinksData) SetRelationships(v UpdateOrganizationLinksDataRelationships)`

SetRelationships sets Relationships field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


