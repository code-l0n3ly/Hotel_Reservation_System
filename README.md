# Graduation Project API Documentation

Welcome to the API documentation for the Graduation Project. This document provides comprehensive information on the RESTful API endpoints for various handlers. 

Use the following URL to connect to the API: `aiscbackend-production.up.railway.app`

## Table of Contents

1. [UserHandler API](#userhandler-api)
2. [UnitHandler API](#unithandler-api)
3. [ReviewHandler API](#reviewhandler-api)
4. [ReportHandler API](#reporthandler-api)
5. [BookingHandler API](#bookinghandler-api)
6. [MaintenanceTicketHandler API](#maintenancetickethandler-api)
7. [MessageHandler API](#messagehandler-api)
8. [FinancialTransactionHandler API](#financialtransactionhandler-api)
9. [PropertyHandler API](#propertyhandler-api)

---

## UserHandler API

### Endpoints

#### `POST /user`
Creates a new user.

##### Parameters
- `Name`: string
- `Email`: string
- `PhoneNumber`: string
- `Password`: string
- `UserRole`: string

##### Returns
- The created User object

#### `GET /user/{id}`
Retrieves a user by ID.

##### Parameters
- `id`: string (path parameter)

##### Returns
- The requested User object

#### `PUT /user/{id}`
Updates a user by ID.

##### Parameters
- `id`: string (path parameter)
- `Name`: string
- `Email`: string
- `UserRole`: string

##### Returns
- A message indicating the update was successful

#### `DELETE /user/{id}`
Deletes a user by ID.

##### Parameters
- `id`: string (path parameter)

##### Returns
- A message indicating the deletion was successful

#### `GET /user`
Retrieves all users.

##### Returns
- An array of User objects

#### `POST /user/login`
Authenticates a user.

##### Parameters
- `Email`: string
- `Password`: string

##### Returns
- A message indicating the login was successful or failed

## PropertyHandler API

### Endpoints

#### `POST /property`
Creates a new property.

##### Parameters
- `OwnerID`: string
- `Name`: string
- `Address`: string
- `Type`: string
- `Description`: string
- `Rules`: string
- `Photos`: array of base64-encoded strings (images)

##### Returns
- The created Property object

#### `GET /property/{id}`
Retrieves a property by ID.

##### Parameters
- `id`: string (path parameter)

##### Returns
- The requested Property object

#### `PUT /property/{id}`
Updates a property by ID.

##### Parameters
- `id`: string (path parameter)
- `OwnerID`: string
- `Name`: string
- `Address`: string
- `Type`: string
- `Description`: string
- `Rules`: string
- `Photos`: array of base64-encoded strings (images)

##### Returns
- A message indicating the update was successful

#### `DELETE /property/{id}`
Deletes a property by ID.

##### Parameters
- `id`: string (path parameter)

##### Returns
- A message indicating the deletion was successful

#### `GET /property/user/{id}`
Retrieves all properties by UserID.

##### Parameters
- `id`: string (path parameter)

##### Returns
- An array of Property objects owned by the user

---

## UnitHandler API

### Endpoints

#### `POST /unit`
Creates a new unit.

##### Parameters
- `PropertyID`: string
- `Name`: string
- `Description`: string
- `OccupancyStatus`: string
- `StructuralProperties`: string
- `Images`: array of base64-encoded strings (images)

##### Returns
- The created Unit object

#### `GET /unit/{id}`
Retrieves a unit by ID.

##### Parameters
- `id`: string (path parameter)

##### Returns
- The requested Unit object

#### `PUT /unit/{id}`
Updates a unit by ID.

##### Parameters
- `id`: string (path parameter)
- `PropertyID`: string
- `Name`: string
- `Description`: string
- `OccupancyStatus`: string
- `StructuralProperties`: string
- `Images`: array of base64-encoded strings (images)

##### Returns
- A message indicating the update was successful

#### `DELETE /unit/{id}`
Deletes a unit by ID.

##### Parameters
- `id`: string (path parameter)

##### Returns
- A message indicating the deletion was successful

#### `GET /unit`
Retrieves all units.

##### Returns
- An array of Unit objects