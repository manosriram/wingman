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

	The above is the context of the repository. Answer all the questions in a simple and understandable way. If you need any other files for more context, do ask.

	`
)

// Context Algorithm types
const (
	PAGERANK_CONTEXT_ALGORITHM ContextAlgorithmType = "pagerank"

	DEFAULT_CONTEXT_ALGORITHM ContextAlgorithmType = PAGERANK_CONTEXT_ALGORITHM
)
