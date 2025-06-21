# Welcome to *AI-Doc-Optimizer*

<p align="center">
    <img src="static/astrolabe.jpg" alt="resolve"/>
</p>

AI-Doc-Optimizer transforms human-readable documentation into AI/RAG-optimized content. Inspired by Kapa.ai documentation best practices.

## Features

- **Context Analysis**: Detects content that depends on previous context
- **Semantic Optimization**: Ensures product names and specific terminology for better AI discoverability  
- **Knowledge Gap Detection**: Identifies implicit assumptions that may confuse AI systems
- **Visual Content Auditing**: Flags visual dependencies without text alternatives
- **Structure Validation**: Checks heading hierarchy and content organization
- **Multi-format Support**: Works with Markdown, HTML, reStructuredText, and plain text

## Installation

### Binary Release
```bash
# Download from releases page
curl -L https://github.com/ghartsel/ai-doc-optimizer/releases/latest/download/ai-doc-optimizer-linux-amd64.tar.gz | tar xz
sudo mv ai-doc-optimizer /usr/local/bin/
```

### Build from Source
```bash
git clone https://github.com/ghartsel/ai-doc-optimizer.git
cd ai-doc-optimizer
go mod tidy
go build -o ai-doc-optimizer
```

## Quick Start

```bash
# Analyze a single file
ai-doc-optimizer docs/README.md

# Analyze directory recursively  
ai-doc-optimizer -recursive docs/

# Use custom configuration
ai-doc-optimizer -config .ai-doc-optimizer.yml docs/

# Output in JSON format
ai-doc-optimizer -output json docs/
```

## Arguments

```bash
  -config string
      Path to configuration file
  -fix
      Attempt to automatically fix issues
  -output string
      Output format: standard (default), json
  -recursive
      Process directories recursively
```

## Configuration

Create `.ai-doc-optimizer.yml` in your project root:

```yaml
StylesPath: "./styles"
MinWordCount: 10

Rules:
  - Name: "contextual-dependency"
    Description: "Detect sections that depend on previous context"
    Pattern: '(?i)\b(this|that|these|those)\b(?:\s+\w+){0,3}\s+(?:will|should|must)'
    Severity: "warning"
    Type: "suggest"
```

## Common Issues Detected

### Contextual Dependencies
❌ **Bad**: "This will configure the webhook endpoint."
✅ **Good**: "This CloudSync configuration will set up the webhook endpoint."

### Missing Product Context  
❌ **Bad**: "## Installation"
✅ **Good**: "## CloudSync Installation"

### Visual Dependencies
❌ **Bad**: "See the diagram above for the workflow."
✅ **Good**: "The CloudSync workflow includes: 1. Authentication, 2. Data upload, 3. Processing confirmation."

### Implicit Knowledge
❌ **Bad**: "Simply configure the endpoint URL."
✅ **Good**: "Configure the endpoint URL in Settings > Webhooks by entering your HTTPS endpoint."

## Output Formats

- **Standard**: Human-readable console output
- **JSON**: Machine-readable for CI integration  

### Standard

This format makes it easy to identify and fix AI optimization issues in documentation.

```
{file}:{line}:{column}: {SEVERITY} [{rule}] {message}
    Suggestion: {improvement_suggestion}
```

**Location Information:**
- **File path**: Exact file with the issue
- **Line number**: Specific line for quick navigation
- **Column number**: Precise character position

**Issue Classification:**
- **SEVERITY**: Color-coded levels (ERROR/WARNING/SUGGESTION)
- **Rule name**: Which optimization rule triggered
- **Clear message**: Human-readable explanation

**Actionable Guidance:**
- **Suggestion line**: Specific improvement recommendation
- **Empty line separator**: Clean visual separation between issues

**IDE Integration Ready:**
- Format matches standard linter output (file:line:column)
- Most editors can parse this for clickable navigation
- Compatible with VS Code problem matcher patterns

### JSON

The JSON format makes it easy to integrate with CI systems, create dashboards, or build additional tooling around the AI documentation optimization results!

**Benefits:**
- **CI/CD Integration**: Easy parsing for automated pipelines
- **Dashboard Visualization**: Summary stats for reporting
- **Filtering**: JSON structure enables post-processing
- **Pretty Printing**: Indented output for readability

```json
{
  "version": "1.0.0",
  "issues": [
    {
      "File": "docs/api.md",
      "Line": 15,
      "Column": 12,
      "Rule": "contextual-dependency",
      "Message": "This text may depend on previous context...",
      "Severity": "warning",
      "Suggestion": "Replace contextual references with specific details",
      "OriginalText": "this will configure"
    }
  ],
  "summary": {
    "total": 8,
    "by_severity": {
      "warning": 5,
      "error": 1,
      "suggestion": 2
    },
    "by_rule": {
      "contextual-dependency": 3,
      "visual-dependency": 1,
      "implicit-knowledge": 4
    }
  }
}
```

## Integration

### CI/CD Pipeline (GitHub Actions)
```yaml
- name: Check Documentation
  run: |
    ai-doc-optimizer -recursive docs/
    if [ $? -ne 0 ]; then exit 1; fi
```

### Pre-commit Hook
```bash
#!/bin/sh
ai-doc-optimizer $(git diff --cached --name-only --diff-filter=ACM | grep -E '\.(md|html)$')
```

## Similar Tools

- [Vale](https://vale.sh/) - Prose linting with style guides
- [textlint](https://textlint.github.io/) - JavaScript-based text linting
- [write-good](https://github.com/btford/write-good) - Naive English linter

## TODO

- Chunk boundary analysis
- Semantic similarity detection
- Error message pattern recognition
- Table semantics recognition

## License

MIT License - see [LICENSE](LICENSE) file for details.