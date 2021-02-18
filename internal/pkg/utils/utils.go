package utils

import (
	"fmt"
	"github.com/gotomicro/egoctl/logger"
)

func DumpWrapper(msg string, dump func()) {
	fmt.Println()
	logger.Log.Infof("------------------------%s--------------------------", msg)
	dump()
	logger.Log.Infof("------------------------%s--------------------------\n\n", msg)
}
