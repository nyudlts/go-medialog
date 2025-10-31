//go:build !ignore
// +build !ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var (
	OS      string
	Version string
	bd      string
)

func init() {
	flag.StringVar(&OS, "os", "", "Operating System")
	flag.StringVar(&Version, "version", "", "Version")
}

func main() {
	flag.Parse()

	if !(OS == "windows" || OS == "linux" || OS == "darwin") {
		fmt.Println("Error: Unsupported OS. Please use 'windows', 'linux', or 'darwin'.")
		os.Exit(1)
	}

	if Version == "" {
		fmt.Println("Error: version is required") // will grab the version from the source if this is empty later
		os.Exit(1)
	}

	bd = filepath.Join(fmt.Sprintf("go-medialog-%s-%s", OS, Version))

	if _, err := os.Stat(bd); err == nil {
		if err := os.RemoveAll(bd); err != nil {
			fmt.Printf("ERROR could not remove build directory: %s", bd)
			os.Exit(1)
		}
		fmt.Printf("removed existing build directory: %s\n", bd)
	}

	if err := os.Mkdir(bd, 0755); err != nil {
		fmt.Printf("ERROR could not create build directory: %s", bd)
		os.Exit(1)
	}
	fmt.Printf("created build directory: %s\n", bd)

	var buildCommand *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		buildCommand = exec.Command("powershell", "-File", "build.ps1", "-Os", OS)
	case "linux":
		buildCommand = exec.Command("./build.sh")
	default:
		fmt.Println("builds systems other than windows or linux are not supported")
		os.Exit(1)
	}

	//build the binary
	buildCommand.Stdout = os.Stdout
	buildCommand.Stderr = os.Stderr
	if err := buildCommand.Run(); err != nil {
		fmt.Printf("failed to execute build script: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("binary built")

	//move the binary to the build directory
	if err := os.Rename("medialog", filepath.Join(bd, "medialog")); err != nil {
		fmt.Printf("ERROR could not move binary to build directory: %v", err)
		os.Exit(1)
	}

	//copy the needed directories
	public := os.DirFS("../public")
	templates := os.DirFS("../templates")
	files := os.DirFS("../files")
	if err := os.CopyFS(filepath.Join(bd, "public"), public); err != nil {
		fmt.Printf("ERROR could not copy public directory: %v", err)
		os.Exit(1)
	}

	if err := os.CopyFS(filepath.Join(bd, "templates"), templates); err != nil {
		fmt.Printf("ERROR could not copy templates directory: %v", err)
		os.Exit(1)
	}

	if err := os.CopyFS(filepath.Join(bd, "files"), files); err != nil {
		fmt.Printf("ERROR could not copy files directory: %v", err)
		os.Exit(1)
	}

	//copy the Makfile
	mf, err := os.ReadFile("Makefile")
	if err != nil {
		fmt.Printf("ERROR could not read Makefile: %v", err)
		os.Exit(1)
	}

	if err := os.WriteFile(filepath.Join(bd, "Makefile"), mf, 0755); err != nil {
		fmt.Printf("ERROR could not copy Makefile: %v", err)
		os.Exit(1)
	}
	fmt.Println("moved resources to build directory")

	//compress the build directory
	tgzFile := fmt.Sprintf("%s.tar.gz", bd)
	var compressCmd = exec.Command("tar", "-czvf", tgzFile, bd)
	compressCmd.Stdout = os.Stdout
	compressCmd.Stderr = os.Stderr
	if err := compressCmd.Run(); err != nil {
		fmt.Printf("failed to compress build directory: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("build directory compressed")

	//remove the build directory
	if err := os.RemoveAll(bd); err != nil {
		fmt.Printf("ERROR could not remove build directory: %v", err)
		os.Exit(1)
	}
	fmt.Println("build directory removed")

	if _, err := os.Stat(filepath.Join("../bin", tgzFile)); err == nil {
		if err := os.Remove(filepath.Join("../bin", tgzFile)); err != nil {
			fmt.Printf("ERROR could not remove existing tarball: %v", err)
		}
		fmt.Println("removed existing tarball from bin directory")
	}

	//move the build tar to bin dir
	if err := os.Rename(tgzFile, filepath.Join("../bin", tgzFile)); err != nil {
		fmt.Printf("ERROR could not move build tar to bin directory: %v", err)
		os.Exit(1)
	}
	fmt.Println("build tarball moved to bin directory")
}
