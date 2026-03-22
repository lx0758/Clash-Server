## ADDED Requirements

### Requirement: Static assets embedded in binary

The system SHALL embed all static web assets (HTML, JS, CSS, images) into the binary using Go's `embed` package.

#### Scenario: Assets accessible after build
- **WHEN** the server binary is built
- **THEN** all files under `server/res/web/` are embedded in the binary

#### Scenario: Assets served without external files
- **WHEN** the server starts without `server/res/web/` directory present
- **THEN** static assets are still accessible via HTTP

### Requirement: SPA routing fallback

The system SHALL serve `index.html` for any unmatched routes to support client-side routing.

#### Scenario: SPA route request
- **WHEN** a request is made to a non-API route like `/dashboard` or `/settings`
- **THEN** the system returns `index.html` with status 200

#### Scenario: Static file request
- **WHEN** a request is made to `/assets/index-xxx.js`
- **THEN** the system returns the actual JS file with correct MIME type

### Requirement: Correct MIME types

The system SHALL return appropriate MIME types for static assets.

#### Scenario: JavaScript file
- **WHEN** a `.js` file is requested
- **THEN** the response Content-Type is `text/javascript` or `application/javascript`

#### Scenario: CSS file
- **WHEN** a `.css` file is requested
- **THEN** the response Content-Type is `text/css`

#### Scenario: HTML file
- **WHEN** an `.html` file is requested
- **THEN** the response Content-Type is `text/html`
