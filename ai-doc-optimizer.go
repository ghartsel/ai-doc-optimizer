// AI-Doc-Optimizer: Transform human documentation for AI/RAG consumption

package main

import (
//    "bufio"
    "flag"
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "regexp"
    "strings"
//    "unicode"

    "gopkg.in/yaml.v3"
)

// Config represents the main configuration structure
type Config struct {
    StylesPath   string            `yaml:"StylesPath"`
    MinWordCount int               `yaml:"MinWordCount"`
    Formats      map[string]Format `yaml:"Formats"`
    Rules        []Rule            `yaml:"Rules"`
}

// Format defines file format configurations
type Format struct {
    Extensions []string `yaml:"Extensions"`
    Parser     string   `yaml:"Parser"`
}

// Rule defines transformation rules
type Rule struct {
    Name        string `yaml:"Name"`
    Description string `yaml:"Description"`
    Pattern     string `yaml:"Pattern"`
    Replacement string `yaml:"Replacement,omitempty"`
    Severity    string `yaml:"Severity"`
    Type        string `yaml:"Type"` // "suggest", "error", "warning"
}

// Issue represents a found issue in documentation
type Issue struct {
    File        string
    Line        int
    Column      int
    Rule        string
    Message     string
    Severity    string
    Suggestion  string
    OriginalText string
}

// Analyzer handles document analysis
type Analyzer struct {
    config *Config
    rules  []Rule
}

// NewAnalyzer creates a new analyzer instance
func NewAnalyzer(configPath string) (*Analyzer, error) {
    config, err := loadConfig(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load config: %w", err)
    }

    return &Analyzer{
        config: config,
        rules:  config.Rules,
    }, nil
}

// loadConfig loads configuration from YAML file
func loadConfig(configPath string) (*Config, error) {
    if configPath == "" {
        return getDefaultConfig(), nil
    }

    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    return &config, nil
}

// getDefaultConfig returns default AI optimization rules
func getDefaultConfig() *Config {
    return &Config{
        StylesPath:   "./styles",
        MinWordCount: 10,
        Formats: map[string]Format{
            "markdown": {
                Extensions: []string{".md", ".markdown"},
                Parser:     "markdown",
            },
            "html": {
                Extensions: []string{".html", ".htm"},
                Parser:     "html",
            },
        },
        Rules: []Rule{
            {
                Name:        "contextual-dependency",
                Description: "Detect sections that depend on previous context",
                Pattern:     `(?i)\b(this|that|these|those|above|below|previously|earlier)\b(?:\s+\w+){0,3}\s+(?:will|should|must|can|may)`,
                Severity:    "warning",
                Type:        "suggest",
            },
            {
                Name:        "semantic-discoverability",
                Description: "Ensure product names are included in relevant sections",
                Pattern:     `^##+\s+(?:Configure|Setup|Install|Enable)\s+\w+(?:\s+\w+)*$`,
                Severity:    "suggestion",
                Type:        "suggest",
            },
            {
                Name:        "implicit-knowledge",
                Description: "Detect assumed knowledge without explanation",
                Pattern:     `(?i)\b(?:simply|just|obviously|clearly|of course|naturally)\b`,
                Severity:    "warning",
                Type:        "suggest",
            },
            {
                Name:        "visual-dependency",
                Description: "Detect references to visual elements without text alternatives",
                Pattern:     `(?i)(?:see\s+(?:the\s+)?(?:diagram|image|figure|chart|screenshot)|(?:above|below)\s+(?:image|diagram|figure))`,
                Severity:    "error",
                Type:        "error",
            },
            {
                Name:        "generic-headings",
                Description: "Detect generic headings that lack context",
                Pattern:     `^##+\s+(?:Overview|Introduction|Getting Started|Configuration|Setup|Installation)$`,
                Severity:    "suggestion",
                Type:        "suggest",
            },
            {
                Name:        "incomplete-context",
                Description: "Detect incomplete procedural instructions",
                Pattern:     `(?i)^(?:\d+\.\s*|[-*]\s*)?(?:configure|set up|enable|disable|update|modify)\s+\w+(?:\s+\w+)*\.?\s*$`,
                Severity:    "warning",
                Type:        "suggest",
            },
        },
    }
}

// AnalyzeFile analyzes a single file for AI optimization issues
func (a *Analyzer) AnalyzeFile(filePath string) ([]Issue, error) {
    content, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    return a.analyzeContent(filePath, string(content)), nil
}

