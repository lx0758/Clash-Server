## ADDED Requirements

### Requirement: Subscription Customization Storage

The system SHALL store subscription customization in a dedicated `SubscriptionCustomization` table with the following fields:
- `subscription_id`: Reference to subscription (unique, 1:1)
- `proxy_insert`, `proxy_append`, `proxy_remove`: Proxy customization YAML
- `proxy_group_insert`, `proxy_group_append`, `proxy_group_remove`: Proxy group customization YAML
- `rule_insert`, `rule_append`, `rule_remove`: Rule customization YAML
- `global_override`: Global config override YAML
- `script`: Post-processing script

#### Scenario: Create customization on first save
- **WHEN** user saves customization for a subscription without existing customization
- **THEN** system creates a new `SubscriptionCustomization` record

#### Scenario: Update existing customization
- **WHEN** user saves customization for a subscription with existing customization
- **THEN** system updates the existing record

### Requirement: Proxy Customization

The system SHALL support proxy customization with insert, append, and remove operations.

#### Scenario: Insert proxies
- **WHEN** user provides YAML in `proxy_insert` field
- **THEN** system inserts proxies at the beginning of subscription's proxy list during merge

#### Scenario: Append proxies
- **WHEN** user provides YAML in `proxy_append` field
- **THEN** system appends proxies at the end of subscription's proxy list during merge

#### Scenario: Remove proxies
- **WHEN** user provides YAML array in `proxy_remove` field
- **THEN** system removes matching proxy names from subscription during merge

### Requirement: Proxy Group Customization

The system SHALL support proxy group customization with insert, append, and remove operations.

#### Scenario: Insert proxy groups
- **WHEN** user provides YAML in `proxy_group_insert` field
- **THEN** system inserts proxy groups at the beginning of subscription's proxy-groups list during merge

#### Scenario: Append proxy groups
- **WHEN** user provides YAML in `proxy_group_append` field
- **THEN** system appends proxy groups at the end of subscription's proxy-groups list during merge

#### Scenario: Remove proxy groups
- **WHEN** user provides YAML array in `proxy_group_remove` field
- **THEN** system removes matching proxy group names from subscription during merge

### Requirement: Rule Customization

The system SHALL support rule customization with insert, append, and remove operations.

#### Scenario: Insert rules
- **WHEN** user provides YAML in `rule_insert` field
- **THEN** system inserts rules at the beginning of subscription's rules list during merge

#### Scenario: Append rules
- **WHEN** user provides YAML in `rule_append` field
- **THEN** system appends rules at the end of subscription's rules list during merge

#### Scenario: Remove rules
- **WHEN** user provides YAML array in `rule_remove` field
- **THEN** system removes exactly matching rules from subscription during merge

### Requirement: Global Config Override

The system SHALL support global configuration override.

#### Scenario: Override global config
- **WHEN** user provides YAML in `global_override` field
- **THEN** system merges the config at the end of merge process (before script execution)

#### Scenario: Invalid YAML format
- **WHEN** user provides invalid YAML in `global_override` field
- **THEN** system returns error message without saving

### Requirement: Post-processing Script

The system SHALL execute post-processing script after all other customizations.

#### Scenario: Execute script with main function
- **WHEN** user provides JavaScript with `main(config)` function
- **THEN** system calls `main` with `config` parameter and uses returned config

#### Scenario: Execute script with direct return
- **WHEN** user provides JavaScript that directly modifies and returns `config`
- **THEN** system uses the returned config

#### Scenario: Script execution error
- **WHEN** script execution fails
- **THEN** system returns error message and does not save customization

### Requirement: Merge Order

The system SHALL apply customizations in the following order:
1. Remove operations (proxy, proxy_group, rule)
2. Insert operations (proxy, proxy_group, rule)
3. Append operations (proxy, proxy_group, rule)
4. Global override
5. Script post-processing

#### Scenario: Apply full customization
- **WHEN** subscription has all customization types configured
- **THEN** system applies them in the defined order

### Requirement: Configuration Validation

The system SHALL validate merged configuration before applying.

#### Scenario: Validate proxies
- **WHEN** merged config is validated
- **THEN** system checks that proxies is an array with valid name field

#### Scenario: Validate proxy groups
- **WHEN** merged config is validated
- **THEN** system checks that proxy-groups is an array with valid name field

#### Scenario: Validate rules format
- **WHEN** merged config is validated
- **THEN** system checks that rules is an array of strings with correct format

#### Scenario: Validate rule target exists
- **WHEN** a rule references a proxy or proxy-group
- **THEN** system checks that the target exists in proxies, proxy-groups, or is a built-in target (DIRECT, REJECT, REJECT-DROP, PASS, COMPATIBLE)

#### Scenario: Validate rule target missing
- **WHEN** a rule references a non-existent proxy or proxy-group
- **THEN** system returns error message indicating which rule has invalid target

#### Scenario: Validate proxy-group proxies exist
- **WHEN** a proxy-group references proxies
- **THEN** system checks that each proxy exists in proxies or other proxy-groups

#### Scenario: Validate proxy-group proxy missing
- **WHEN** a proxy-group references a non-existent proxy or proxy-group
- **THEN** system returns error message indicating which proxy-group has invalid proxy reference

#### Scenario: Validation failure
- **WHEN** validation fails
- **THEN** system restores previous customization and returns error message

### Requirement: Error Handling and Rollback

The system SHALL rollback customization on any error.

#### Scenario: YAML format error
- **WHEN** user provides invalid YAML
- **THEN** system returns error message without saving

#### Scenario: Merge error
- **WHEN** merge process fails
- **THEN** system restores previous customization and returns error message

#### Scenario: Core application error
- **WHEN** applying config to core fails
- **THEN** system restores previous customization, reapplies old config, and returns error message

### Requirement: Customization API

The system SHALL provide REST API for customization management.

#### Scenario: Get customization
- **WHEN** client requests `GET /subscriptions/:id/customization`
- **THEN** system returns customization data or empty object if not exists

#### Scenario: Update customization with validation
- **WHEN** client requests `PUT /subscriptions/:id/customization` with customization data
- **THEN** system validates, saves, merges, and applies config; returns error if any step fails

### Requirement: Frontend Customization Dialog

The system SHALL provide a unified customization dialog with tabs.

#### Scenario: Open customization dialog
- **WHEN** user clicks "自定义配置" from subscription card menu
- **THEN** system displays dialog with tabs: [代理] [代理组] [规则] [全局] [脚本]

#### Scenario: YAML editor with placeholder examples
- **WHEN** user views any tab
- **THEN** system displays textarea with placeholder showing example YAML format

#### Scenario: Error display
- **WHEN** save operation fails
- **THEN** system displays red error alert at top of dialog and does not close dialog

#### Scenario: Frontend YAML validation
- **WHEN** user clicks save button
- **THEN** system validates all YAML fields before sending to backend

## REMOVED Requirements

### Requirement: Rule Management (from rule-management)

**Reason**: Replaced by unified subscription customization.

**Migration**: Use `rule_insert`, `rule_append`, `rule_remove` fields in customization.

### Requirement: Script Management (from script-management)

**Reason**: Replaced by unified subscription customization.

**Migration**: Use `script` field in customization.
