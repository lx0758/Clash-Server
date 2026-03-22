## MODIFIED Requirements

### Requirement: Active subscription visual indicator

The system SHALL indicate active subscriptions with a subtle visual style.

#### Scenario: Active subscription card styling
- **WHEN** a subscription is active
- **THEN** the card displays a green left border (3px solid)
- **AND** no green box-shadow is applied

#### Scenario: Non-active subscription card styling
- **WHEN** a subscription is not active
- **THEN** the card displays no special border indicator

#### Scenario: Active tag still visible
- **WHEN** a subscription is active
- **THEN** a green "激活中" tag is still displayed next to the subscription name