// analyzeContent analyzes content string for issues
func (a *Analyzer) analyzeContent(filePath, content string) []Issue {
    var issues []Issue
    lines := strings.Split(content, "\n")

    for i, line := range lines {
        lineNum := i + 1
        issues = append(issues, a.analyzeLine(filePath, line, lineNum)...)
    }

    // Additional content-level analysis
    issues = append(issues, a.analyzeStructure(filePath, content)...)

    return issues
}

// analyzeLine analyzes a single line for issues
func (a *Analyzer) analyzeLine(filePath, line string, lineNum int) []Issue {
    var issues []Issue

    for _, rule := range a.rules {
        regex, err := regexp.Compile(rule.Pattern)
        if err != nil {
            continue
        }

        matches := regex.FindAllStringSubmatchIndex(line, -1)
        for _, match := range matches {
            if len(match) >= 2 {
                matchText := line[match[0]:match[1]]
                issue := Issue{
                    File:         filePath,
                    Line:         lineNum,
                    Column:       match[0] + 1,
                    Rule:         rule.Name,
                    Message:      a.generateMessage(rule, matchText),
                    Severity:     rule.Severity,
                    Suggestion:   a.generateSuggestion(rule, matchText, line),
                    OriginalText: matchText,
                }
                issues = append(issues, issue)
            }
        }
    }

    return issues
}

// analyzeStructure performs document-level structural analysis
func (a *Analyzer) analyzeStructure(filePath, content string) []Issue {
    var issues []Issue

    // Check for missing product context in headings
    headingRegex := regexp.MustCompile(`(?m)^#{1,6}\s+(.+)$`)
    headings := headingRegex.FindAllStringSubmatch(content, -1)

    productNames := a.extractProductNames(content)
    
    for _, heading := range headings {
        if len(heading) > 1 {
            headingText := heading[1]
            if a.isGenericHeading(headingText) && !a.containsProductContext(headingText, productNames) {
                issues = append(issues, Issue{
                    File:     filePath,
                    Line:     a.findLineNumber(content, heading[0]),
                    Rule:     "missing-product-context",
                    Message:  "Heading lacks product-specific context",
                    Severity: "suggestion",
                    Suggestion: fmt.Sprintf("Consider adding product name: '%s %s'", 
                        a.inferProductName(productNames), headingText),
                    OriginalText: headingText,
                })
            }
        }
    }

    return issues
}

// generateMessage creates a human-readable message for the issue
func (a *Analyzer) generateMessage(rule Rule, matchText string) string {
    switch rule.Name {
    case "contextual-dependency":
        return "This text may depend on previous context. Consider making it self-contained."
    case "semantic-discoverability":
        return "Consider including product name for better AI discoverability."
    case "implicit-knowledge":
        return "Avoid assuming user knowledge. Provide explicit context."
    case "visual-dependency":
        return "Visual reference detected. Provide text alternative."
    case "generic-headings":
        return "Generic heading detected. Add specific context."
    case "incomplete-context":
        return "Instruction may lack sufficient context. Include prerequisites and specific steps."
    default:
        return rule.Description
    }
}

// generateSuggestion creates improvement suggestions
func (a *Analyzer) generateSuggestion(rule Rule, matchText, fullLine string) string {
    switch rule.Name {
    case "contextual-dependency":
        return "Replace contextual references with specific details"
    case "implicit-knowledge":
        return "Replace assumption words with explicit explanations"
    case "visual-dependency":
        return "Add text description alongside visual reference"
    case "generic-headings":
        return "Add product/feature name to heading"
    case "incomplete-context":
        return "Include prerequisite steps and specific system/location details"
    default:
        return "Consider rewriting for AI clarity"
    }
}

// Helper functions
func (a *Analyzer) extractProductNames(content string) []string {
    // Simple heuristic to find potential product names
    // Look for capitalized words that appear frequently
    words := regexp.MustCompile(`\b[A-Z][a-zA-Z]+\b`).FindAllString(content, -1)
    frequency := make(map[string]int)
    
    for _, word := range words {
        if len(word) > 3 && !a.isCommonWord(word) {
            frequency[word]++
        }
    }

    var products []string
    for word, count := range frequency {
        if count >= 3 { // Appears at least 3 times
            products = append(products, word)
        }
    }

    return products
}

func (a *Analyzer) isCommonWord(word string) bool {
    commonWords := []string{"The", "This", "That", "With", "From", "Your", "When", "Where", "What", "How"}
    for _, common := range commonWords {
        if word == common {
            return true
        }
    }
    return false
}

func (a *Analyzer) isGenericHeading(heading string) bool {
    generic := []string{"overview", "introduction", "getting started", "configuration", "setup", "installation"}
    lower := strings.ToLower(heading)
    for _, g := range generic {
        if strings.Contains(lower, g) {
            return true
        }
    }
    return false
}

