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
curl -L https://github.com/yourorg/ai-doc-optimizer/releases/latest/download/ai-doc-optimizer-linux-amd64.tar.gz | tar xz
sudo mv ai-doc-optimizer /usr/local/bin/
```

### Build from Source
```bash
git clone https://github.com/yourorg/ai-doc-optimizer.git
cd ai-doc-optimizer
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
      Output format (standard, json, sarif) (default "standard")
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

## Output Formats

- **Standard**: Human-readable console output
- **JSON**: Machine-readable for CI integration  
- **SARIF**: Static Analysis Results Interchange Format

## Similar Tools

- [Vale](https://vale.sh/) - Prose linting with style guides
- [textlint](https://textlint.github.io/) - JavaScript-based text linting
- [write-good](https://github.com/btford/write-good) - Naive English linter

## License

MIT License - see [LICENSE](LICENSE) file for details.