# Custom Filters API Documentation

## Overview

The Custom Filters API provides a type-aware, dynamic query builder that allows the frontend to construct complex filter queries. This system supports various data types, operators, and compound conditions using AND/OR logic.

## Filter Types

The following data types are supported:

- `string` - Text fields
- `number` - Numeric fields (integers, floats)
- `date` - Date fields (YYYY-MM-DD)
- `datetime` - Date and time fields
- `boolean` - True/false fields
- `enum` - Predefined list of values
- `array` - Array/list fields

## Operators

Each data type supports specific operators:

### String Operators
- `equals` - Exact match
- `not_equals` - Not equal to
- `contains` - Contains substring
- `not_contains` - Does not contain substring
- `starts_with` - Starts with string
- `ends_with` - Ends with string
- `is_empty` - Field is empty
- `is_not_empty` - Field is not empty
- `regex_match` - Matches regex pattern
- `is_null` - Field is null
- `is_not_null` - Field is not null

### Number Operators
- `equals` - Equal to
- `not_equals` - Not equal to
- `greater_than` - Greater than
- `less_than` - Less than
- `greater_than_or_equal` - Greater than or equal to
- `less_than_or_equal` - Less than or equal to
- `between` - Between two values (inclusive)
- `not_between` - Not between two values
- `is_null` - Field is null
- `is_not_null` - Field is not null

### Date/DateTime Operators
- `equals` - Exact date match
- `not_equals` - Not equal to date
- `before` - Before date
- `after` - After date
- `between` - Between two dates
- `not_between` - Not between two dates
- `is_today` - Is today's date
- `is_yesterday` - Is yesterday's date
- `is_this_week` - Is in current week
- `is_this_month` - Is in current month
- `is_this_year` - Is in current year
- `last_n_days` - In the last N days
- `next_n_days` - In the next N days
- `is_null` - Field is null
- `is_not_null` - Field is not null

### Boolean Operators
- `is_true` - Value is true
- `is_false` - Value is false
- `is_null` - Field is null
- `is_not_null` - Field is not null

### Enum Operators
- `equals` - Equal to value
- `not_equals` - Not equal to value
- `in` - In list of values
- `not_in` - Not in list of values
- `is_null` - Field is null
- `is_not_null` - Field is not null

### Array Operators
- `contains` - Contains element
- `not_contains` - Does not contain element
- `contains_any` - Contains any of the elements
- `contains_all` - Contains all elements
- `is_empty` - Array is empty
- `is_not_empty` - Array is not empty

## API Endpoints

### 1. Get Filter Metadata
Retrieve available filters, operators, and configuration for a resource.

**Endpoint:** `GET /api/{resource}/filters`

**Example:** `GET /api/books/filters`

**Response:**
```json
{
  "success": true,
  "data": {
    "filters": [
      {
        "field": "price",
        "label": "Price",
        "type": "number",
        "operators": [
          "equals",
          "greater_than",
          "less_than",
          "greater_than_or_equal",
          "less_than_or_equal",
          "between"
        ],
        "validation": {
          "min": 0,
          "max": 999.99
        }
      },
      {
        "field": "status",
        "label": "Status",
        "type": "enum",
        "operators": ["equals", "not_equals", "in", "not_in"],
        "enum_values": ["AVAILABLE", "BORROWED", "MAINTENANCE", "RESERVED", "LOST"]
      },
      {
        "field": "published_at",
        "label": "Published Date",
        "type": "date",
        "operators": [
          "before",
          "after",
          "between",
          "is_today",
          "is_this_month",
          "is_this_year",
          "last_n_days"
        ],
        "format": "2006-01-02"
      }
    ],
    "logic_operators": ["AND", "OR"],
    "resource": "book",
    "searchable_fields": ["title", "author", "isbn", "description"],
    "sortable_fields": ["id", "title", "author", "price", "created_at"],
    "filterable_fields": ["status", "author", "price", "published_at"]
  }
}
```

### 2. Apply Filters to List
Apply custom filters when retrieving a list of resources.

**Endpoint:** `GET /api/{resource}?filters={filterJSON}`

**Filter JSON Structure:**

#### Simple Filter
```json
{
  "field": "price",
  "operator": "greater_than",
  "value": 30
}
```

#### Compound Filter (AND/OR)
```json
{
  "logic": "AND",
  "conditions": [
    {
      "field": "price",
      "operator": "between",
      "value": [20, 50]
    },
    {
      "field": "status",
      "operator": "equals",
      "value": "AVAILABLE"
    }
  ]
}
```

#### Nested Compound Filter
```json
{
  "logic": "OR",
  "conditions": [
    {
      "logic": "AND",
      "conditions": [
        {
          "field": "price",
          "operator": "greater_than",
          "value": 30
        },
        {
          "field": "status",
          "operator": "equals",
          "value": "AVAILABLE"
        }
      ]
    },
    {
      "field": "author",
      "operator": "equals",
      "value": "Special Author"
    }
  ]
}
```

