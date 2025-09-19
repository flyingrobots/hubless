# hubless

Home for the off-platform roadmap tooling that keeps GitMind's Features Ledger fresh.

## Usage

1. Clone this repo beside `git-mind` (e.g., `/path/to/git/hubless` and `/path/to/git/git-mind`).
2. Run the updater:

   ```bash
   python3 update_progress.py --root ../git-mind
   # or rely on auto-detection if the repos are siblings
   python3 update_progress.py
   ```

3. Commit the regenerated Markdown in `git-mind`.

Set `GITMIND_ROOT` or pass `--root` if your checkout lives elsewhere.

## Makefile integration

Inside `git-mind`, `make features-update` shells out to this script. Point the
`HUBLESS_REPO` (or `HUBLESS_PATH`) environment variable at your hubless clone if
it isn't located at `../hubless`.
