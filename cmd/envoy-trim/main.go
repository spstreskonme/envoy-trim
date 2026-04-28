// main is the entry point for the envoy-trim CLI tool.
// It wires together the env parser, source scanner, and pruner
// to audit and optionally remove unused environment variables.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourorg/envoy-trim/internal/envparser"
	"github.com/yourorg/envoy-trim/internal/pruner"
	"github.com/yourorg/envoy-trim/internal/scanner"
)

const version = "0.1.0"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("envoy-trim", flag.ContinueOnError)

	var (
		envFile  = fs.String("env", ".env", "path to the .env file to audit")
		scanDir  = fs.String("dir", ".", "root directory to scan for variable references")
		dryRun   = fs.Bool("dry-run", false, "print unused variables without modifying the .env file")
		showVer  = fs.Bool("version", false, "print version and exit")
		verbose  = fs.Bool("verbose", false, "print all variables and their usage status")
	)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *showVer {
		fmt.Printf("envoy-trim %s\n", version)
		return nil
	}

	// Resolve paths to absolute so all downstream operations are consistent.
	absEnv, err := filepath.Abs(*envFile)
	if err != nil {
		return fmt.Errorf("resolving env file path: %w", err)
	}
	absDir, err := filepath.Abs(*scanDir)
	if err != nil {
		return fmt.Errorf("resolving scan directory path: %w", err)
	}

	// Parse the .env file and collect declared variable keys.
	envEntries, err := envparser.ParseFile(absEnv)
	if err != nil {
		return fmt.Errorf("parsing env file %q: %w", absEnv, err)
	}
	keys := envparser.Keys(envEntries)

	if len(keys) == 0 {
		fmt.Println("No variables found in", absEnv)
		return nil
	}

	// Scan the source tree for references to those keys.
	sc := scanner.New(absDir)
	refs, err := sc.ScanDir(keys)
	if err != nil {
		return fmt.Errorf("scanning directory %q: %w", absDir, err)
	}

	report := scanner.BuildReport(keys, refs)

	// Print results.
	if *verbose {
		fmt.Println("Variable usage report:")
		for _, key := range keys {
			status := "unused"
			if report.Used[key] {
				status = "used"
			}
			fmt.Printf("  %-40s %s\n", key, status)
		}
		fmt.Println()
	}

	if len(report.Unused) == 0 {
		fmt.Println("✓ No unused variables found.")
		return nil
	}

	fmt.Printf("Found %d unused variable(s):\n", len(report.Unused))
	for _, key := range report.Unused {
		fmt.Printf("  - %s\n", key)
	}

	if *dryRun {
		fmt.Println("\n(dry-run) No changes written.")
		return nil
	}

	// Prune unused variables from the .env file.
	removed, err := pruner.Prune(absEnv, report.Unused, false)
	if err != nil {
		return fmt.Errorf("pruning env file: %w", err)
	}

	fmt.Printf("\nRemoved %d variable(s) from %s\n", removed, absEnv)
	return nil
}
