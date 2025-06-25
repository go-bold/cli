package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	templateType   string
	templateRepo   string
	templateBranch string
)

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new Go Bold application",
	Long: `Create a new Go Bold application with all the necessary boilerplate.

This command creates a new directory with your project name and sets up
a complete Go Bold application`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		return createProject(projectName)
	},
}

func createProject(name string) error {
	fmt.Printf("🚀 Creating Go Bold application '%s'...\n\n", name)

	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", name)
	}

	fmt.Printf("📦 Downloading template '%s' from %s...\n", templateType, templateRepo)

	if err := downloadTemplateFromAPI(name, templateType); err != nil {
		return fmt.Errorf("failed to download template: %w", err)
	}

	if err := processTemplateVariables(name); err != nil {
		fmt.Printf("⚠️  Warning: Could not process some template variables: %v\n", err)
	}

	fmt.Printf("✅  Successfully created Go Bold application '%s'\n\n", name)
	fmt.Printf("Get started:\n")
	fmt.Printf("  cd %s\n", name)
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  bold serve\n\n")
	fmt.Printf("Happy coding! 🎉\n")

	return nil
}

func downloadTemplateFromAPI(projectName, templateType string) error {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s?ref=%s",
		templateRepo, templateType, templateBranch)

	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return fmt.Errorf("template '%s' not found", templateType)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var files []GitHubFile
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return fmt.Errorf("failed to parse API response: %w", err)
	}

	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	return downloadFilesRecursively(files, projectName, "")
}

type GitHubFile struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
	URL         string `json:"url"`
}

func downloadFilesRecursively(files []GitHubFile, projectName, currentPath string) error {
	for _, file := range files {
		localPath := filepath.Join(projectName, currentPath, file.Name)

		if file.Type == "dir" {
			if err := os.MkdirAll(localPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", localPath, err)
			}

			resp, err := http.Get(file.URL)
			if err != nil {
				return fmt.Errorf("failed to get directory contents for %s: %w", file.Name, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				return fmt.Errorf("failed to fetch directory %s: HTTP %d", file.Name, resp.StatusCode)
			}

			var subFiles []GitHubFile
			if err := json.NewDecoder(resp.Body).Decode(&subFiles); err != nil {
				return fmt.Errorf("failed to parse directory contents for %s: %w", file.Name, err)
			}

			if err := downloadFilesRecursively(subFiles, projectName, filepath.Join(currentPath, file.Name)); err != nil {
				return err
			}
		} else {
			if err := downloadFile(file.DownloadURL, localPath); err != nil {
				return fmt.Errorf("failed to download %s: %w", file.Name, err)
			}
		}
	}
	return nil
}

func downloadFile(url, destPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP GET failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, resp.Body); err != nil {
		return fmt.Errorf("failed to write file content: %w", err)
	}

	return nil
}

func processTemplateVariables(projectName string) error {
	fmt.Printf("🔧 Processing template variables...\n")

	variables := map[string]string{
		"{{PROJECT_NAME}}": projectName,
		"{{MODULE_NAME}}":  projectName,
		"{{APP_NAME}}":     strings.Title(projectName),
		"{{PACKAGE_NAME}}": strings.ToLower(projectName),
	}

	return filepath.Walk(projectName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk path %s: %w", path, err)
		}

		if info.IsDir() || isBinaryFile(path) {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		contentStr := string(content)
		modified := false

		for placeholder, value := range variables {
			if strings.Contains(contentStr, placeholder) {
				contentStr = strings.ReplaceAll(contentStr, placeholder, value)
				modified = true
			}
		}

		if modified {
			if err := os.WriteFile(path, []byte(contentStr), info.Mode()); err != nil {
				return fmt.Errorf("failed to write file %s: %w", path, err)
			}
		}

		return nil
	})
}

func isBinaryFile(path string) bool {
	binaryExtensions := []string{
		".jpg", ".jpeg", ".png", ".gif", ".bmp",
		".pdf", ".zip", ".tar", ".gz",
		".exe", ".bin", ".so", ".dll",
	}

	ext := strings.ToLower(filepath.Ext(path))
	for _, binExt := range binaryExtensions {
		if ext == binExt {
			return true
		}
	}
	return false
}

func init() {
	newCmd.Flags().StringVarP(&templateType, "template", "t", "default", "Template type (default)")
	newCmd.Flags().StringVar(&templateRepo, "repo", "go-bold/templates", "Template repository")
	newCmd.Flags().StringVar(&templateBranch, "branch", "main", "Template repository branch")

	rootCmd.AddCommand(newCmd)
}
