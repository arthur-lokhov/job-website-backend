# API Documentation

This document provides a summary of the API endpoints, including request and response formats.

## Public Endpoints

### GET /api/v1/vacancies

**Description:** Retrieves a list of public vacancies.

**Authentication:** Public

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/vacancies
```

**Response:**
*   **Status:** `200 OK`
*   **Body:**
```json
[
  {
    "id": "6fb5c58f-d7c1-4c18-8ed7-81712fcfb5f2",
    "name": "Haskell Developer",
    "location": {
      "id": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
      "name": "Удаленно"
    },
    "department": {
      "id": "c9d0e1f2-a3b4-5678-9012-cdef01234567",
      "name": ""
    },
    "level": {
      "id": "44444444-4444-4444-4444-444444444444",
      "name": "Senior"
    },
    "important": true,
    "priority": 101
  }
]
```

### GET /api/v1/vacancies/{id}

**Description:** Retrieves details for a specific vacancy.

**Authentication:** Public

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/vacancies/6fb5c58f-d7c1-4c18-8ed7-81712fcfb5f2
```

**Response:**
*   **Status:** `200 OK`
*   **Body:**
```json
{
  "id": "6fb5c58f-d7c1-4c18-8ed7-81712fcfb5f2",
  "name": "Haskell Developer",
  "location": {
    "id": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    "name": "Удаленно"
  },
  "department": {
    "id": "c9d0e1f2-a3b4-5678-9012-cdef01234567",
    "name": ""
  },
  "level": {
    "id": "44444444-4444-4444-4444-444444444444",
    "name": "Senior"
  },
  "important": true,
  "priority": 101,
  "info": "Haskell разработчик."
}
```

### GET /api/v1/vacancies/{id}/form

**Description:** Retrieves the application form schema for a specific vacancy.

**Authentication:** Public

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/vacancies/6fb5c58f-d7c1-4c18-8ed7-81712fcfb5f2/form
```

**Response:**
*   **Status:** `200 OK`
*   **Body:**
```json
{
  "formSchema": {
    "formSchema": {
      "email": {
        "type": "string"
      },
      "fullName": {
        "type": "string"
      }
    }
  }
}
```

### POST /api/v1/applications

**Description:** Submits a new application for a vacancy.

**Authentication:** Public

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/applications \
  -H "Content-Type: application/json" \
  -d '{
    "vacancyId": "6fb5c58f-d7c1-4c18-8ed7-81712fcfb5f2",
    "answers": {
      "email": "test@example.com",
      "fullName": "Test User"
    }
  }'
```

**Response:**
*   **Status:** `201 Created`
*   **Body:**
```json
{
  "success": true,
  "id": "f74d0446-5662-437a-b9ab-f6138b7b1849"
}
```

### GET /api/v1/departments

**Description:** Retrieves a list of departments.

**Authentication:** Public

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/departments
```

**Response:**
*   **Status:** `200 OK`
*   **Body:**
```json
[
  {
    "label": "Тестовый отдел N2",
    "value": "job-test-2"
  },
  {
    "label": "Тестовый отдел 1",
    "value": "job-test-group-1"
  }
]
```
**Note:** The `value` field for departments is not a UUID, which is inconsistent with other parts of the API.

### GET /api/v1/levels

**Description:** Retrieves a list of job levels.

**Authentication:** Public

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/levels
```

**Response:**
*   **Status:** `200 OK`
*   **Body:**
```json
[
  {
    "label": "Intern",
    "value": "11111111-1111-1111-1111-111111111111"
  },
  {
    "label": "Junior",
    "value": "22222222-2222-2222-2222-222222222222"
  },
  {
    "label": "Middle",
    "value": "33333333-3333-3333-3333-333333333333"
  },
  {
    "label": "Senior",
    "value": "44444444-4444-4444-4444-444444444444"
  },
  {
    "label": "Lead",
    "value": "55555555-5555-5555-5555-555555555555"
  }
]
```

### GET /api/v1/locations

**Description:** Retrieves a list of locations.

