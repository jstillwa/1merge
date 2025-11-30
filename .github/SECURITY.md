# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in 1merge, please **do not** create a public GitHub issue. Instead, please report it privately using GitHub's security advisory feature.

### How to Report

1. Go to the [Security tab](../../security/advisories) of this repository
2. Click "Report a vulnerability"
3. Fill in the vulnerability details with:
   - A clear description of the vulnerability
   - Steps to reproduce (if applicable)
   - Potential impact
   - Suggested fix (if you have one)

### What to Expect

- We will acknowledge receipt of your report within 48 hours
- We will investigate and provide updates on progress
- We will work with you to understand and resolve the issue
- We ask that you maintain confidentiality until we release a fix

## Security Considerations

### 1Password CLI Integration

This tool requires the 1Password CLI (`op`) to be installed and authenticated. Be aware that:

- The tool operates with the permissions of the authenticated 1Password account
- The `--auto` flag merges duplicates without confirmation - use with caution
- Always use `--dry-run` first to preview changes before applying them
- Archived items can be restored from the 1Password Archive

### Best Practices

- Keep your 1Password CLI updated to the latest version
- Maintain active 1Password authentication
- Review merged items after operations to ensure correctness
- Use `--dry-run` in automated environments before enabling actual merges
- Restrict access to your 1Password vault for sensitive credentials

## Dependencies

This project depends on:
- [Cobra](https://github.com/spf13/cobra) - CLI framework

We keep dependencies up to date with Dependabot. Security updates are prioritized and applied promptly.

## Contact

For questions about security practices, you can open a regular issue on the repository.
