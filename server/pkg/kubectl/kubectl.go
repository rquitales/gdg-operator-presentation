package kubectl

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetCRDSpec(crd string) ([]byte, error) {
	return kubectlExec(fmt.Sprintf("get crd %s -o json", crd))
}

func GetAllCRDs() ([]byte, error) {
	return kubectlExec("get crd -o json")
}

func kubectlExec(args string) ([]byte, error) {
	argsSplit := strings.Split(args, " ")
	cmd := exec.Command("kubectl", argsSplit...)
	return cmd.CombinedOutput()
}