## Frontend Integration Examples

### JavaScript/TypeScript

```typescript
// Define filter types
interface FilterCondition {
  field: string;
  operator: string;
  value: any;
}

interface CompoundFilter {
  logic: 'AND' | 'OR';
  conditions: (FilterCondition | CompoundFilter)[];
}

// Simple filter example
const priceFilter: FilterCondition = {
  field: 'price',
  operator: 'greater_than',
  value: 30
};

// Compound filter example
const compoundFilter: CompoundFilter = {
  logic: 'AND',
  conditions: [
    {
      field: 'price',
      operator: 'between',
      value: [20, 50]
    },
    {
      field: 'status',
      operator: 'equals',
      value: 'AVAILABLE'
    }
  ]
};

// Apply filters
const filterJSON = JSON.stringify(compoundFilter);
const url = `/api/books?filters=${encodeURIComponent(filterJSON)}`;

fetch(url, {
  headers: {
    'Accept': 'application/json',
    'Authorization': 'Bearer ' + token
  }
})
.then(response => response.json())
.then(data => {
  console.log('Filtered results:', data);
});
```

### React Component Example

```tsx
import React, { useState, useEffect } from 'react';
import { FilterBuilder } from './FilterBuilder'; // Custom component

const BookList = () => {
  const [filters, setFilters] = useState({});
  const [filterMetadata, setFilterMetadata] = useState(null);
  const [books, setBooks] = useState([]);

  // Fetch filter metadata on mount
  useEffect(() => {
    fetch('/api/books/filters')
      .then(res => res.json())
      .then(data => setFilterMetadata(data.data));
  }, []);

  // Apply filters
  const applyFilters = (filterConditions) => {
    const filterJSON = JSON.stringify(filterConditions);
    const params = new URLSearchParams({
      filters: filterJSON,
      page: 1,
      pageSize: 20
    });

    fetch(`/api/books?${params}`)
      .then(res => res.json())
      .then(data => setBooks(data.data.data));
  };

  return (
    <div>
      {filterMetadata && (
        <FilterBuilder
          metadata={filterMetadata}
          onApply={applyFilters}
        />
      )}
      <BookGrid books={books} />
    </div>
  );
};
```

## Value Formats

### Date Values
- Format: `YYYY-MM-DD` (e.g., `2024-01-15`)
- Example: `{ "field": "published_at", "operator": "after", "value": "2024-01-01" }`

### DateTime Values
- Format: RFC3339 (e.g., `2024-01-15T14:30:00Z`)
- Example: `{ "field": "created_at", "operator": "after", "value": "2024-01-01T00:00:00Z" }`

### Between Operator Values
- Format: Array of two values `[min, max]`
- Number example: `{ "field": "price", "operator": "between", "value": [10, 50] }`
- Date example: `{ "field": "published_at", "operator": "between", "value": ["2023-01-01", "2023-12-31"] }`

### In/Not In Operator Values
- Format: Array of values
- Example: `{ "field": "status", "operator": "in", "value": ["AVAILABLE", "BORROWED"] }`

### Last/Next N Days Values
- Format: Number of days
- Example: `{ "field": "created_at", "operator": "last_n_days", "value": 7 }`

## Error Handling

### Invalid Filter Structure
**Status:** 400 Bad Request
```json
{
  "success": false,
  "message": "Invalid filter JSON: field is required"
}
```

### Invalid Operator for Type
**Status:** 400 Bad Request
```json
{
  "success": false,
  "message": "Operator 'contains' is not valid for field 'price' of type 'number'"
}
```

### Invalid Value Format
**Status:** 400 Bad Request
```json
{
  "success": false,
  "message": "Invalid date format: 2024-13-45"
}
```

## Best Practices

1. **Cache Filter Metadata**: Filter metadata doesn't change often, so cache it on the frontend.

2. **Validate Before Sending**: Use the metadata to validate filters on the frontend before sending to the API.

3. **Use Compound Filters Wisely**: Deeply nested compound filters can impact performance. Keep nesting to 2-3 levels maximum.

4. **Handle Empty Results**: Always handle cases where filters return no results.

5. **Provide Clear UI**: Use the filter metadata to build intuitive UI components that guide users to valid filter combinations.

6. **URL Encoding**: Always URL-encode the filter JSON when sending as a query parameter.

7. **Batch Filters**: When possible, combine multiple conditions into a single compound filter rather than making multiple API calls.

## Performance Considerations

- Filters are applied at the database level for optimal performance
- Indexed fields provide better filter performance
- Complex compound filters may impact query performance
- The `between` operator is generally more efficient than multiple comparison operators
- String operations like `contains` may be slower on large datasets

## Security

- All filter values are parameterized to prevent SQL injection
- Field names are validated against allowed filter fields
- Operators are validated for each field type
- Values are type-checked and sanitized before query execution