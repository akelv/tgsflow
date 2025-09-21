package cmd

import (
	"archive/zip"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kelvin/tgsflow/src/core/config"
	"github.com/kelvin/tgsflow/src/core/thoughts"
	"github.com/kelvin/tgsflow/src/templates"
	"github.com/kelvin/tgsflow/src/util/logx"
	"github.com/spf13/cobra"
)

// CmdInit implements `tgs init --decorate [--ci-template github|gitlab|none]`.
func CmdInit(args []string) int {
	fs := flag.NewFlagSet("tgs init", flag.ContinueOnError)
	decorate := fs.Bool("decorate", true, "Create tgs/ layout and optional CI templates")
	ciTemplate := fs.String("ci-template", "github", "CI template: github|gitlab|none")
	interactive := fs.Bool("interactive", false, "Interactive config setup")
	templatesSrc := fs.String("templates", "", "Path, URL, or git repo for tgs templates (overrides embedded)")
	templatesRef := fs.String("templates-ref", "", "Git ref (branch/tag/commit) when --templates is a git repo")
	templatesSubdir := fs.String("templates-subdir", "", "Subdirectory inside the templates source that contains tgs templates")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	logx.Infof("initializing TGS (decorate=%v, ci-template=%s, interactive=%v, templates=%s, templates-ref=%s, templates-subdir=%s)", *decorate, *ciTemplate, *interactive, *templatesSrc, *templatesRef, *templatesSubdir)
	if !*decorate {
		logx.Infof("Nothing to do (decorate=false)")
		return 0
	}

	// Mirror tgs templates into ./tgs (idempotent)
	// Use config template data for placeholders (e.g., tgs.yml.tmpl)
	tmplData := config.DefaultTemplateData("")
	if *interactive {
		tmplData = config.PromptInteractive(os.Stdin, os.Stdout, tmplData)
	}
	if err := thoughts.EnsureDir("tgs"); err != nil {
		logx.Errorf("failed to create tgs dir: %v", err)
		return 1
	}
	// Choose template source: remote/local dir or embedded
	if *templatesSrc != "" {
		// Three cases: archive URL, git repo, or local directory
		if isArchiveURL(*templatesSrc) {
			tmpDir, err := downloadAndExtractTemplates(*templatesSrc)
			if err != nil {
				logx.Errorf("download templates: %v", err)
				return 1
			}
			defer os.RemoveAll(tmpDir)
			root, err := findTemplatesRoot(tmpDir, *templatesSubdir)
			if err != nil {
				logx.Errorf("locate templates root (archive): %v", err)
				return 1
			}
			if err := templates.RenderTGSTreeFromDir(root, ".", tmplData); err != nil {
				logx.Errorf("render tgs templates (archive): %v", err)
				return 1
			}
		} else if isLikelyGitURL(*templatesSrc) {
			tmpDir, err := cloneGitRepo(*templatesSrc, *templatesRef)
			if err != nil {
				logx.Errorf("clone templates repo: %v", err)
				return 1
			}
			defer os.RemoveAll(tmpDir)
			root, err := findTemplatesRoot(tmpDir, *templatesSubdir)
			if err != nil {
				logx.Errorf("locate templates root (git): %v", err)
				return 1
			}
			if err := templates.RenderTGSTreeFromDir(root, ".", tmplData); err != nil {
				logx.Errorf("render tgs templates (git): %v", err)
				return 1
			}
		} else {
			root, err := findTemplatesRoot(*templatesSrc, *templatesSubdir)
			if err != nil {
				logx.Errorf("locate templates root (local): %v", err)
				return 1
			}
			if err := templates.RenderTGSTreeFromDir(root, ".", tmplData); err != nil {
				logx.Errorf("render tgs templates (local): %v", err)
				return 1
			}
		}
	} else if err := templates.RenderTGSTree(".", tmplData); err != nil {
		logx.Errorf("render tgs templates: %v", err)
		return 1
	}

	// Optional CI templates
	switch strings.ToLower(*ciTemplate) {
	case "github":
		wfDir := filepath.Join(".github", "workflows")
		if err := thoughts.EnsureDir(wfDir); err != nil {
			logx.Errorf("failed to ensure workflows dir: %v", err)
			return 1
		}
		approve := filepath.Join(wfDir, "tgs-approve.yml")
		content, err := templates.Render("ci/github-approve.yml.tmpl", nil)
		if err != nil {
			logx.Errorf("render workflow: %v", err)
			return 1
		}
		_, err = thoughts.EnsureFile(approve, []byte(content))
		if err != nil {
			logx.Errorf("failed to write workflow: %v", err)
			return 1
		}
		logx.Infof("ensured %s", approve)
	case "gitlab":
		// Minimal stub .gitlab-ci.yml
		content, err := templates.Render("ci/gitlab-ci.yml.tmpl", nil)
		if err != nil {
			logx.Errorf("render gitlab ci: %v", err)
			return 1
		}
		_, err = thoughts.EnsureFile(".gitlab-ci.yml", []byte(content))
		if err != nil && !errors.Is(err, os.ErrExist) {
			logx.Errorf("failed to ensure gitlab ci: %v", err)
			return 1
		}
	case "none":
		// no-op
		logx.Infof("skipping CI templates (ci-template=none)")
	default:
		fmt.Fprintln(os.Stderr, "unknown --ci-template value; expected github|gitlab|none")
		return 2
	}

	// tgs.yml is rendered above via templates; keep message for parity
	logx.Infof("ensured tgs/ tree and tgs.yml (idempotent)")
	logx.Infof("tgs init complete (idempotent)")
	return 0
}

