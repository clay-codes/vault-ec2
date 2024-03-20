package utils

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)
var keyPath string
// addAliasIfNotExists checks for the "vault01" alias in ~/.zshrc; if not found, it adds the alias using the current working directory to determine the key path.
// updateAlias ensures the "vault01" alias is always set to the correct command in ~/.zshrc.

func GetKeyPath() (string, error) {
    // Determine the key.pem path based on the current working directory
    cwd, err := os.Getwd()
    if err != nil {
        return "", fmt.Errorf("failed to get current working directory: %w", err)
    }
    keyPath = filepath.Join(cwd, "key.pem")
    return keyPath, nil
}

func UpdateAlias(ipv4dns string) error {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("failed to get user home directory: %w", err)
    }
    zshrcPath := fmt.Sprintf("%s/.zshrc", homeDir)
    file, err := os.Open(zshrcPath)
    if err != nil {
        return fmt.Errorf("failed to open .zshrc: %w", err)
    }
    defer file.Close()

    // Read the file and check if the alias exists
    var lines []string
    aliasCmd := fmt.Sprintf("alias vault01='ssh -i %s -o StrictHostKeyChecking=no ec2-user@%s'", keyPath, ipv4dns)
    aliasExists := false
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "alias vault01=") {
            aliasExists = true
            lines = append(lines, aliasCmd) // Replace existing alias
        } else {
            lines = append(lines, line)
        }
    }

    if err := scanner.Err(); err != nil {
        return fmt.Errorf("error reading .zshrc: %w", err)
    }

    if !aliasExists {
        lines = append(lines, aliasCmd) // Add alias if not found
    }

    // Rewrite the .zshrc file
    file, err = os.Create(zshrcPath)
    if err != nil {
        return fmt.Errorf("failed to open .zshrc for writing: %w", err)
    }
    defer file.Close()

    for _, line := range lines {
        if _, err := file.WriteString(line + "\n"); err != nil {
            return fmt.Errorf("failed to write to .zshrc: %w", err)
        }
    }

    return nil
}