**Authentication:** Public

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/locations
```

**Response:**
*   **Status:** `200 OK`
*   **Body:**
```json
[
  {
    "label": "Удаленно",
    "value": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
  },
  {
    "label": "Офис Москва",
    "value": "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
  },
  {
    "label": "Офис Санкт-Петербург",
    "value": "cccccccc-cccc-cccc-cccc-cccccccccccc"
  },
  {
    "label": "Гибрид",
    "value": "dddddddd-dddd-dddd-dddd-dddddddddddd"
  }
]
```

## Admin Endpoints

### GET /api/v1/admin/vacancies

**Description:** Retrieves a list of all vacancies for admin users.

**Authentication:** Bearer Token (`jobservice.vacancies:read` permission)

**Request:**
```bash
curl -X GET -H "Authorization: Bearer <YOUR_TOKEN>" http://localhost:8080/api/v1/admin/vacancies
```

**Response:**
*   **Status:** `200 OK`
*   **Body:**
```json
[
  {
    "id": "255057b8-1e6c-42c4-a54b-54701e06a6a8",
    "name": "Data Scientist (Python)",
    "location": {
      "id": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
      "name": "Удаленно"
    },
    "department": {
      "id": "f6a7b8c9-d0e1-2345-6789-0abcdef01234",
      "name": ""
    },
    "level": {
      "id": "44444444-4444-4444-4444-444444444444",
      "name": "Senior"
    },
    "important": false,
    "priority": 6,
    "info": "Исследователь данных для построения моделей машинного обучения.",
    "applicationForm": {
      "formSchema": null
    },
    "isActive": false,
    "createdAt": "2025-08-27T16:33:54Z",
    "updatedAt": "2025-08-27T16:33:54Z"
  }
]
```

### POST /api/v1/admin/vacancies

**Description:** Creates a new vacancy.

**Authentication:** Bearer Token (`jobservice.vacancies:create` permission)

**Request:**
```bash
curl -X POST -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Vacancy",
    "locationId": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    "departmentId": "f6a7b8c9-d0e1-2345-6789-0abcdef01234",
    "levelId": "11111111-1111-1111-1111-111111111111",
    "info": "This is a test vacancy.",
    "applicationForm": {
      "formSchema": {
        "test_field": {
          "type": "string"
        }
      }
    },
    "important": false,
    "priority": 1,
    "isActive": true
  }' http://localhost:8080/api/v1/admin/vacancies
```

**Response:**
*   **Status:** `201 Created`
*   **Body:**
```json
{
  "success": true,
  "id": "2fddc4c7-bc0b-4e36-a9d8-251e85144c73",
  "createdAt": ""
}
```

### PATCH /api/v1/admin/vacancies/{id}

**Description:** Updates an existing vacancy.

**Authentication:** Bearer Token (`jobservice.vacancies:update` permission)

**Request:**
```bash
curl -X PATCH -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Vacancy Updated",
    "important": true
  }' http://localhost:8080/api/v1/admin/vacancies/2fddc4c7-bc0b-4e36-a9d8-251e85144c73
```

**Response:**
*   **Status:** `200 OK`
*   **Body:**
```json
{
  "success": true,
  "id": "2fddc4c7-bc0b-4e36-a9d8-251e85144c73"
}
```

### DELETE /api/v1/admin/vacancies/{id}

**Description:** Deletes a vacancy.

**Authentication:** Bearer Token (`jobservice.vacancies:delete` permission)

**Request:**
```bash
curl -X DELETE -H "Authorization: Bearer <YOUR_TOKEN>" http://localhost:8080/api/v1/admin/vacancies/2fddc4c7-bc0b-4e36-a9d8-251e85144c73 -i
```

**Response:**
*   **Status:** `204 No Content`

### GET /api/v1/admin/applications

**Description:** Retrieves a list of all applications for admin users.

**Authentication:** Bearer Token (`jobservice.applications:read` permission)

**Request:**
```bash
curl -X GET -H "Authorization: Bearer <YOUR_TOKEN>" http://localhost:8080/api/v1/admin/applications
```

**Response:**
*   **Status:** `200 OK`
*   **Body:**
```json
[
  {
    "id": "3b08c602-e636-4c90-8249-120949fd7e4d",
    "vacancyId": "6fb5c58f-d7c1-4c18-8ed7-81712fcfb5f2",
    "vacancyName": "",
    "answers": {
      "email": "test@example.com",
      "fullName": "Test User"
    },
    "status": "PENDING",
    "createdAt": "2025-08-28T09:46:01Z",
    "updatedAt": "2025-08-28T09:46:01Z"
  }
]
```

### PATCH /api/v1/admin/applications/{id}

**Description:** Updates the status of an application.

**Authentication:** Bearer Token (`jobservice.applications:update` permission)

**Request:**
```bash
curl -X PATCH -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "VIEWED"
  }' http://localhost:8080/api/v1/admin/applications/f74d0446-5662-437a-b9ab-f6138b7b1849
```

**Response:**
*   **Status:** `200 OK`
*   **Body:**
```json
{
  "success": true,
  "id": "f74d0446-5662-437a-b9ab-f6138b7b1849"
}
```

### DELETE /api/v1/admin/applications/{id}

**Description:** Deletes an application.

**Authentication:** Bearer Token (`jobservice.applications:delete` permission)

**Request:**
```bash
curl -X DELETE -H "Authorization: Bearer <YOUR_TOKEN>" http://localhost:8080/api/v1/admin/applications/f74d0446-5662-437a-b9ab-f6138b7b1849 -i
```

**Response:**
*   **Status:** `204 No Content`
