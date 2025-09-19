package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	groupRe = regexp.MustCompile("(?is)<!-- group-progress:([a-z0-9\\-]+):begin -->\n```text\n(.*?)\n```\n<!-- group-progress:\\1:end -->")

	readmeBlockRe     = regexp.MustCompile("(?s)<!-- features-progress:begin -->\n```text\n(.*?)\n```\n<!-- features-progress:end -->")
	overallSectionRe  = regexp.MustCompile("(?s)<!-- progress-overall:begin -->\n```text\n(.*?)\n```\n<!-- progress-overall:end -->")
	numberRe          = regexp.MustCompile("-?\\d+(?:\\.\\d+)?")
	threePlusNewlines = regexp.MustCompile("\\n{3,}")

	milestonePatterns = map[string]*regexp.Regexp{
		"mvp":   regexp.MustCompile("(?s)<!-- progress-mvp:begin -->\n```text\n(.*?)\n```\n<!-- progress-mvp:end -->"),
		"alpha": regexp.MustCompile("(?s)<!-- progress-alpha:begin -->\n```text\n(.*?)\n```\n<!-- progress-alpha:end -->"),
		"beta":  regexp.MustCompile("(?s)<!-- progress-beta:begin -->\n```text\n(.*?)\n```\n<!-- progress-beta:end -->"),
		"v1":    regexp.MustCompile("(?s)<!-- progress-v1:begin -->\n```text\n(.*?)\n```\n<!-- progress-v1:end -->"),
	}

	taskBlockPatterns = map[string]*regexp.Regexp{
		"mvp":   regexp.MustCompile("(?s)> <!-- tasks-mvp:begin -->\n(.*?)\n> <!-- tasks-mvp:end -->"),
		"alpha": regexp.MustCompile("(?s)> <!-- tasks-alpha:begin -->\n(.*?)\n> <!-- tasks-alpha:end -->"),
		"beta":  regexp.MustCompile("(?s)> <!-- tasks-beta:begin -->\n(.*?)\n> <!-- tasks-beta:end -->"),
		"v1":    regexp.MustCompile("(?s)> <!-- tasks-v1:begin -->\n(.*?)\n> <!-- tasks-v1:end -->"),
	}

	milestoneLabels = map[string]string{
		"MVP":    "mvp",
		"Alpha":  "alpha",
		"Beta":   "beta",
		"v1.0.0": "v1",
		"V1":     "v1",
	}

	milestoneWeights = map[string]float64{
		"mvp":   0.3,
		"alpha": 0.3,
		"beta":  0.2,
		"v1":    0.2,
	}
)

type featureEntry struct {
	milestone string
	percent   float64
	weight    float64
}

type totalWeight struct {
	total  float64
	weight float64
}

