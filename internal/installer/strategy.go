package installer

import (
	"github.com/hubble-works/pim/internal/config"
)

// Strategy defines the interface for different installation strategies.
//
// Each strategy must implement methods to prepare the output, add files,
// and close any resources when done.
type Strategy interface {
	Prepare() error
	AddFile(srcPath, relativePath string) error
	Close() error
}

func createStrategy(
	strategyType config.Strategy,
	outputPath string,
) Strategy {
	switch strategyType {
	case config.StrategyConcat:
		return NewConcatStrategy(outputPath)
	case config.StrategyFlatten:
		return NewFlattenStrategy(outputPath)
	case config.StrategyPreserve:
		return NewPreserveStrategy(outputPath)
	default:
		return NewFlattenStrategy(outputPath)
	}
}
