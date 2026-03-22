## ADDED Requirements

### Requirement: PC sidebar logout entry

The system SHALL provide a logout entry in the PC sidebar for user authentication management.

#### Scenario: Logged in user sees logout option
- **WHEN** a logged in user views the PC sidebar
- **THEN** the sidebar displays the username and a logout button at the bottom

#### Scenario: User clicks logout
- **WHEN** user clicks the logout button in the sidebar
- **THEN** the user session is terminated
- **AND** the user is redirected to the login page

#### Scenario: Consistent with mobile experience
- **WHEN** comparing PC sidebar logout with mobile header dropdown
- **THEN** both provide the same logout functionality