// Cobra command constructor colocated for cleanliness
func newInitCommand() *cobra.Command {
	var (
		flagDecorate        bool
		flagCITemplate      string
		flagInteractive     bool
		flagTemplates       string
		flagTemplatesRef    string
		flagTemplatesSubdir string
	)
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Initialize TGS layout (idempotent)",
		Long:    "Create the standard TGS scaffolding and ensure tgs/tgs.yml exists. Optionally seed CI templates and select a remote or local templates source.",
		Example: "  tgs init\n  tgs init --ci-template gitlab\n  tgs init --interactive\n  tgs init --decorate=false\n  tgs init --templates /path/to/templates/tgs\n  tgs init --templates https://example.com/org-tgs-templates.zip\n  tgs init --templates https://github.com/org/repo.git --templates-ref main --templates-subdir path/in/repo",
		RunE: func(c *cobra.Command, args []string) error {
			// Reconstruct args for CmdInit which uses the stdlib flag parser to keep behavior consistent with tests.
			forward := []string{}
			if c.Flags().Changed("decorate") {
				if flagDecorate {
					forward = append(forward, "--decorate")
				} else {
					forward = append(forward, "--decorate=false")
				}
			}
			if c.Flags().Changed("ci-template") {
				forward = append(forward, "--ci-template", flagCITemplate)
			}
			if c.Flags().Changed("interactive") {
				if flagInteractive {
					forward = append(forward, "--interactive")
				} else {
					forward = append(forward, "--interactive=false")
				}
			}
			if c.Flags().Changed("templates") {
				forward = append(forward, "--templates", flagTemplates)
			}
			if c.Flags().Changed("templates-ref") {
				forward = append(forward, "--templates-ref", flagTemplatesRef)
			}
			if c.Flags().Changed("templates-subdir") {
				forward = append(forward, "--templates-subdir", flagTemplatesSubdir)
			}
			code := CmdInit(forward)
			return codeToErr(code)
		},
	}
	cmd.Flags().BoolVar(&flagDecorate, "decorate", true, "Create tgs/ layout and optional CI templates")
	cmd.Flags().StringVar(&flagCITemplate, "ci-template", "github", "CI template: github|gitlab|none")
	cmd.Flags().BoolVar(&flagInteractive, "interactive", false, "Interactive config setup")
	cmd.Flags().StringVar(&flagTemplates, "templates", "", "Path, URL (.zip/.tar.gz), or git repo for tgs templates")
	cmd.Flags().StringVar(&flagTemplatesRef, "templates-ref", "", "Git ref (branch/tag/commit) when --templates is a git repo")
	cmd.Flags().StringVar(&flagTemplatesSubdir, "templates-subdir", "", "Subdirectory containing tgs templates in source (default auto-detect)")
	return cmd
}

// downloadAndExtractTemplates downloads a URL pointing to either a .zip or .tar.gz archive
// and extracts it to a temporary directory, returning the root directory containing the tgs
// templates. It assumes the archive contains the tgs structure at its root.
func downloadAndExtractTemplates(url string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "tgs-templates-*")
	if err != nil {
		return "", err
	}
	// pick extension
	lower := strings.ToLower(url)
	switch {
	case strings.HasSuffix(lower, ".zip"):
		if err := downloadZip(url, tmpDir); err != nil {
			os.RemoveAll(tmpDir)
			return "", err
		}
	case strings.HasSuffix(lower, ".tar.gz") || strings.HasSuffix(lower, ".tgz"):
		if err := downloadTarGz(url, tmpDir); err != nil {
			os.RemoveAll(tmpDir)
			return "", err
		}
	default:
		// treat as a directory-like URL not supported
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("unsupported templates URL (expect .zip or .tar.gz): %s", url)
	}
	return tmpDir, nil
}

