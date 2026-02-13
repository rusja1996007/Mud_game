package main

import (
	"go.uber.org/zap"
)

func foo(log *zap.Logger) {
	log.Error("some error")

}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger.Debug("heelo its DEBUG")

	foo(logger)

}

//11 10 56
