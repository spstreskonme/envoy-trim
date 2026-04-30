// Package differ computes and formats the delta between two sets of
// environment variable keys — typically the original .env file and the
// pruned result after envoy-trim has removed unused entries.
//
// Usage:
//
//	base := envparser.Keys(originalEntries)
//	next := envparser.Keys(prunedEntries)
//	delta := differ.Compare(base, next)
//
//	// Human-readable output
//	differ.FormatText(os.Stdout, delta)
//
//	// JSON output
//	differ.FormatJSON(os.Stdout, delta)
//
//	// One-line summary
//	fmt.Println(differ.Summary(delta))
package differ
