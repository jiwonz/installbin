package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func cloneFile(sourceFile, destinationDir string) error {
	// Open the source file for reading
	src, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer src.Close()

	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return err
	}

	// Create the destination file with the same name as the source file
	destinationFile := filepath.Join(destinationDir, filepath.Base(sourceFile))
	dst, err := os.Create(destinationFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}

func isExecutable(file string) (bool, error) {
	info, err := os.Stat(file)
	if err != nil {
		return false, err
	}

	// Check if the file has the execute permission bit set
	mode := info.Mode()
	isExecutable := (mode & 0111) != 0

	return isExecutable, nil
}

func ls() {
	cmd := exec.Command("sudo", "ls", "/usr/local/bin/")

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the output
	fmt.Println("installed bins in /usr/local/bin/:")
	fmt.Println(string(output))
}

func main() {
	args := os.Args
	if len(args) < 2 {
		executableName := filepath.Base(args[0])
		fmt.Printf("Usage: %s <executable file or ls>\n", executableName)
		return
	}
	exe := args[1]
	isExe, err := isExecutable(exe)
	if isExe == false || err != nil {
		if exe == "ls" {
			ls()
			return
		}
		fmt.Println("Error: given file is not executable!")
	}
	if isExe == true {
		cmd := exec.Command("sudo", "cp", exe, "/usr/local/bin/")
		err = cmd.Run()
		if err != nil {
			fmt.Println("Error:", err)
		}

		msg := "installed"
		fileName := filepath.Base(exe)
		_, err := os.Stat(exe)
		if err == nil {
			// File exists
			//fmt.Println("File exists.")
			msg = "updated"
		} else if os.IsNotExist(err) {
			// File does not exist
			//fmt.Println("File does not exist.")
		} else {
			// An error occurred (e.g., permission denied)
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("installbin: Successfully %s %s in /usr/local/bin/\n", msg, fileName)
	}
}
