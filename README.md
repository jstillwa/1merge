# 1Merge

A CLI tool to merge duplicate 1Password login entries.

## Prerequisites

- Go 1.21 or higher
- 1Password CLI (`op`) must be installed and available in PATH
  - Install from: <https://developer.1password.com/docs/cli/get-started/>
  - After installation, sign in with: `op signin`
  - You must be signed into 1Password CLI before running 1merge

## Installation

1. Clone the repository
2. Navigate to the project directory
3. Download dependencies:

   ```bash
   go mod download
   ```

4. Build the application:

   ```bash
   go build
   ```

## Usage

### Basic Usage

```bash
./1merge [flags]
```

### Interactive Mode

By default, 1merge runs in interactive mode, prompting you for confirmation before merging each duplicate group.

For each group of duplicates found, you'll see:

- The domain and username that identifies the group
- A list of all duplicate items with their titles, IDs, and last updated timestamps
- A prompt asking whether to merge: `(y/n/q)`

Response options:

- `y` (yes): Merge this group of duplicates
- `n` (no): Skip this group and move to the next
- `q` (quit): Exit the program immediately without processing remaining groups

Example interaction:

```Example
Found 3 duplicate groups

=== Duplicate Group: google.com | user@example.com ===
Found 3 duplicate items:
  1. "Google Account" (ID: abc12345...) - Updated: 2024-01-15 14:30:00
     URL: https://accounts.google.com
  2. "Gmail Login" (ID: def67890...) - Updated: 2024-01-10 09:15:00
     URL: https://mail.google.com
  3. "Google" (ID: ghi11121...) - Updated: 2023-12-20 16:45:00
     URL: https://google.com

Merge these items? (y/n/q): y
Successfully merged 2 items into abc12345

[Next group appears...]
```

### Flags

- `--vault` (string): Specifies which 1Password vault to scan. If not specified, uses the default vault.
- `--dry-run` (bool): Prevents any write operations and only prints what would happen.
- `--auto` (bool): Automatically merges all duplicates without prompting (skips interactive mode).

### Merge Operation

The merge operation works by:

1. **Winner Selection**: The most recently updated item (by `updated_at` timestamp) becomes the winner
2. **Field Merging**: Unique fields from duplicate items are merged into the winner. The winner's `AdditionalInformation` (notes) field is preserved.
3. **Conflict Handling**: Conflicting fields (same label but different values) are preserved in an "Archived Conflicts" section
4. **URL Consolidation**: All unique URLs from duplicate items are added to the winner, preserving URL labels. If multiple items have primary URLs, only the winner's primary URL remains marked as primary.
5. **Archive Duplicates**: The duplicate items are archived (not permanently deleted) and can be restored from 1Password Archive

### Understanding the Merge Process

When you confirm a merge (or use `--auto` mode), 1merge:

1. **Selects a Winner**: The item with the most recent `updated_at` timestamp
2. **Merges All Losers**: Iteratively merges each duplicate into the winner:
   - Unique fields from each duplicate are added to the winner
   - Conflicting fields (same label, different values) are preserved in "Archived Conflicts" section
   - All unique URLs are consolidated
3. **Updates the Winner**: The merged item replaces the winner in your vault
4. **Archives Duplicates**: All other items in the group are archived (not deleted)

The merge is atomic per group: if any step fails, the group is skipped and processing continues with the next group.

### Dry-Run Mode

Use `--dry-run` to preview changes without modifying your vault:

```bash
./1merge --dry-run
```

Dry-run mode will:

- Display the merged item in JSON format
- Show which items would be archived
- List all operations that would occur
- Make **no** actual changes to your vault

Example output:

```JSON
[DRY RUN] Would edit item: <winner_id> (<winner_title>)
{
  "id": "<winner_id>",
  "title": "<winner_title>",
  "fields": [...],
  "urls": [...],
  ...
}
[DRY RUN] Would archive item: <loser_id> (<loser_title>)
[DRY RUN] Would archive item: <loser_id_2> (<loser_title_2>)
```

### Examples

Run in interactive mode (default):

```bash
./1merge
```

Preview merges without making changes:

```bash
./1merge --dry-run
```

Merge duplicates in a specific vault:

```bash
./1merge --vault "MyVault"
```

Automatically merge duplicates without confirmation:

```bash
./1merge --auto
```

Combine flags for dry-run in a specific vault:

```bash
./1merge --vault "MyVault" --dry-run
```

### Safety Considerations

- **Always test with `--dry-run` first** to review what will be changed
- Interactive mode allows you to review each duplicate group before merging
- Use 'n' to skip groups you're unsure about, or 'q' to exit and review your vault first
- Archived items can be restored from the 1Password Archive if needed
- The tool requires the `op` CLI to be installed and authenticated
- Merge operations are fail-fast: if archiving a loser fails, no further items are archived
- The tool uses temporary files for item updates, which are automatically cleaned up after each operation

### Known Limitations

- The tool uses the 1Password CLI's template file approach for item updates, which requires write access to the system's temp directory
- Section references in merged items rely on the `op` CLI automatically creating sections when fields reference them
- URL label preservation requires 1Password CLI v2.0+ (earlier versions may not support the `label` field)

## Verification

To verify the installation and 1Password CLI connection, run:

```bash
go run main.go --dry-run
```

This will print:

- `Dry Run Mode Enabled`
- A success message confirming connection to 1Password CLI

If you encounter authentication errors, ensure you've signed in with `op signin`.

## Development

### Project Architecture

The project uses a modular architecture with an internal `op` package that encapsulates all 1Password CLI interactions:

- **`internal/op/client.go`**: Core wrapper interface for executing `op` CLI commands
  - `Client` interface: Defines `RunOpCmd` method
  - `DefaultClient`: Production implementation using `os/exec`
  - Injectable design enables testing with mock clients
  - `CheckOpInstalled()`: Verifies op binary is in PATH
  - `CheckOpSignedIn()`: Verifies authentication status
  - `VerifyOpReady()`: Combined check for installation and authentication

- **`internal/items/`**: Core business logic for fetching, grouping, merging, and applying changes
  - `fetcher.go`: Retrieves login items from 1Password
  - `grouper.go`: Groups duplicates by base domain and username
  - `merger.go`: Implements superset merge strategy
  - `applier.go`: Applies merged items back to 1Password vault using template files

### Testing

Run unit tests:

```bash
go test ./...
```

Run tests for specific packages:

```bash
go test ./internal/op/
go test ./internal/items/
```

Note: Unit tests do not require the `op` CLI or authentication. The dry-run mode tests verify merge logic without executing actual 1Password operations.

To run integration tests that interact with actual 1Password:

```bash
go test -v ./internal/items -run Integration
```

**Important**: Integration tests require:

- `op` CLI to be installed and in PATH
- Active 1Password CLI authentication (`op signin`)
- A test vault available in your 1Password account
