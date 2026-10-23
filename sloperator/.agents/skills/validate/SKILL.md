---
name: validate
description: Run Sloperator's tests, formatting check, and linter
---

# Validate

From the project root, run:

1. `uv run pytest`
2. `uv run ruff format --check`
3. `uv run ruff check`

Summarize which checks passed and which failed. Report failures before proposing fixes.
