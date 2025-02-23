# API Documentation

## REST API

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication
Currently, the API does not require authentication. For production use, implement appropriate authentication mechanisms.

### Endpoints

#### Messages API

##### Create Message
```http
POST /messages
Content-Type: application/json

{
    "content": "string"
}
```

**Response**
```json
{
    "id": "uuid",
    "content": "string",
    "created_at": "timestamp",
    "updated_at": "timestamp"
}
```

##### Get Message
```http
GET /messages/{id}
```

**Response**
```json
{
    "id": "uuid",
    "content": "string",
    "created_at": "timestamp",
    "updated_at": "timestamp"
}
```

##### Update Message
```http
PUT /messages/{id}
Content-Type: application/json

{
    "content": "string"
}
```

**Response**
```json
{
    "id": "uuid",
    "content": "string",
    "created_at": "timestamp",
    "updated_at": "timestamp"
}
```

##### Delete Message
```http
DELETE /messages/{id}
```

**Response**
```
204 No Content
```

##### List Messages
```http
GET /messages?page=1&page_size=10
```

**Response**
```json
{
    "messages": [
        {
            "id": "uuid",
            "content": "string",
            "created_at": "timestamp",
            "updated_at": "timestamp"
        }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
}
```

## gRPC Service

### Service Definition
```protobuf
service MessageService {
    rpc CreateMessage(CreateMessageRequest) returns (MessageResponse) {}
    rpc GetMessage(GetMessageRequest) returns (MessageResponse) {}
    rpc UpdateMessage(UpdateMessageRequest) returns (MessageResponse) {}
    rpc DeleteMessage(DeleteMessageRequest) returns (google.protobuf.Empty) {}
    rpc ListMessages(ListMessagesRequest) returns (ListMessagesResponse) {}
    rpc StreamMessages(google.protobuf.Empty) returns (stream MessageResponse) {}
}
```

### Message Types
```protobuf
message CreateMessageRequest {
    string content = 1;
}

message GetMessageRequest {
    string id = 1;
}

message UpdateMessageRequest {
    string id = 1;
    string content = 2;
}

message DeleteMessageRequest {
    string id = 1;
}

message ListMessagesRequest {
    int32 page = 1;
    int32 page_size = 2;
}

message ListMessagesResponse {
    repeated MessageResponse messages = 1;
    int32 total = 2;
}

message MessageResponse {
    string id = 1;
    string content = 2;
    google.protobuf.Timestamp created_at = 3;
    google.protobuf.Timestamp updated_at = 4;
}
```

## Error Handling

### HTTP Error Responses
All error responses follow this format:
```json
{
    "error": {
        "code": 400,
        "message": "Error message",
        "details": "Detailed error information"
    }
}
```

### Common HTTP Status Codes
- `200 OK`: Successful request
- `201 Created`: Resource successfully created
- `204 No Content`: Resource successfully deleted
- `400 Bad Request`: Invalid request payload
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

### Validation Errors
Validation errors include details about which fields failed validation:
```json
{
    "error": {
        "code": 400,
        "message": "Validation failed",
        "details": "content: failed validation for 'required'"
    }
}
```

## Rate Limiting
The API implements rate limiting with the following defaults:
- 100 requests per minute per IP address
- Rate limit headers are included in responses:
  - `X-RateLimit-Limit`
  - `X-RateLimit-Remaining`
  - `X-RateLimit-Reset`
