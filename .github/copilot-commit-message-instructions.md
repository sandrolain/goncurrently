Conventional Commits (English, ≤50-char subject)

- Format: type(scope?): subject
- Types: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert
- Scope: optional, lowercase, describes affected area
- Subject: imperative, lowercase, English, ≤50 chars, no period, no emojis
- Body (optional): wrap at 72 chars, explain what and why
- Breaking changes: use type!: subject and add footer "BREAKING CHANGE: details"
- Footer (optional): "Fixes #123", "Refs #456", co-authors, etc.
- Keep commits atomic; one logical change per commit

Examples:

- feat(api): add retry policy
- fix(ui): handle null event state
- docs(readme): add setup guide
