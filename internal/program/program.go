package program

type Program struct {
	TargetDir string
}

func NewProgram(targetDir string) *Program {
	return &Program{
		TargetDir: targetDir,
	}
}

func (*Program) Run() error {
	return nil
}
