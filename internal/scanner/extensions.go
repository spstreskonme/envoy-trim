package scanner

// supportedExtensions defines the set of file extensions that the scanner
// will inspect when searching for environment variable references.
// Keys are extensions (including the leading dot); values are a human-readable
// description used in debug output.
var supportedExtensions = map[string]string{
	".go":    "Go source",
	".js":    "JavaScript",
	".jsx":   "JavaScript (JSX)",
	".ts":    "TypeScript",
	".tsx":   "TypeScript (JSX)",
	".py":    "Python",
	".rb":    "Ruby",
	".sh":    "Shell script",
	".bash":  "Bash script",
	".zsh":   "Zsh script",
	".env":   "dotenv file",
	".yaml":  "YAML",
	".yml":   "YAML",
	".toml":  "TOML",
	".json":  "JSON",
	".tf":    "Terraform",
	".hcl":   "HCL",
	".dockerfile": "Dockerfile",
	".conf":  "configuration file",
	".ini":   "INI configuration",
	".php":   "PHP",
	".java":  "Java",
	".kt":    "Kotlin",
	".rs":    "Rust",
	".cs":    "C#",
	".cpp":   "C++",
	".c":     "C",
	".swift": "Swift",
	".ex":    "Elixir",
	".exs":   "Elixir script",
	".lua":   "Lua",
	".r":     "R",
	".scala": "Scala",
	".groovy": "Groovy",
	".pl":    "Perl",
	".pm":    "Perl module",
}

// skippedDirs contains directory names that should never be traversed during
// scanning. These are typically dependency or generated-artifact directories
// that would produce noisy false-positives and slow down the scan.
var skippedDirs = map[string]bool{
	"node_modules":  true,
	".git":          true,
	".hg":           true,
	".svn":          true,
	"vendor":        true,
	"dist":          true,
	"build":         true,
	".next":         true,
	".nuxt":         true,
	".turbo":        true,
	"__pycache__":   true,
	".pytest_cache": true,
	".mypy_cache":   true,
	".venv":         true,
	"venv":          true,
	"env":           true,
	".tox":          true,
	"target":        true,
	".gradle":       true,
	"bin":           true,
	"obj":           true,
	"coverage":      true,
	".nyc_output":   true,
	"tmp":           true,
	"temp":          true,
	".cache":        true,
}

// IsSupported reports whether the given file extension (e.g. ".go") is
// included in the scanner's supported extension list.
func IsSupported(ext string) bool {
	_, ok := supportedExtensions[ext]
	return ok
}

// IsSkippedDir reports whether the given directory name should be excluded
// from recursive scanning.
func IsSkippedDir(name string) bool {
	return skippedDirs[name]
}
