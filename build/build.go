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

	"github.com/nyudlts/go-medialog/version"
)

var (
	OS    string
	all   bool
	clean bool
)

func init() {
	flag.StringVar(&OS, "os", "", "")
	flag.BoolVar(&all, "all", false, "")
	flag.BoolVar(&clean, "clean", false, "")
}

func main() {

	fmt.Println("go-medialog build system v0.1.0")
	flag.Parse()

	wd, _ := os.Getwd()
	if clean {
		fmt.Println("  * Cleaning build artifacts")

		binDirectory := filepath.Join(wd, "bin")
		artifacts, _ := os.ReadDir(binDirectory)
		for _, artifact := range artifacts {
			if err := os.RemoveAll(filepath.Join(binDirectory, artifact.Name())); err != nil {
				fmt.Printf("ERROR could not remove build artifact: %s", artifact.Name())
				os.Exit(1)
			}
		}

		fmt.Printf("  * removed artifacts from build directory: %s\n", binDirectory)
		return
	}

	if OS == "" {
		OS = runtime.GOOS
	}
	if !(OS == "windows" || OS == "linux" || OS == "darwin") {
		fmt.Println("Error: Unsupported OS. Please use 'windows', 'linux', or 'darwin'.")
		os.Exit(1)
	}

	buildDirectory := filepath.Join(wd, "build")

	if all {
		for _, target := range []string{"windows", "linux"} {
			binDirectory := filepath.Join(wd, "bin", fmt.Sprintf("go-medialog-%s-v%s", target, version.AppVersion))
			build(binDirectory, buildDirectory, wd, target)
		}
		return
	} else {

		binDirectory := filepath.Join(wd, "bin", fmt.Sprintf("go-medialog-%s-v%s", OS, version.AppVersion))
		build(binDirectory, buildDirectory, wd, OS)
	}
}

func build(binDirectory string, buildDirectory string, wd string, targetSystem string) {
	fmt.Printf("  * building for %s\n", targetSystem)
	if _, err := os.Stat(binDirectory); err == nil {
		if err := os.RemoveAll(binDirectory); err != nil {
			fmt.Printf("ERROR could not remove build directory: %s", binDirectory)
			os.Exit(1)
		}
		fmt.Printf("  * removed existing build directory: %s\n", binDirectory)
	}

	if err := os.Mkdir(binDirectory, 0755); err != nil {
		fmt.Printf("ERROR could not create build directory: %s", binDirectory)
		os.Exit(1)
	}
	fmt.Printf("  * created build directory: %s\n", binDirectory)

	var buildCommand *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		buildCommand = exec.Command("powershell", "-File", filepath.Join(buildDirectory, "build.ps1"), "-Os", targetSystem, "-Od", filepath.Join(binDirectory, "medialog"), "-Path", wd)
	case "linux":
		buildCommand = exec.Command("./build.sh", "--os", targetSystem, "--od", filepath.Join(binDirectory, "medialog"), "--path", wd)
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

	fmt.Println("  * binary built")

	//copy the needed directories
	public := os.DirFS("public")
	templates := os.DirFS("templates")
	files := os.DirFS("files")
	if err := os.CopyFS(filepath.Join(binDirectory, "public"), public); err != nil {
		fmt.Printf("ERROR could not copy public directory: %v", err)
		os.Exit(1)
	}

	if err := os.CopyFS(filepath.Join(binDirectory, "templates"), templates); err != nil {
		fmt.Printf("ERROR could not copy templates directory: %v", err)
		os.Exit(1)
	}

	if err := os.CopyFS(filepath.Join(binDirectory, "files"), files); err != nil {
		fmt.Printf("ERROR could not copy files directory: %v", err)
		os.Exit(1)
	}

	//copy the Makefile
	mf, err := os.ReadFile(filepath.Join(buildDirectory, "Makefile"))
	if err != nil {
		fmt.Printf("ERROR could not read Makefile: %v", err)
		os.Exit(1)
	}

	if err := os.WriteFile(filepath.Join(binDirectory, "Makefile"), mf, 0755); err != nil {
		fmt.Printf("ERROR could not copy Makefile: %v", err)
		os.Exit(1)
	}
	fmt.Println("  * moved resources to build directory")

	//compress the build directory

	tgzFile := fmt.Sprintf("%s.tar.gz", binDirectory)
	var compressCmd = exec.Command(
		"tar",
		"-czvf",
		tgzFile,
		"-C",
		filepath.Dir(binDirectory),
		filepath.Base(binDirectory),
	)
	compressCmd.Stdout = os.Stdout
	compressCmd.Stderr = os.Stderr
	if err := compressCmd.Run(); err != nil {
		fmt.Printf("failed to compress build directory: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("  * build directory compressed")

	if err := os.RemoveAll(binDirectory); err != nil {
		fmt.Printf("ERROR could not remove build directory: %s", binDirectory)
	}
	fmt.Printf("  * removed build directory: %s\n", binDirectory)

}
