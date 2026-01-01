package language

import "github.com/manosriram/wingman/internal/types"

func GetStrategy(args StrategyArgs) LangStrategy {
	switch args.StrategyLanguage {
	case types.GOLANG:
		return NewGolangStrategy(args)
	}
	return NewDefaultStrategy(args)
}
