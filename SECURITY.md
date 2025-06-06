# Security Policy

## Supported Versions

| Version | Supported          | Status |
| ------- | ------------------ | ------ |
| 1.1.0   | :white_check_mark: | Latest |
| 1.0.3   | :white_check_mark: | Supported |
| 1.0.2   | :warning:          | Limited Support |
| < 1.0.2 | :x:                | Unsupported |

## Reporting a Vulnerability

If you discover a security vulnerability, please follow these steps:

1. **Do not** open a public issue
2. Email details to: jad@madi.se
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Any suggested fixes
   - Go-OPML version affected
   - Operating system and Go version

## Response Timeline

- **Acknowledgment**: Within 48 hours
- **Initial Assessment**: Within 1 week
- **Fix Development**: Depends on severity (Critical: 1-3 days, High: 1 week, Medium: 2 weeks)
- **Release**: Security fixes are released as patch versions
- **Disclosure**: After fix is released and users have time to update

## Security Best Practices

When using Go-OPML:

1. **Input Validation**: Always validate OPML input files before processing
2. **Network Security**: Use appropriate timeout settings for RSS fetching (default: 30s)
3. **Permissions**: Run with minimal required permissions
4. **Updates**: Keep Go-OPML updated to the latest version (currently v1.1.0)
5. **Dependencies**: Regularly check for dependency updates using `go list -u -m all`
6. **Isolation**: Consider running in containerized environments for additional isolation
7. **Monitoring**: Monitor RSS fetch timeouts and failures for unusual patterns

## Known Security Considerations

- **RSS Feed Fetching**: Go-OPML fetches external RSS feeds, which could potentially expose you to malicious content
- **File Processing**: OPML files are parsed as XML - ensure input files are from trusted sources
- **Network Requests**: The tool makes HTTP requests to RSS endpoints - consider network policies
- **Concurrent Processing**: Default concurrency is 5 - adjust based on your security requirements

## Changelog of Security Updates

### v1.1.0 (Current)
- Enhanced CI/CD pipeline with automated security checks
- Comprehensive linting configuration for code quality
- Improved project governance and security policies
- Added automated testing across multiple Go versions

### v1.0.3
- Updated 11 dependencies for enhanced security
- Improved timeout handling for RSS fetching
- Enhanced error handling to prevent information disclosure

### v1.0.2
- Dependency security updates
- Improved input validation

For the complete changelog, see [README.md](README.md#changelog).
