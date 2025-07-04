package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

type BufConfig struct {
	Deps []string `yaml:"deps"`
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		processProtoFiles()
	}()

	go func() {
		defer wg.Done()
		generatePulsarCode()
	}()

	go func() {
		defer wg.Done()
		generatePythonDependencies()
	}()

	wg.Wait()
}

func processProtoFiles() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	fmt.Printf("Program executed from: %s\n", cwd)

	fmt.Println("Generating proto code")

	// Collect all .proto files
	var protoFiles []string
	if err = filepath.WalkDir("./proto", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".proto" {
			protoFiles = append(protoFiles, path)
		}
		return nil
	}); err != nil {
		log.Fatalf("Error walking directory: %v", err)
	}

	fmt.Println(protoFiles)

	var wg sync.WaitGroup

	// Process each proto file concurrently
	for _, file := range protoFiles {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()

			content, err := os.ReadFile(file)
			if err != nil {
				log.Printf("Failed to read file %s: %v", file, err)
				return
			}

			matched, err := regexp.Match(`(?i)option\s+go_package\s*=\s*".*Zenrock-Foundation.*"`, content)
			if err != nil {
				log.Printf("Failed to match regex in file %s: %v", file, err)
				return
			}

			if matched {
				fmt.Printf("Processing file %s\n", file)

				cmd1 := exec.Command("buf", "generate", "--template", "proto/buf.gen.gogo.yaml", file)
				cmd1.Stdout = os.Stdout
				cmd1.Stderr = os.Stderr
				if err := cmd1.Run(); err != nil {
					log.Printf("Failed to run buf generate for file %s: %v", file, err)
					return
				}

				cmd2 := exec.Command("buf", "generate", "--template", "proto/buf.gen.python.yaml", file)
				cmd2.Stdout = os.Stdout
				cmd2.Stderr = os.Stderr
				if err := cmd2.Run(); err != nil {
					log.Printf("Failed to run buf generate for file %s: %v", file, err)
					return
				}
			}
		}(file)
	}

	wg.Wait()
	fmt.Println("Proto files generated.")

	// Move proto files to the right places
	srcDir := filepath.Join("github.com", "Zenrock-Foundation", "zrchain", "v6")
	if err := copyDir(srcDir, "./"); err != nil {
		log.Fatalf("Failed to copy files: %v", err)
	}
	if err := os.RemoveAll("./github.com"); err != nil {
		log.Fatalf("Failed to remove github.com directory: %v", err)
	}
}

func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			// Create the directory in the destination
			return os.MkdirAll(dstPath, info.Mode())
		} else {
			// Copy the file
			return copyFile(path, dstPath)
		}
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.ReadFrom(in)
	if err != nil {
		return err
	}
	return out.Close()
}

func generatePythonDependencies() {
	// Read the buf.yaml file
	bufYamlPath := "proto/buf.yaml"
	bufYamlContent, err := os.ReadFile(bufYamlPath)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", bufYamlPath, err)
	}

	// Parse the YAML content
	var bufConfig BufConfig
	err = yaml.Unmarshal(bufYamlContent, &bufConfig)
	if err != nil {
		log.Fatalf("Failed to parse %s: %v", bufYamlPath, err)
	}

	var wg sync.WaitGroup

	// Generate Python dependencies concurrently
	for _, dep := range bufConfig.Deps {
		wg.Add(1)
		go func(dep string) {
			defer wg.Done()
			fmt.Printf("Generating python dependencies for %s\n", dep)
			cmd := exec.Command("buf", "generate", "--template", "proto/buf.gen.python.yaml", dep)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Printf("Failed to generate dependencies for %s: %v", dep, err)
			}
		}(dep)
	}

	wg.Wait()
	fmt.Println("Python dependencies generated.")
}

func generatePulsarCode() {
	// Check if protoc-gen-go-pulsar is installed
	if _, err := exec.LookPath("protoc-gen-go-pulsar"); err != nil {
		fmt.Println("Error: protoc-gen-go-pulsar is not installed. Please install it before proceeding.")
		fmt.Println("go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar")
		os.Exit(1)
	}

	// Get the project root directory
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to get project root directory: %v", err)
	}
	projectRootDir := strings.TrimSpace(string(output))

	fmt.Println("Generating proto pulsar code")

	// Run buf generate command for main proto files
	bufGenCmd := exec.Command("buf", "generate", "-v", "--template", filepath.Join(projectRootDir, "proto", "buf.gen.pulsar.yaml"))
	bufGenCmd.Dir = projectRootDir
	bufGenCmd.Stdout = os.Stdout
	bufGenCmd.Stderr = os.Stderr
	if err := bufGenCmd.Run(); err != nil {
		log.Fatalf("Failed to run buf generate: %v", err)
	}

	// Also generate pulsar code for sidecar proto files
	fmt.Println("Generating sidecar proto pulsar code")
	sidecarBufGenCmd := exec.Command("buf", "generate", "-v", "--template", filepath.Join(projectRootDir, "proto", "buf.gen.pulsar.yaml"), "buf.build/zenrock-foundation/sidecar")
	sidecarBufGenCmd.Dir = projectRootDir
	sidecarBufGenCmd.Stdout = os.Stdout
	sidecarBufGenCmd.Stderr = os.Stderr
	if err := sidecarBufGenCmd.Run(); err != nil {
		log.Printf("Warning: Failed to generate sidecar pulsar code: %v", err)
	}

	fmt.Println("Pulsar files generated.")
}
