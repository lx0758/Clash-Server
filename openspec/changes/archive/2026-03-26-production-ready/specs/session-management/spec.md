## ADDED Requirements

### Requirement: Session management via gin middleware

The system SHALL use `gin-contrib/sessions` with `memstore` for session management.

#### Scenario: Session initialized on startup
- **WHEN** the server starts
- **THEN** session middleware is configured with a generated secret key

#### Scenario: Session accessible in handlers
- **WHEN** a handler calls `sessions.Default(c)`
- **THEN** a valid session object is returned

### Requirement: Login user helper function

The system SHALL provide a `LoginUser(c, userId)` helper function in the middleware package.

#### Scenario: User login
- **WHEN** `LoginUser(c, userId)` is called
- **THEN** the user ID is stored in the session

### Requirement: Logout user helper function

The system SHALL provide a `LogoutUser(c)` helper function in the middleware package.

#### Scenario: User logout
- **WHEN** `LogoutUser(c)` is called
- **THEN** the user ID is removed from the session

### Requirement: Get user ID helper function

The system SHALL provide a `GetUserID(c)` helper function that returns the current user's ID.

#### Scenario: Authenticated user
- **WHEN** `GetUserID(c)` is called for an authenticated session
- **THEN** the user ID is returned

#### Scenario: Unauthenticated user
- **WHEN** `GetUserID(c)` is called for an unauthenticated session
- **THEN** zero (0) is returned

## REMOVED Requirements

### Requirement: Manual session store initialization

**Reason**: Replaced by gin-contrib/sessions middleware
**Migration**: Remove `InitSessionStore()` calls, use `sessions.Sessions()` middleware instead

### Requirement: Global session store variable

**Reason**: Session store is now managed by gin middleware
**Migration**: Use `sessions.Default(c)` instead of `middleware.Store.Get()`