func (a *Analyzer) containsProductContext(heading string, products []string) bool {
    lower := strings.ToLower(heading)
    for _, product := range products {
        if strings.Contains(lower, strings.ToLower(product)) {
            return true
        }
    }
    return false
}

func (a *Analyzer) inferProductName(products []string) string {
    if len(products) > 0 {
        return products[0] // Return most frequent
    }
    return "[PRODUCT_NAME]"
}

func (a *Analyzer) findLineNumber(content, target string) int {
    lines := strings.Split(content, "\n")
    for i, line := range lines {
        if strings.Contains(line, target) {
            return i + 1
        }
    }
    return 1
}

// Output formatting
func printIssues(issues []Issue, format string) {
    switch format {
    case "json":
        printJSONIssues(issues)
    case "sarif":
        printSARIFIssues(issues)
    default:
        printStandardIssues(issues)
    }
}

func printStandardIssues(issues []Issue) {
    for _, issue := range issues {
        severity := strings.ToUpper(issue.Severity)
        fmt.Printf("%s:%d:%d: %s [%s] %s\n",
            issue.File, issue.Line, issue.Column, severity, issue.Rule, issue.Message)
        
        if issue.Suggestion != "" {
            fmt.Printf("    Suggestion: %s\n", issue.Suggestion)
        }
        fmt.Println()
    }
}

func printJSONIssues(issues []Issue) {
    fmt.Println("JSON output not implemented yet")
}

func printSARIFIssues(issues []Issue) {
    fmt.Println("SARIF output not implemented yet")
}

// CLI interface
func main() {
    var (
        configPath = flag.String("config", "", "Path to configuration file")
        outputFormat = flag.String("output", "standard", "Output format (standard, json, sarif)")
        fix = flag.Bool("fix", false, "Attempt to automatically fix issues")
        recursive = flag.Bool("recursive", false, "Process directories recursively")
    )
    flag.Parse()

    if len(flag.Args()) == 0 {
        fmt.Fprintf(os.Stderr, "Usage: %s [options] <file_or_directory>\n", os.Args[0])
        flag.PrintDefaults()
        os.Exit(1)
    }

    analyzer, err := NewAnalyzer(*configPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error creating analyzer: %v\n", err)
        os.Exit(1)
    }

    var allIssues []Issue

    for _, path := range flag.Args() {
        issues, err := processPath(analyzer, path, *recursive)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", path, err)
            continue
        }
        allIssues = append(allIssues, issues...)
    }

    if *fix {
        fmt.Println("Auto-fix functionality not yet implemented")
    }

    printIssues(allIssues, *outputFormat)

    if len(allIssues) > 0 {
        os.Exit(1)
    }
}

func processPath(analyzer *Analyzer, path string, recursive bool) ([]Issue, error) {
    var allIssues []Issue

    stat, err := os.Stat(path)
    if err != nil {
        return nil, err
    }

    if stat.IsDir() {
        if recursive {
            err = filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
                if err != nil {
                    return err
                }

                if !d.IsDir() && isSupportedFile(filePath) {
                    issues, err := analyzer.AnalyzeFile(filePath)
                    if err != nil {
                        fmt.Fprintf(os.Stderr, "Warning: failed to analyze %s: %v\n", filePath, err)
                        return nil
                    }
                    allIssues = append(allIssues, issues...)
                }
                return nil
            })
        } else {
            entries, err := os.ReadDir(path)
            if err != nil {
                return nil, err
            }

            for _, entry := range entries {
                if !entry.IsDir() {
                    filePath := filepath.Join(path, entry.Name())
                    if isSupportedFile(filePath) {
                        issues, err := analyzer.AnalyzeFile(filePath)
                        if err != nil {
                            fmt.Fprintf(os.Stderr, "Warning: failed to analyze %s: %v\n", filePath, err)
                            continue
                        }
                        allIssues = append(allIssues, issues...)
                    }
                }
            }
        }
    } else {
        if isSupportedFile(path) {
            issues, err := analyzer.AnalyzeFile(path)
            if err != nil {
                return nil, err
            }
            allIssues = append(allIssues, issues...)
        }
    }

    return allIssues, err
}

func isSupportedFile(path string) bool {
    ext := strings.ToLower(filepath.Ext(path))
    supportedExts := []string{".md", ".markdown", ".html", ".htm", ".txt", ".rst"}
    
    for _, supported := range supportedExts {
        if ext == supported {
            return true
        }
    }
    return false
}

// Additional utility functions for advanced analysis could be added here:
// - Chunk boundary analysis
// - Semantic similarity detection
// - Error message pattern recognition
// - Table structure validation
// - Link and reference validation