# UserHandler API

## Endpoints

### POST /user

Creates a new user.

#### Parameters

- `Name`: string
- `Email`: string
- `Password`: string

#### Returns

- `Status`: string
- `Message`: string
- `Data`: User object (if creation is successful)

### GET /user/{id}

Retrieves a user by ID.

#### Parameters

- `id`: string (path parameter)

#### Returns

- `Status`: string
- `Message`: string
- `Data`: User object (if user is found)

### GET /users

Retrieves all users.

#### Parameters

None

#### Returns

- Array of User objects

### PUT /user/{id}

Updates a user by ID.

#### Parameters

- `id`: string (path parameter)
- `Name`: string
- `Email`: string
- `UserRole`: string

#### Returns

- `Status`: string
- `Message`: string

### DELETE /user/{id}

Deletes a user by ID.

#### Parameters

- `id`: string (path parameter)

#### Returns

- `Status`: string
- `Message`: string

### POST /login

Authenticates a user.

#### Parameters

- `Email`: string
- `Password`: string

#### Returns

- `Status`: string
- `Message`: string

# UnitHandler API

## UNIT

### POST /unit

Creates a new unit.

#### Parameters

- `PropertyID`: string
- `RentalPrice`: float
- `OccupancyStatus`: string
- `StructuralProperties`: string
- `CreateTime`: string

#### Returns

- The created Unit object

### GET /unit/{id}

Retrieves a unit by ID.

#### Parameters

- `id`: string (path parameter)

#### Returns

- The requested Unit object

### PUT /unit/{id}

Updates a unit by ID.

#### Parameters

- `id`: string (path parameter)
- `PropertyID`: string
- `RentalPrice`: float
- `OccupancyStatus`: string
- `StructuralProperties`: string
- `CreateTime`: string

#### Returns

- A message indicating the update was successful

### DELETE /unit/{id}

Deletes a unit by ID.

#### Parameters

- `id`: string (path parameter)

#### Returns

- A message indicating the deletion was successful


# ReviewHandler API

## REVIEW

### POST /review

Creates a new review.

#### Parameters

- `ReviewID`: string
- `UserID`: string
- `UnitID`: string
- `Rating`: float
- `Comment`: string
- `CreateTime`: string

#### Returns

- The created Review object

### GET /review/{id}

Retrieves a review by ID.

#### Parameters

- `id`: string (path parameter)

#### Returns

- The requested Review object

### PUT /review/{id}

Updates a review by ID.

#### Parameters

- `id`: string (path parameter)
- `UserID`: string
- `UnitID`: string
- `Rating`: float
- `Comment`: string
- `CreateTime`: string

#### Returns

- A message indicating the update was successful

### DELETE /review/{id}

Deletes a review by ID.

#### Parameters

- `id`: string (path parameter)

#### Returns

- A message indicating the deletion was successful


# ReportHandler API

## REPORT

### POST /report

Creates a new report.

#### Parameters

- `ReportID`: string
- `UserID`: string
- `Type`: string
- `CreateTime`: string
- `Data`: string

#### Returns

- The created Report object

### GET /report/{id}

Retrieves a report by ID.

#### Parameters

- `id`: string (path parameter)

#### Returns

- The requested Report object

### PUT /report/{id}

Updates a report by ID.

#### Parameters

- `id`: string (path parameter)
- `UserID`: string
- `Type`: string
- `CreateTime`: string
- `Data`: string

#### Returns

- A message indicating the update was successful

### DELETE /report/{id}

Deletes a report by ID.

#### Parameters

- `id`: string (path parameter)

#### Returns

- A message indicating the deletion was successful

# BookingHandler API

## BOOKING

### POST /booking

Creates a new booking.

#### Parameters

- `BookingID`: string
- `UserID`: string
- `UnitID`: string
- `StartDate`: string
- `EndDate`: string
- `CreateTime`: string
- `Summary`: string

#### Returns

- The created Booking object

### GET /booking/{id}

Retrieves a booking by ID.

#### Parameters

- `id`: string (path parameter)

#### Returns

- The requested Booking object

### PUT /booking/{id}

Updates a booking by ID.

#### Parameters

- `id`: string (path parameter)
- `UserID`: string
- `UnitID`: string
- `StartDate`: string
- `EndDate`: string
- `Summary`: string

#### Returns

- A message indicating the update was successful

### DELETE /booking/{id}

Deletes a booking by ID.

#### Parameters

- `id`: string (path parameter)

#### Returns

- A message indicating the deletion was successful
