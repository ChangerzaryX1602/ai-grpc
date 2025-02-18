package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// copyFile copies the contents of the file named src to dst.
// It also attempts to preserve file permissions.
func CopyFile(src, dst string) error {
	// Open the source file for reading.
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %q: %w", src, err)
	}
	defer in.Close()

	// Create the destination file.
	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %q: %w", dst, err)
	}
	// Ensure the file is closed properly.
	defer func() {
		if cerr := out.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Copy file content.
	if _, err = io.Copy(out, in); err != nil {
		return fmt.Errorf("failed to copy contents from %q to %q: %w", src, dst, err)
	}

	// Retrieve the source file's mode (permissions).
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source file %q: %w", src, err)
	}

	// Set the same permissions on the destination file.
	if err = os.Chmod(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to set file mode on %q: %w", dst, err)
	}

	return nil
}

// copyDir recursively copies a directory tree, attempting to preserve permissions.
func CopyDir(src string, dst string) error {
	// Get properties of the source directory.
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source directory %q: %w", src, err)
	}

	// Create the destination directory with the same permissions.
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to create destination directory %q: %w", dst, err)
	}

	// Read entries in the source directory.
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read directory %q: %w", src, err)
	}

	// Iterate over each entry.
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectories.
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy files.
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}
func CreateDir(paths []string) {

	for _, path := range paths {
		// Extract the directory from the path.
		dir := filepath.Dir(path)
		// Create the directory (and parents) if it doesn't exist.
		// os.ModePerm uses 0777 permissions modified by the umask.
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Printf("Error creating directory %q: %v\n", dir, err)
			continue
		}

		// Create an empty file.
		// Note: os.Create truncates the file if it already exists.
		f, err := os.Create(path)
		if err != nil {
			fmt.Printf("Error creating file %q: %v\n", path, err)
			continue
		}
		f.Close() // Close the file handle.
	}

	fmt.Println("Project structure created successfully.")
}
