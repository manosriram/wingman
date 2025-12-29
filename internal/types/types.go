package types

type Language string
type ContextAlgorithmType string

const (
	GO         Language = "golang"
	PYTHON     Language = "python"
	JAVASCRIPT Language = "javascript"
	UNKNOWN    Language = "unknown"
)

const (
	PAGERANK_CONTEXT_ALGORITHM ContextAlgorithmType = "pagerank"

	DEFAULT_CONTEXT_ALGORITHM ContextAlgorithmType = "pagerank"
)

type NodeImport struct {
	ImportPath string
	FilePath   string
}