func main() {
	rootFlag := flag.String("root", "", "Path to the git-mind repository root (defaults to sibling repo detection).")
	flag.Parse()

	root, err := resolveRoot(*rootFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ledgerPath := filepath.Join(root, "docs", "features", "Features_Ledger.md")
	readmePath := filepath.Join(root, "README.md")

	overall, hasOverall, _, _, err := updateLedger(ledgerPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := updateReadme(readmePath, overall, hasOverall); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func resolveRoot(cliRoot string) (string, error) {
	var candidates []string
	appendCandidate := func(path string) {
		if path == "" {
			return
		}
		abs, err := filepath.Abs(path)
		if err == nil {
			candidates = append(candidates, abs)
		} else {
			candidates = append(candidates, path)
		}
	}

	if cliRoot != "" {
		appendCandidate(cliRoot)
	}
	if env := os.Getenv("GITMIND_ROOT"); env != "" {
		appendCandidate(env)
	}

	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		appendCandidate(filepath.Join(exeDir, "git-mind"))
		appendCandidate(filepath.Join(filepath.Dir(exeDir), "git-mind"))
		appendCandidate(exeDir)
	}

	if cwd, err := os.Getwd(); err == nil {
		appendCandidate(filepath.Join(cwd, "git-mind"))
		appendCandidate(filepath.Join(filepath.Dir(cwd), "git-mind"))
		appendCandidate(cwd)
	}

	seen := make(map[string]struct{})
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if _, ok := seen[candidate]; ok {
			continue
		}
		seen[candidate] = struct{}{}
		ledger := filepath.Join(candidate, "docs", "features", "Features_Ledger.md")
		if _, err := os.Stat(ledger); err == nil {
			return candidate, nil
		}
	}

	return "", errors.New("unable to locate git-mind repo. Pass --root /path/to/git-mind or set GITMIND_ROOT environment variable")
}

func updateLedger(ledgerPath string) (float64, bool, map[string]float64, map[string][]string, error) {
	data, err := os.ReadFile(ledgerPath)
	if err != nil {
		return 0, false, nil, nil, err
	}

	original := string(data)
	updated := original

	featureEntries := make([]featureEntry, 0)

	matches := groupRe.FindAllStringSubmatchIndex(original, -1)
	for _, loc := range matches {
		if len(loc) < 4 {
			continue
		}
		groupName := original[loc[2]:loc[3]]
		pct, weight, count := computeGroupPercent(original, loc[1], &featureEntries)
		bar := progressBar(pct)
		featureTag := fmt.Sprintf("features=%d", count)
		block := fmt.Sprintf("<!-- group-progress:%s:begin -->\n```text\n%s\n------------|-------------|------------|\n           MVP          Alpha    v1.0.0 \n%s\n```\n<!-- group-progress:%s:end -->", groupName, bar, featureTag, groupName)
		segment := original[loc[0]:loc[1]]
		updated = strings.Replace(updated, segment, block, 1)
		_ = weight
	}

	totals := make(map[string]totalWeight)
	for key := range milestonePatterns {
		totals[key] = totalWeight{}
	}
	for _, entry := range featureEntries {
		key, ok := resolveMilestoneKey(entry.milestone)
		if !ok {
			continue
		}
		tw := totals[key]
		tw.total += entry.percent * entry.weight
		tw.weight += entry.weight
		totals[key] = tw
	}

	milestoneProgress := make(map[string]float64)
	for key, tw := range totals {
		if tw.weight > 0 {
			milestoneProgress[key] = tw.total / tw.weight
		} else {
			milestoneProgress[key] = 0
		}
	}

	order := []string{"mvp", "alpha", "beta", "v1"}
	var prev *float64
	for _, key := range order {
		current := milestoneProgress[key]
		if prev == nil {
			prev = new(float64)
			*prev = current
			continue
		}
		if current > *prev {
			milestoneProgress[key] = *prev
		} else {
			milestoneProgress[key] = current
		}
		*prev = milestoneProgress[key]
	}

	totalWeight := 0.0
	weightedSum := 0.0
	for key, w := range milestoneWeights {
		totalWeight += w
		weightedSum += milestoneProgress[key] * w
	}

	hasOverall := totalWeight > 0
	overall := 0.0
	if hasOverall && totalWeight > 0 {
		overall = weightedSum / totalWeight
	}

	if loc := overallSectionRe.FindStringIndex(updated); loc != nil {
		legend := []string{
			fmt.Sprintf("MVP %d%%", int(math.Round(milestoneProgress["mvp"]*100))),
			fmt.Sprintf("Alpha %d%%", int(math.Round(milestoneProgress["alpha"]*100))),
			fmt.Sprintf("Beta %d%%", int(math.Round(milestoneProgress["beta"]*100))),
			fmt.Sprintf("v1.0.0 %d%%", int(math.Round(milestoneProgress["v1"]*100))),
		}
		block := fmt.Sprintf("<!-- progress-overall:begin -->\n```text\n%s\n%s\n```\n<!-- progress-overall:end -->", progressBar(overall), strings.Join(legend, " | "))
		updated = updated[:loc[0]] + block + updated[loc[1]:]
	}

	for key, pattern := range milestonePatterns {
		if match := pattern.FindStringSubmatch(updated); len(match) > 0 {
			block := fmt.Sprintf("<!-- progress-%s:begin -->\n```text\n%s\n```\n<!-- progress-%s:end -->", key, progressBar(milestoneProgress[key]), key)
			updated = strings.Replace(updated, match[0], block, 1)
		}
	}

	tasks := extractTasks(updated)
	for key, pattern := range taskBlockPatterns {
		match := pattern.FindStringSubmatch(updated)
		if len(match) == 0 {
			continue
		}
		items := tasks[key]
		var body string
		if len(items) == 0 {
			body = "> _All tracked tasks complete_"
		} else {
			prefixed := make([]string, 0, len(items))
			for _, item := range items {
				trimmed := strings.TrimSpace(item)
				prefixed = append(prefixed, "> "+trimmed)
			}
			body = strings.Join(prefixed, "\n")
		}
		replacement := fmt.Sprintf("> <!-- tasks-%s:begin -->\n%s\n> <!-- tasks-%s:end -->", key, body, key)
		updated = strings.Replace(updated, match[0], replacement, 1)
	}

	if updated != original {
		if err := os.WriteFile(ledgerPath, []byte(updated), 0644); err != nil {
			return 0, false, nil, nil, err
		}
	}

	return overall, hasOverall, milestoneProgress, tasks, nil
}

func updateReadme(readmePath string, overall float64, hasOverall bool) error {
	data, err := os.ReadFile(readmePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	original := string(data)
	if !strings.Contains(original, "## 📊 Status") {
		return nil
	}

	updated := original

	if !readmeBlockRe.MatchString(updated) {
		parts := strings.SplitN(updated, "## 📊 Status", 2)
		if len(parts) == 2 {
			progress := "\n<!-- features-progress:begin -->\n```text\nFeature progress to be updated via hubless/update_progress.py\n```\n<!-- features-progress:end -->\n\n"
			updated = parts[0] + "## 📊 Status\n\n" + progress + parts[1]
		}
	}

	if hasOverall {
		block := fmt.Sprintf("<!-- features-progress:begin -->\n```text\n%s\n```\n<!-- features-progress:end -->\n", progressBar(overall))
		updated = readmeBlockRe.ReplaceAllString(updated, "")
		parts := strings.SplitN(updated, "## 📊 Status", 2)
		if len(parts) == 2 {
			updated = parts[0] + "## 📊 Status\n\n" + block + "\n" + parts[1]
		}
	}

	updated = threePlusNewlines.ReplaceAllString(updated, "\n\n")

	if updated != original {
		return os.WriteFile(readmePath, []byte(updated), 0644)
	}
	return nil
}

func computeGroupPercent(md string, anchor int, entries *[]featureEntry) (float64, float64, int) {
	if anchor >= len(md) {
		return 0, 0, 0
	}
	tail := md[anchor:]
	idx := strings.Index(tail, "| Emoji")
	if idx == -1 {
		return 0, 0, 0
	}
	table := tail[idx:]
	lines := strings.Split(table, "\n")
	if len(lines) < 3 {
		return 0, 0, 0
	}
	header := strings.TrimSpace(lines[0])
	rows := lines[2:]

	headerCells := splitTableRow(header)

	findIdx := func(name string) int {
		lower := strings.ToLower(name)
		for i, cell := range headerCells {
			if strings.Contains(strings.ToLower(cell), lower) {
				return i
			}
		}
		return -1
	}

	progressIdx := findIdx("progress")
	klocIdx := findIdx("kloc")
	milestoneIdx := findIdx("milestone")

	var values []float64
	var weights []float64
	count := 0

	for _, rawLine := range rows {
		line := strings.TrimRight(rawLine, "\r")
		if !strings.HasPrefix(line, "|") {
			break
		}
		cells := splitTableRow(line)
		prog, ok := parseProgress(cells, progressIdx)
		if !ok {
			continue
		}
		weight := parseWeight(cells, klocIdx)
		values = append(values, float64(prog))
		weights = append(weights, weight)
		count++

		if entries != nil {
			milestone := "Unassigned"
			if milestoneIdx >= 0 && milestoneIdx < len(cells) {
				if cells[milestoneIdx] != "" {
					milestone = cells[milestoneIdx]
				}
			}
			*entries = append(*entries, featureEntry{milestone: strings.TrimSpace(milestone), percent: float64(prog) / 100.0, weight: weight})
		}
	}

	if len(values) == 0 {
		return 0, 0, 0
	}

	weightSum := 0.0
	weighted := 0.0
	for i, val := range values {
		weightSum += weights[i]
		weighted += val * weights[i]
	}

	pct := 0.0
	if weightSum > 0 {
		pct = (weighted / weightSum) / 100.0
	} else {
		pct = average(values) / 100.0
	}

	return pct, weightSum, count
}

func splitTableRow(row string) []string {
	trimmed := strings.TrimSpace(row)
	parts := strings.Split(trimmed, "|")
	cells := make([]string, 0, len(parts))
	for _, part := range parts {
		cells = append(cells, strings.TrimSpace(part))
	}
	return cells
}

func parseProgress(cells []string, progressIdx int) (int, bool) {
	if progressIdx >= 0 && progressIdx < len(cells) {
		if prog, ok := extractProgress(cells[progressIdx]); ok {
			return prog, true
		}
	}
	for _, cell := range cells {
		if prog, ok := extractProgress(cell); ok {
			return prog, true
		}
	}
	return 0, false
}

func extractProgress(cell string) (int, bool) {
	trimmed := strings.TrimSpace(cell)
	if trimmed == "" {
		return 0, false
	}
	if strings.HasSuffix(trimmed, "%") {
		core := strings.TrimSpace(trimmed[:len(trimmed)-1])
		if isDigits(core) {
			val, err := strconv.Atoi(core)
			if err == nil {
				return val, true
			}
		}
	}
	digits := extractDigits(trimmed)
	if digits == "" {
		return 0, false
	}
	val, err := strconv.Atoi(digits)
	if err != nil {
		return 0, false
	}
	return val, true
}

func parseWeight(cells []string, klocIdx int) float64 {
	weight := 0.1
	if klocIdx >= 0 && klocIdx < len(cells) {
		raw := strings.ReplaceAll(cells[klocIdx], ",", "")
		match := numberRe.FindString(raw)
		if match != "" {
			if v, err := strconv.ParseFloat(match, 64); err == nil && v > 0 {
				weight = v
			}
		}
	}
	return weight
}

func average(vals []float64) float64 {
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	if len(vals) == 0 {
		return 0
	}
	return sum / float64(len(vals))
}

func extractDigits(s string) string {
	var builder strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func resolveMilestoneKey(label string) (string, bool) {
	trimmed := strings.TrimSpace(label)
	for raw, key := range milestoneLabels {
		if strings.EqualFold(trimmed, raw) {
			return key, true
		}
	}
	return "", false
}

func extractTasks(md string) map[string][]string {
	tasks := map[string][]string{
		"mvp":   {},
		"alpha": {},
		"beta":  {},
		"v1":    {},
	}
	lower := strings.ToLower(md)
	idx := strings.Index(lower, "## tasklist")
	if idx == -1 {
		return tasks
	}
	segment := md[idx:]

	lines := strings.Split(segment, "\n")
	for _, line := range lines {
		stripped := strings.TrimSpace(line)
		if !strings.HasPrefix(stripped, "- [ ]") {
			continue
		}
		tag, entry := parseTask(stripped)
		if tag == "" {
			tag = "mvp"
		}
		tasks[tag] = append(tasks[tag], entry)
	}
	return tasks
}

func parseTask(line string) (string, string) {
	rest := strings.TrimSpace(line[len("- [ ]"):])

	resolveFromToken := func(token string) string {
		token = strings.TrimSpace(token)
		if token == "" {
			return ""
		}
		parts := []string{token}
		if strings.Contains(token, ".") {
			parts = append(parts, strings.Split(token, ".")...)
		}
		for _, part := range parts {
			if key, ok := resolveMilestoneKey(part); ok {
				return key
			}
		}
		return ""
	}

	if strings.HasPrefix(rest, "[") && strings.Contains(rest, "]") {
		tagName := rest[1:strings.Index(rest, "]")]
		if key := resolveFromToken(tagName); key != "" {
			remainder := strings.TrimSpace(rest[strings.Index(rest, "]")+1:])
			if remainder == "" {
				return key, "- [ ]"
			}
			return key, "- [ ] " + remainder
		}
	}

	if strings.HasPrefix(rest, "(") && strings.Contains(rest, ")") {
		tagName := rest[1:strings.Index(rest, ")")]
		if key := resolveFromToken(tagName); key != "" {
			remainder := strings.TrimSpace(rest[strings.Index(rest, ")")+1:])
			if remainder == "" {
				return key, "- [ ]"
			}
			return key, "- [ ] " + remainder
		}
	}

	remainder := strings.TrimSpace(rest)
	if remainder == "" {
		return "", "- [ ]"
	}
	return "", "- [ ] " + remainder
}

func progressBar(pct float64) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}
	width := 40
	filled := int(math.Round(pct * float64(width)))
	if filled > width {
		filled = width
	}
	remainder := pct*float64(width) - float64(filled)
	edge := 0
	if filled < width && remainder > 0 {
		edge = 1
	}
	if filled+edge > width {
		edge = width - filled
	}
	bar := strings.Repeat("█", filled)
	if edge > 0 {
		bar += "▓"
	}
	remainderCount := width - filled - edge
	if remainderCount < 0 {
		remainderCount = 0
	}
	bar += strings.Repeat("░", remainderCount)
	bar += fmt.Sprintf(" %d%%", int(math.Round(pct*100)))
	return bar
}
