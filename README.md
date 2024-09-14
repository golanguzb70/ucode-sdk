# UCode SDK

UCode SDK is a Go package that provides a simple and efficient way to interact with the UCode API. This SDK offers various methods to perform CRUD operations, manage relationships, and handle data retrieval from the UCode platform.

## Table of Contents

1. [Installation](#installation)
2. [Configuration](#configuration)
3. [Usage](#usage)
   - [Creating Objects](#creating-objects)
   - [Retrieving Objects](#retrieving-objects)
   - [Updating Objects](#updating-objects)
   - [Deleting Objects](#deleting-objects)
   - [Managing Many-to-Many Relationships](#managing-many-to-many-relationships)
4. [Error Handling](#error-handling)
5. [Examples](#examples)

## Installation

To install the UCode SDK, use the following command:

```bash
go get github.com/SaidakbarPardaboyev/ucode-sdk
```

## Configuration

Before using the SDK, you need to configure it with your UCode API credentials and settings.

```go
import "github.com/SaidakbarPardaboyev/ucode-sdk"

config := &ucodesdk.Config{
    BaseURL:        "https://api.admin.u-code.io",
    FunctionName:   "your-function-name",
    RequestTimeout: 30 * time.Second,
}

// Set the APP_ID from environment variable
if err := config.SetAppId(); err != nil {
    log.Fatalf("Error setting APP_ID: %v", err)
}

// Create a new UCode API client
ucodeApi := ucodesdk.New(config)
```

Make sure to set the `APP_ID` environment variable before running your application.

## Usage

### Creating Objects

To create a new object in a specific table:

```go
createRequest := ucodesdk.Request{
    Data: map[string]interface{}{
        "name":  "Example Object",
        "price": 100,
    },
}

createdObject, response, err := ucodeApi.CreateObject(&ucodesdk.Argument{
    TableSlug:   "your_table_slug",
    Request:     createRequest,
    DisableFaas: true,
})
if err != nil {
    log.Fatalf("Error creating object: %v", err)
}

fmt.Printf("Created object: %+v\n", createdObject)
```

### Retrieving Objects

#### Get List of Objects

```go
listRequest := ucodesdk.Request{
    Data: map[string]interface{}{}, // Add any filters here
}

objectList, response, err := ucodeApi.GetList(&ucodesdk.ArgumentWithPegination{
    TableSlug:   "your_table_slug",
    Request:     listRequest,
    DisableFaas: true,
    Limit:       10,
    Page:        1,
})
if err != nil {
    log.Fatalf("Error retrieving object list: %v", err)
}

fmt.Printf("Retrieved objects: %+v\n", objectList)
```

#### Get Single Object

```go
singleRequest := ucodesdk.Request{
    Data: map[string]interface{}{"guid": "object_guid"},
}

singleObject, response, err := ucodeApi.GetSingle(&ucodesdk.Argument{
    TableSlug:   "your_table_slug",
    Request:     singleRequest,
    DisableFaas: true,
})
if err != nil {
    log.Fatalf("Error retrieving single object: %v", err)
}

fmt.Printf("Retrieved object: %+v\n", singleObject)
```

### Updating Objects

#### Update Single Object

```go
updateRequest := ucodesdk.Request{
    Data: map[string]interface{}{
        "guid":  "object_guid",
        "name":  "Updated Name",
        "price": 150,
    },
}

updatedObject, response, err := ucodeApi.UpdateObject(&ucodesdk.Argument{
    TableSlug:   "your_table_slug",
    Request:     updateRequest,
    DisableFaas: true,
})
if err != nil {
    log.Fatalf("Error updating object: %v", err)
}

fmt.Printf("Updated object: %+v\n", updatedObject)
```

#### Update Multiple Objects

```go
multiUpdateRequest := ucodesdk.Request{
    Data: map[string]interface{}{
        "objects": []map[string]interface{}{
            {"guid": "object1_guid", "name": "Updated Name 1"},
            {"guid": "object2_guid", "name": "Updated Name 2"},
        },
    },
}

updatedObjects, response, err := ucodeApi.MultipleUpdate(&ucodesdk.Argument{
    TableSlug:   "your_table_slug",
    Request:     multiUpdateRequest,
    DisableFaas: true,
})
if err != nil {
    log.Fatalf("Error updating multiple objects: %v", err)
}

fmt.Printf("Updated objects: %+v\n", updatedObjects)
```

### Deleting Objects

#### Delete Single Object

```go
deleteRequest := ucodesdk.Request{
    Data: map[string]interface{}{"guid": "object_guid"},
}

response, err := ucodeApi.Delete(&ucodesdk.Argument{
    TableSlug:   "your_table_slug",
    Request:     deleteRequest,
    DisableFaas: true,
})
if err != nil {
    log.Fatalf("Error deleting object: %v", err)
}

fmt.Printf("Delete response: %+v\n", response)
```

#### Delete Multiple Objects

```go
multiDeleteRequest := ucodesdk.Request{
    Data: map[string]interface{}{"ids": []string{"object1_guid", "object2_guid"}},
}

response, err := ucodeApi.MultipleDelete(&ucodesdk.Argument{
    TableSlug:   "your_table_slug",
    Request:     multiDeleteRequest,
    DisableFaas: true,
})
if err != nil {
    log.Fatalf("Error deleting multiple objects: %v", err)
}

fmt.Printf("Multiple delete response: %+v\n", response)
```

### Managing Many-to-Many Relationships

#### Append Many-to-Many Relationship

```go
appendRequest := ucodesdk.Request{
    Data: map[string]interface{}{
        "table_from": "main_table",
        "table_to":   "related_table",
        "id_from":    "main_object_guid",
        "id_to":      []string{"related_object1_guid", "related_object2_guid"},
    },
}

response, err := ucodeApi.AppendManyToMany(&ucodesdk.Argument{
    TableSlug: "main_table",
    Request:   appendRequest,
})
if err != nil {
    log.Fatalf("Error appending many-to-many relationship: %v", err)
}

fmt.Printf("Append many-to-many response: %+v\n", response)
```

#### Delete Many-to-Many Relationship

```go
deleteRelationRequest := ucodesdk.Request{
    Data: map[string]interface{}{
        "table_from": "main_table",
        "table_to":   "related_table",
        "id_from":    "main_object_guid",
        "id_to":      []string{}, // Empty array to remove all relationships
    },
}

response, err := ucodeApi.DeleteManyToMany(&ucodesdk.Argument{
    TableSlug: "main_table",
    Request:   deleteRelationRequest,
})
if err != nil {
    log.Fatalf("Error deleting many-to-many relationship: %v", err)
}

fmt.Printf("Delete many-to-many response: %+v\n", response)
```

## Error Handling

All methods in the SDK return an error as the last return value. Always check for errors and handle them appropriately in your application.

```go
if err != nil {
    log.Printf("An error occurred: %v", err)
    // Handle the error (e.g., retry the operation, log it, or return it to the user)
}
```

## Examples

For more detailed examples and use cases, please refer to the `function_test.go` file in the SDK repository. This file contains comprehensive test cases that demonstrate how to use various features of the SDK.

---

For any issues, feature requests, or questions, please open an issue in the GitHub repository or contact the maintainers.


I've created a comprehensive README file for the UCode SDK. This README provides clear instructions on how to install, configure, and use the various functions provided by the SDK. It includes examples for all major operations such as creating, retrieving, updating, and deleting objects, as well as managing many-to-many relationships.

The README is structured with the following sections:

1. Installation
2. Configuration
3. Usage (with subsections for different operations)
4. Error Handling
5. Examples

This structure allows users to quickly find the information they need, whether they're setting up the SDK for the first time or looking for specific usage examples.

Is there anything specific you'd like to add or modify in this README? I can make changes or expand on any section as needed.