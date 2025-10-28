# Security Policy

## Supported Versions

Phoenix TUI Framework is currently in beta. We provide security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.0-beta.x | :white_check_mark: |
| < 0.1.0 | :x:                |

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability in Phoenix TUI Framework, please report it responsibly.

### How to Report

**DO NOT** open a public GitHub issue for security vulnerabilities.

Instead, please report security issues by emailing:

**security@phoenix-tui.org** (preferred)

Or open a **private security advisory** on GitHub:
https://github.com/phoenix-tui/phoenix/security/advisories/new

### What to Include

Please include the following information in your report:

- **Description** of the vulnerability
- **Steps to reproduce** the issue
- **Affected versions** (which libraries/versions are impacted)
- **Potential impact** (what can an attacker do?)
- **Suggested fix** (if you have one)
- **Your contact information** (for follow-up questions)

### Response Timeline

- **Initial Response**: Within 48 hours
- **Triage & Assessment**: Within 1 week
- **Fix & Disclosure**: Coordinated with reporter

We aim to:
1. Acknowledge receipt within 48 hours
2. Provide an initial assessment within 1 week
3. Work with you on a coordinated disclosure timeline
4. Credit you in the security advisory (unless you prefer to remain anonymous)

## Security Best Practices

When using Phoenix TUI Framework:

### Terminal Input Validation

Phoenix components handle user input. Always validate and sanitize:

```go
// ❌ BAD - Don't trust raw input
cmd := exec.Command("sh", "-c", userInput)

// ✅ GOOD - Validate first
if !isValidCommand(userInput) {
    return errors.New("invalid command")
}
```

### ANSI Escape Sequence Injection

Phoenix generates ANSI sequences. Be careful with untrusted content:

```go
// ❌ BAD - Raw user content in terminal
fmt.Fprintf(terminal, "\x1b[31m%s\x1b[0m", userContent)

// ✅ GOOD - Use Phoenix's safe rendering
style.New().Foreground(color.Red).Render(userContent)
```

### Mouse Event Handling

Phoenix parses mouse events. Validate coordinates:

```go
// Phoenix automatically validates mouse positions
// But verify they're within your UI bounds
if mouse.X >= 0 && mouse.X < width && mouse.Y >= 0 && mouse.Y < height {
    // Safe to use
}
```

### Clipboard Operations

Phoenix clipboard supports OSC 52. Be aware:

- **Reading clipboard**: May expose sensitive data (passwords, tokens)
- **Writing clipboard**: User must trust your application
- **OSC 52 in SSH**: Terminal emulator controls permissions

```go
// ✅ Ask user permission before clipboard access
if userConfirmedClipboardAccess {
    clipboard.Write("data")
}
```

## Known Security Considerations

### 1. Terminal Escape Sequence Parsing

Phoenix parses ANSI/CSI sequences from terminal input. Malicious terminals could send crafted sequences.

**Mitigation**: Phoenix validates all parsed sequences and bounds-checks all values.

### 2. Mouse Protocol Buffer Overflow

SGR mouse protocol uses unbounded integers for coordinates.

**Mitigation**: Phoenix limits coordinates to `int` range and validates against terminal dimensions.

### 3. Clipboard OSC 52 Injection

OSC 52 sequences can write to clipboard without user interaction in some terminals.

**Mitigation**:
- Phoenix only sends OSC 52 when explicitly requested by application
- Applications should ask user permission before clipboard writes
- Document recommends users configure terminals with OSC 52 confirmation dialogs

### 4. Dependency Security

Phoenix has minimal dependencies. We monitor security advisories for:

- `github.com/rivo/uniseg` - Unicode grapheme cluster segmentation
- `golang.org/x/sys` - Platform-specific terminal operations
- `golang.org/x/term` - Terminal raw mode management (Unix)

**Monitoring**:
- Dependabot enabled (when public)
- Weekly dependency audit
- Automated CVE scanning in CI

## Security Disclosure History

No security vulnerabilities have been reported or fixed yet (project is in beta).

When vulnerabilities are addressed, they will be listed here with:
- **CVE ID** (if assigned)
- **Affected versions**
- **Fixed in version**
- **Severity** (Critical/High/Medium/Low)
- **Credit** to reporter

## Security Contact

- **Email**: security@phoenix-tui.org
- **GitHub**: https://github.com/phoenix-tui/phoenix/security
- **PGP Key**: (Coming soon)

## Bug Bounty Program

Phoenix does not currently have a bug bounty program. We rely on responsible disclosure from the security community.

If you report a valid security vulnerability:
- ✅ Public credit in security advisory (if desired)
- ✅ Acknowledgment in CHANGELOG
- ✅ Our gratitude and a virtual high-five 🙌

---

**Thank you for helping keep Phoenix TUI Framework secure!** 🔒
