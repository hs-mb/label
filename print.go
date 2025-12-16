package etikett

import (
	"os/exec"
	"strings"
)

func Print(printer string, data string, lprBinArg ...string) error {
	lprBin := "lpr"
	if len(lprBinArg) > 0 {
		lprBin = lprBinArg[0]
	}
	cmd := exec.Command(lprBin, "-P", printer, "-o", "raw")
	cmd.Stdin = strings.NewReader(data)
	return cmd.Run()
}
