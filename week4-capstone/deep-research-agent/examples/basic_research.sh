#!/bin/bash

# Basic Research Example
# This script demonstrates basic research using the Deep Research Agent

echo "====================================="
echo "Basic Research Example"
echo "====================================="
echo ""

# Example 1: Simple research with default settings
echo "Example 1: Basic research on Go concurrency patterns"
./research-agent research "Go concurrency patterns" \
  --output reports/go-concurrency.md
echo ""

# Example 2: Research with specific depth
echo "Example 2: Medium-depth research on AI ethics"
./research-agent research "AI ethics in healthcare" \
  --depth medium \
  --max-sources 15 \
  --output reports/ai-ethics.md
echo ""

# Example 3: Deep research with many sources
echo "Example 3: Deep research on quantum computing"
./research-agent research "quantum computing applications" \
  --depth deep \
  --max-sources 25 \
  --cite-style APA \
  --output reports/quantum-computing.md
echo ""

echo "Research complete! Check the reports/ directory for output."
