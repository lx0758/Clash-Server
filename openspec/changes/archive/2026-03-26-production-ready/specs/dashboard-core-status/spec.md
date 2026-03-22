## MODIFIED Requirements

### Requirement: Core status display format

The system SHALL display core status in a simplified format.

#### Scenario: Core is running
- **WHEN** the core is running
- **THEN** only the version number is displayed (e.g., "v1.18.0")

#### Scenario: Core is stopped
- **WHEN** the core is not running
- **THEN** "已停止" is displayed

#### Scenario: Version unavailable
- **WHEN** the core is running but version is unknown
- **THEN** "运行中" is displayed
