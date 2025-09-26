# Transaction API Guide

This document shows how to use the transaction endpoints, including creating transactions and using the advanced filtering capabilities using the morkid/paginate package.

## Endpoints

### 1. Create Transaction
```
POST /v1/transaction
```

Creates a new transaction record.

**Request Body:**
```json
{
  "userId": "user123",
  "amount": 100.50,
  "type": "credit",
  "status": "pending",
  "description": "Payment for order #123"
}
```

**Response:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "transactionId": "Tr123xyz",
    "userId": "user123",
    "amount": 100.50,
    "type": "credit",
    "status": "pending",
    "description": "Payment for order #123",
    "createdAt": "2024-01-15T10:30:00Z"
  }
}
```

**Example cURL:**
```bash
curl -X POST http://localhost:8000/v1/transaction \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "user123",
    "amount": 100.50,
    "type": "credit",
    "status": "pending",
    "description": "Payment for order #123"
  }'
```

### 2. Get Single Transaction
```
GET /v1/transaction/{id}
```

### 3. Get Transactions List with Filtering
```
GET /v1/transactions
```

## Advanced Filtering Usage

The transactions list endpoint uses the `filters` parameter with JSON array format for advanced filtering:

```
GET /v1/transactions?filters=[["column","operator","value"]]
```

## Filter Examples

### 1. Simple Equality Filter
Filter transactions by type:
```
GET /v1/transactions?filters=[["type","=","credit"]]
```

### 2. LIKE Filter for Text Search
Search in transaction descriptions:
```
GET /v1/transactions?filters=[["description","like","payment"]]
```

### 3. Multiple Filters with OR
Get credit transactions OR completed transactions:
```
GET /v1/transactions?filters=[["type","=","credit"],["or"],["status","=","completed"]]
```

### 4. Multiple Filters with AND
Get credit transactions that are completed:
```
GET /v1/transactions?filters=[["type","=","credit"],["and"],["status","=","completed"]]
```

### 5. Date Range Filter
Get transactions from a specific date range:
```
GET /v1/transactions?filters=[["created_at","between",["2024-01-01","2024-12-31"]]]
```

### 6. IN Filter
Get transactions with specific statuses:
```
GET /v1/transactions?filters=[["status","in",["completed","pending","processing"]]]
```

### 7. NULL Checks
Get transactions without description:
```
GET /v1/transactions?filters=[["description","is",null]]
```

### 8. Complex Nested Filters
Complex condition: (type=credit AND status=completed) OR (amount > 100):
```
GET /v1/transactions?filters=[[["type","=","credit"],["and"],["status","=","completed"]],["or"],["amount",">",100]]
```

### 9. User Relationship Filtering
Filter by user name (using JOIN relationship):
```
GET /v1/transactions?filters=[["user.name","like","john"]]
```

## Supported Operators

- `=` - Equal
- `!=`, `<>` - Not equal
- `>` - Greater than
- `>=` - Greater than or equal
- `<` - Less than
- `<=` - Less than or equal
- `like` - Pattern matching with wildcards
- `not like` - Negated pattern matching
- `in` - Value in list
- `not in` - Value not in list
- `between` - Value between two values
- `is` - IS NULL check
- `is not` - IS NOT NULL check

## Logical Operators

- `["and"]` - AND condition
- `["or"]` - OR condition (default if not specified)

## Pagination and Sorting

Combine filters with pagination and sorting:
```
GET /v1/transactions?page=1&size=20&sort=-created_at&filters=[["type","=","credit"]]
```

## URL Encoding

When using filters in URLs, make sure to properly encode the JSON:
```javascript
const filters = [["type","=","credit"],["or"],["status","=","completed"]];
const encodedFilters = encodeURIComponent(JSON.stringify(filters));
const url = `/v1/transactions?filters=${encodedFilters}`;
```

## Migration from Old Parameters

### Before (custom parameters):
```
GET /transactions?type=credit&status=completed&search=payment
```

### After (morkid/paginate filters):
```
GET /v1/transactions?filters=[["type","=","credit"],["and"],["status","=","completed"],["and"],["description","like","payment"]]
```

This new approach provides much more flexibility and supports complex filtering scenarios that weren't possible with the old custom parameters.