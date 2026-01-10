package types

type Language string
type ContextAlgorithmType string

const (
	GOLANG     Language = "golang"
	PYTHON     Language = "python"
	JAVASCRIPT Language = "javascript"
	UNKNOWN    Language = "unknown"
)

type NodeImport struct {
	ImportPackage string
	FilePath      string
}

// Prompting types
const (
	BASE_LLM_PROMPT = `

	The below is the context of the repository. The paths are given along with the signatures of the files. If you want more context for a specific file, ask the user to use this command to add more files in this format:
	/add <paths>

	For example:
	/add /a/b/c.go /b/c/d.go

	`
)

// Context Algorithm types
const (
	PAGERANK_CONTEXT_ALGORITHM ContextAlgorithmType = "pagerank"

	DEFAULT_CONTEXT_ALGORITHM ContextAlgorithmType = PAGERANK_CONTEXT_ALGORITHM
)
