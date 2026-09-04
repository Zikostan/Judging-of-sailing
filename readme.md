# Application for judging of sailing

## API Endpoints

All API routes are prefixed with `/api/v1`. Protected routes require a valid JWT token in the `Authorization` header.

### Public Routes

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/auth/register` | Register a new user account |
| `POST` | `/api/v1/auth/login` | Authenticate and receive a JWT token |

### Swagger Documentation

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/swagger/doc.json` | Swagger/OpenAPI specification in JSON format |
| `GET` | `/swagger/` | Swagger UI interactive documentation |

### Protected Routes

#### Regatta

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/regattas` | List all regattas |
| `POST` | `/api/v1/regattas` | Create a new regatta |
| `GET` | `/api/v1/regattas/{id}` | Get a regatta by ID |
| `PUT` | `/api/v1/regattas/{id}` | Update an existing regatta |
| `DELETE` | `/api/v1/regattas/{id}` | Delete a regatta |

#### Category

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/categories` | List all categories |
| `POST` | `/api/v1/categories` | Create a new category |
| `GET` | `/api/v1/categories/{id}` | Get a category by ID |
| `PUT` | `/api/v1/categories/{id}` | Update an existing category |
| `DELETE` | `/api/v1/categories/{id}` | Delete a category |

#### Sailor

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/regattas/{regatta_id}/sailors` | List sailors for a specific regatta |
| `GET` | `/api/v1/regattas/{regatta_id}/sailors/find` | Find a sailor by class number within a regatta |
| `POST` | `/api/v1/sailors` | Create a new sailor |
| `GET` | `/api/v1/sailors/{id}` | Get a sailor by ID |
| `PUT` | `/api/v1/sailors/{id}` | Update an existing sailor |
| `DELETE` | `/api/v1/sailors/{id}` | Delete a sailor |

#### Race Group

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/regattas/{regatta_id}/race-groups` | List race groups for a specific regatta |
| `POST` | `/api/v1/race-groups` | Create a new race group |
| `GET` | `/api/v1/race-groups/{id}` | Get a race group by ID |
| `POST` | `/api/v1/race-groups/{id}/start` | Start a race group |
| `POST` | `/api/v1/race-groups/{id}/finish` | Finish a race group |
| `DELETE` | `/api/v1/race-groups/{id}` | Delete a race group |

#### Race

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/race-groups/{group_id}/races` | List races for a specific race group |
| `POST` | `/api/v1/races` | Record a new race |
| `GET` | `/api/v1/races/{id}` | Get a race by ID |
| `PATCH` | `/api/v1/races/{id}/result` | Update the result of a race |
| `POST` | `/api/v1/races/{id}/link-sailor` | Link a sailor to a race |
| `DELETE` | `/api/v1/races/{id}` | Delete a race |

#### Document

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/sailors/{sailor_id}/documents` | List documents for a specific sailor |
| `POST` | `/api/v1/documents` | Upload a new document |
| `GET` | `/api/v1/documents/{id}` | Get a document by ID |
| `POST` | `/api/v1/documents/{id}/approve` | Approve a document |
| `POST` | `/api/v1/documents/{id}/reject` | Reject a document |
| `POST` | `/api/v1/documents/{id}/request-revision` | Request a revision of a document |
| `DELETE` | `/api/v1/documents/{id}` | Delete a document |