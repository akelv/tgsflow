package ears

import (
	"bufio"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type expectedCase struct {
	file          string
	expectShapes  map[int]Shape // line -> shape
	expectInvalid []int         // lines expected to fail parse
}

// isCandidate returns true if a trimmed line should be considered a requirement sentence.
func isCandidate(trimmed string) bool {
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return false
	}
	upper := strings.ToUpper(trimmed)
	return strings.HasPrefix(upper, "WHEN ") || strings.HasPrefix(upper, "WHILE ") || strings.HasPrefix(upper, "IF ") || strings.HasPrefix(upper, "THE ")
}

func Test_EARS_Fixtures_Matrix(t *testing.T) {
	_, thisFile, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(thisFile), "testdata")
	cases := []expectedCase{
		{file: "positive_ubiquitous.md", expectShapes: map[int]Shape{2: ShapeUbiquitous, 4: ShapeUbiquitous, 6: ShapeUbiquitous, 8: ShapeUbiquitous, 10: ShapeUbiquitous}},
		{file: "positive_event.md", expectShapes: map[int]Shape{2: ShapeEvent, 4: ShapeEvent, 6: ShapeEvent, 8: ShapeEvent}},
		{file: "positive_state.md", expectShapes: map[int]Shape{2: ShapeState, 4: ShapeState, 6: ShapeState, 8: ShapeState}},
		{file: "positive_complex.md", expectShapes: map[int]Shape{2: ShapeComplex, 4: ShapeComplex, 6: ShapeComplex}},
		{file: "positive_unwanted.md", expectShapes: map[int]Shape{2: ShapeUnwanted, 4: ShapeUnwanted, 6: ShapeUnwanted}},
		{file: "formatting_bullets_and_skip_blocks.md", expectShapes: map[int]Shape{2: ShapeEvent}},
		// Negatives
		{file: "negative_missing_system.md", expectInvalid: []int{2, 4}},
		{file: "negative_multiple_when.md", expectInvalid: []int{2, 4}},
		{file: "negative_wrong_order.md", expectInvalid: []int{2, 4}},
		{file: "negative_missing_shall.md", expectInvalid: []int{2, 4, 6}},
		// Optional 'Where' is out of scope; ensure we ignore non-candidates and assert nothing
		// Ambiguous phrases are syntactically valid ubiquitous; we assert shapes
		{file: "negative_ambiguous_phrases.md", expectShapes: map[int]Shape{2: ShapeUbiquitous, 4: ShapeEvent, 6: ShapeState}},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			path := filepath.Join(root, tc.file)
			f, err := os.Open(path)
			if err != nil {
				t.Fatalf("open: %v", err)
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
			inFence := false
			line := 0
			seen := make(map[int]bool)
			for scanner.Scan() {
				line++
				raw := scanner.Text()
				trimmed := strings.TrimSpace(raw)
				if strings.HasPrefix(trimmed, "```") {
					inFence = !inFence
					continue
				}
				if inFence {
					continue
				}
				if !isCandidate(trimmed) {
					continue
				}

				// Parse
				res, err := ParseRequirement(trimmed)
				if wantShape, ok := tc.expectShapes[line]; ok {
					seen[line] = true
					if err != nil {
						t.Fatalf("line %d: unexpected err: %v", line, err)
					}
					if res.Shape != wantShape {
						t.Fatalf("line %d: want %s, got %s", line, wantShape, res.Shape)
					}
				}
				// Expected invalid lines
				for _, bad := range tc.expectInvalid {
					if bad == line {
						seen[line] = true
						if err == nil {
							t.Fatalf("line %d: expected parse error, got shape=%s", line, res.Shape)
						}
					}
				}
			}
			if err := scanner.Err(); err != nil {
				t.Fatalf("scan: %v", err)
			}

			// Ensure all expectations were hit
			for ln := range tc.expectShapes {
				if !seen[ln] {
					t.Fatalf("missing assertion for line %d (shape)", ln)
				}
			}
			for _, ln := range tc.expectInvalid {
				if !seen[ln] {
					t.Fatalf("missing assertion for line %d (invalid)", ln)
				}
			}
		})
	}
}
