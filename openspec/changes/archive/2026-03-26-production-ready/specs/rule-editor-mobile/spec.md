## MODIFIED Requirements

### Requirement: Delete button position on small screens

The system SHALL position the delete button at the top-right corner of each rule item on small screens.

#### Scenario: Small screen delete button layout
- **WHEN** the screen width is 768px or less
- **THEN** the delete button is positioned at the top-right corner of the rule item
- **AND** the delete button does not occupy a separate row

#### Scenario: Desktop delete button layout
- **WHEN** the screen width is greater than 768px
- **THEN** the delete button remains at the end of the rule row

#### Scenario: Delete button remains functional
- **WHEN** the delete button is repositioned
- **THEN** clicking it still removes the rule from the list
