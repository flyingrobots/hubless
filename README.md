# Hubless

> *Imagine GitHub… but in your repo. No hub; just Git.*

> [!INFO]
> **EARLY DAYS.** I just started this project yesterday. Expect rapid iteration, rough edges, and breaking changes.
> If you want boring stability, wait. If you want to see Git-native project state come alive, jump in now.

---

## Why Hubless?

> *Art doesn’t need to be explained; its purpose is to create a new reality as powerful and engaging as the one we live in.*

### **Hubless = Freedom**

Keep your entire project in your repo. No SaaS lock-in. Offline-first. Fast. Minimalist. Deeply integrated.
Web frontend optional (for PMs or when you’re on the go), but developers live in Git.

### The Vision

- **Git-native issues, boards, and execution.**
  Every change is a commit. No website required.

- **Conflict-free by design.**
  CRDT event streams, snapshots, and catalogs keep state boringly coherent.

- **Auditable forever.**
  Undo = append-only. No rewrites, no drift.

- **Sort of like Magit, but for project flow.**
  Fast TUI, consistent keystrokes, muscle-memory ergonomics.

- **Optional Play button.**
  Don’t just track issues — press ▶ to execute DAG-style tasks automatically.

- **Boring stuff just happens.**
  **Old way:**

  ```text
  ticket -> website -> bookkeeping -> copy-paste -> PR
  ```

  **New way:**

  ```bash
  git hubless start issue 34
  # branch created, issue assigned, kanban updated, draft PR opened
  ...
  git hubless submit issue 34
  # PR updated, undrafted, review requested
  ```

### Local CI dry run

Mirror the GitHub Actions workflow locally:

```bash
./scripts/ci-local.sh
```

This container runs formatting checks, lightweight linting (gofmt/goimports/revive), `go vet`, unit tests, docs regeneration, and verifies that all `![[…]]` placeholders are resolved.