func httpGetToFile(url string, dst string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("download failed: %s", resp.Status)
	}
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.Copy(f, resp.Body); err != nil {
		return err
	}
	return nil
}

func downloadZip(url string, extractDir string) error {
	// Use the system unzip if available to avoid pulling extra deps
	zipFile := filepath.Join(extractDir, "templates.zip")
	if err := httpGetToFile(url, zipFile); err != nil {
		return err
	}
	// Prefer native unzip via `tar -xf` on some platforms doesn't work for zip
	// Try `unzip -q` if present
	if _, err := exec.LookPath("unzip"); err == nil {
		return execCommand("unzip", []string{"-q", zipFile, "-d", extractDir})
	}
	// Fallback to Go stdlib zip reader
	return unzipFile(zipFile, extractDir)
}

func downloadTarGz(url string, extractDir string) error {
	tarGz := filepath.Join(extractDir, "templates.tar.gz")
	if err := httpGetToFile(url, tarGz); err != nil {
		return err
	}
	// Use system tar if possible
	if _, err := exec.LookPath("tar"); err == nil {
		return execCommand("tar", []string{"-xzf", tarGz, "-C", extractDir})
	}
	return fmt.Errorf("tar not found in PATH; cannot extract %s", tarGz)
}

// isArchiveURL returns true if the URL ends with a supported archive extension.
func isArchiveURL(u string) bool {
	lower := strings.ToLower(u)
	return strings.HasSuffix(lower, ".zip") || strings.HasSuffix(lower, ".tar.gz") || strings.HasSuffix(lower, ".tgz")
}

// isLikelyGitURL heuristically detects git repo URLs.
func isLikelyGitURL(s string) bool {
	return strings.HasSuffix(strings.ToLower(s), ".git") || strings.HasPrefix(s, "git@") || strings.Contains(s, "github.com/") || strings.Contains(s, "gitlab.com/") || strings.Contains(s, "bitbucket.org/")
}

// cloneGitRepo clones a git repository at ref (optional) into a temporary directory and returns the path.
func cloneGitRepo(repoURL string, ref string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "tgs-templates-git-*")
	if err != nil {
		return "", err
	}
	args := []string{"clone", "--depth", "1"}
	if ref != "" {
		args = append(args, "--branch", ref)
	}
	args = append(args, repoURL, tmpDir)
	if err := execCommand("git", args); err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}
	return tmpDir, nil
}

// findTemplatesRoot returns the directory that directly contains the TGS templates (i.e., design/, agentops/, README.md.tmpl, tgs.yml.tmpl)
// If subdir is provided, it returns repoRoot/subdir. Otherwise, it auto-detects by checking for expected files/dirs.
func findTemplatesRoot(repoRoot string, subdir string) (string, error) {
	root := repoRoot
	if subdir != "" {
		root = filepath.Join(repoRoot, subdir)
	} else {
		// Try common locations
		candidates := []string{
			repoRoot,
			filepath.Join(repoRoot, "tgs"),
			filepath.Join(repoRoot, "templates", "data", "tgs"),
		}
		found := ""
		for _, c := range candidates {
			if dirHasTgsTemplates(c) {
				found = c
				break
			}
		}
		if found == "" {
			return "", fmt.Errorf("could not auto-detect templates root under %s; provide --templates-subdir", repoRoot)
		}
		root = found
	}
	// Final validation
	if !dirHasTgsTemplates(root) {
		return "", fmt.Errorf("templates root %s does not appear to contain TGS templates", root)
	}
	return root, nil
}

func dirHasTgsTemplates(dir string) bool {
	// Look for design/ and tgs.yml.tmpl at minimum
	if st, err := os.Stat(filepath.Join(dir, "design")); err != nil || !st.IsDir() {
		return false
	}
	if _, err := os.Stat(filepath.Join(dir, "tgs.yml.tmpl")); err != nil {
		return false
	}
	return true
}

func execCommand(name string, args []string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func unzipFile(zipPath string, destDir string) error {
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer zr.Close()
	for _, f := range zr.File {
		outPath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(outPath, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		// Use file mode if available, else default
		mode := f.Mode()
		if mode == 0 {
			mode = 0o644
		}
		w, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
		if err != nil {
			return err
		}
		copyErr := func() error {
			defer w.Close()
			_, err := io.Copy(w, rc)
			return err
		}()
		cerr := rc.Close()
		if copyErr != nil {
			return copyErr
		}
		if cerr != nil {
			return cerr
		}
	}
	return nil
}
