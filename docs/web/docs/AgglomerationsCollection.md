# AgglomerationsCollection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]AgglomerationData**](AgglomerationData.md) |  | 
**Links** | [**PaginationData**](PaginationData.md) |  | 

## Methods

### NewAgglomerationsCollection

`func NewAgglomerationsCollection(data []AgglomerationData, links PaginationData, ) *AgglomerationsCollection`

NewAgglomerationsCollection instantiates a new AgglomerationsCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAgglomerationsCollectionWithDefaults

`func NewAgglomerationsCollectionWithDefaults() *AgglomerationsCollection`

NewAgglomerationsCollectionWithDefaults instantiates a new AgglomerationsCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *AgglomerationsCollection) GetData() []AgglomerationData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *AgglomerationsCollection) GetDataOk() (*[]AgglomerationData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *AgglomerationsCollection) SetData(v []AgglomerationData)`

SetData sets Data field to given value.


### GetLinks

`func (o *AgglomerationsCollection) GetLinks() PaginationData`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *AgglomerationsCollection) GetLinksOk() (*PaginationData, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *AgglomerationsCollection) SetLinks(v PaginationData)`

SetLinks sets Links field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


