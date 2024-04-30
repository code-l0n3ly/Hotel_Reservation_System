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

#### `POST /users/create`
Creates a new user.

##### Parameters
- `Name`: string
- `Email`: string
- `PhoneNumber`: string
- `Password`: string
- `UserRole`: ENUM('LandLord', 'Tenant', 'MaintenancePresenter')

##### Returns
- The created User object

#### `GET /users/{id}`
Retrieves a user by ID.

##### Parameters
- `id`: string (path parameter)

##### Returns
- The requested User object

#### `PUT /users/{id}`
Updates a user by ID.

##### Parameters
- `id`: string (path parameter)
- `Name`: string
- `Email`: string
- `UserRole`: ENUM('LandLord', 'Tenant', 'MaintenancePresenter')

##### Returns
- A message indicating the update was successful

#### `DELETE /users/{id}`
Deletes a user by ID.

##### Parameters
- `id`: string (path parameter)

##### Returns
- A message indicating the deletion was successful

#### `GET /users`
Retrieves all users.

##### Returns
- An array of User objects

#### `POST /users/login`
Authenticates a user.

##### Parameters
- `Email`: string
- `Password`: string

##### Returns
- A message indicating the login was successful or failed

## PropertyHandler API

### Endpoints

#### `POST /property/create`
Creates a new property.

##### Parameters
- `OwnerID`: string
- `Name`: string
- `Address`: string
- `Type`: ENUM('Residential', 'Commercial')
- `Description`: string
- `Rules`: JSON
- `Photos`: array of base64-encoded strings (images)

##### Returns
- The created Property object

#### `GET /property/{id}`
Retrieves a property by ID.

##### Parameters
- `id`: string (path parameter)

##### Returns
- The requested Property object

#### `GET /property/`
Retrieves all properties.

##### Returns
- An array of Property objects

#### `GET /property/owner/{id}`
Retrieves all properties by UserID.

##### Parameters
- `id`: string (path parameter)

##### Returns
- An array of Property objects owned by the user

#### `GET /property/AllUnits/{id}`
Retrieves all units by PropertyID.

##### Parameters
- `id`: string (path parameter)

##### Returns
- An array of Unit objects owned by the property

#### `GET /property/ByType/{type}`
Retrieves all properties by Type.

##### Parameters
- `type`: string (path parameter)

##### Returns
- An array of Property objects of the specified type

#### `PUT /property/{id}`
Updates a property by ID.

##### Parameters
- `id`: string (path parameter)
- `OwnerID`: string
- `Name`: string
- `Address`: string
- `Type`: ENUM('Residential', 'Commercial')
- `Description`: string
- `Rules`: JSON
- `Photos`: array of base64-encoded strings (images)

##### Returns
- A message indicating the update was successful

#### `DELETE /property/{id}`
Deletes a property by ID.

##### Parameters
- `id`: string (path parameter)

##### Returns
- A message indicating the deletion was successful

---

## UnitHandler API

### Endpoints

#### `POST /unit`
Creates a new unit.

##### Parameters
- `PropertyID`: string
- `Name`: string
- `Description`: string
- `OccupancyStatus`: ENUM('Occupied', 'Available')
- `StructuralProperties`: JSON
- `RentalPrice`: float
- `Rating`: float
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
- `OccupancyStatus`: ENUM('Occupied', 'Available')
- `StructuralProperties`: JSON
- `RentalPrice`: float
- `Rating`: float
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