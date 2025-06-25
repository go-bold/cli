package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
)

var (
	port      string
	hotReload bool
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the development server",
	Long: `Start the Go Bold development server.

This command starts your application in development mode with:
- Hot reload (if enabled)
- Debug logging
- Environment variable loading
- Database connection testing`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func startServer() {
	if !isGoBoldProject() {
		fmt.Println("❌ This is not a Go Bold project directory")
		fmt.Println("Run 'bold new myapp' to create a new project")
		return
	}

	fmt.Printf("🚀 Starting Go Bold development server on port %s...\n", port)

	if hotReload {
		fmt.Println("🔥 Hot reload enabled")
		startWithHotReload()
	} else {
		startNormal()
	}
}

func isGoBoldProject() bool {
	_, err1 := os.Stat("go.mod")
	_, err2 := os.Stat("main.go")
	return err1 == nil && err2 == nil
}

func startNormal() {
	cmd := exec.Command("go", "run", "main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Server error: %v", err)
	}
}

func startWithHotReload() {
	// For hot reload, we'll use air if available
	if _, err := exec.LookPath("air"); err != nil {
		fmt.Println("⚠️  Hot reload requires 'air' to be installed")
		fmt.Println("Install with: go install github.com/air-verse/air@latest")
		fmt.Println("Starting without hot reload...")
		startNormal()
		return
	}

	// Create .air.toml if it doesn't exist
	createAirConfig()

	cmd := exec.Command("air")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Hot reload error: %v", err)
	}
}

func createAirConfig() {
	airConfig := `# .air.toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ."
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
`

	if _, err := os.Stat(".air.toml"); os.IsNotExist(err) {
		err := os.WriteFile(".air.toml", []byte(airConfig), 0644)
		if err != nil {
			fmt.Println("❌ Cannot create .air.toml")
			fmt.Println("Starting server without hot reload...")
			startNormal()
		}
	}
}

func init() {
	serveCmd.Flags().StringVarP(&port, "port", "p", "8000", "Port to run the server on")
	serveCmd.Flags().BoolVar(&hotReload, "hot", true, "Enable hot reload")

	rootCmd.AddCommand(serveCmd)
}
