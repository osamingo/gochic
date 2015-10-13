package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func init() {

	err := exec.Command("which", "go").Run()
	if err != nil {
		fmt.Println("gochic: golang should be install")
		os.Exit(1)
	}

	err = exec.Command("which", "golint").Run()
	if err != nil {
		err = exec.Command("go", "get", "github.com/golang/lint/golint").Run()
		if err != nil {
			fmt.Println("gochic: failed to install golint")
			os.Exit(1)
		}
	}

	err = exec.Command("which", "goimports").Run()
	if err != nil {
		err = exec.Command("go", "get", "golang.org/x/tools/cmd/goimports").Run()
		if err != nil {
			fmt.Println("gochic: failed to install goimports")
			os.Exit(1)
		}
	}

}

func main() {

	if len(os.Args) < 2 {
		usage()
	}

	err := govet()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = golint()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = goimports()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func usage() {
	fmt.Println("Usage: gochic [packages]...")
	os.Exit(0)
}

func govet() (err error) {

	cmd := exec.Command("go", append([]string{"vet"}, os.Args[1:]...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Run: %s\n", cmd.Args)
		fmt.Printf("%s\n", out)
		err = errors.New("gochic: failed to go vet")
	}

	return
}

func golint() (err error) {

	cmd := exec.Command("golint", os.Args[1:]...)
	out, _ := cmd.CombinedOutput()
	if len(out) != 0 {
		fmt.Printf("Run: %s\n", cmd.Args)
		fmt.Printf("%s\n", out)
		err = errors.New("gochic: failed to golint")
	}

	return
}

func goimports() (err error) {

	args := []string{}
	for i := range os.Args {
		if strings.HasSuffix(os.Args[i], "/...") {
			out, err := exec.Command("go", "list", "-f='{{range .GoFiles}}{{.}} {{end}}'", os.Args[i]).CombinedOutput()
			if err != nil {
				return err
			}
			args = append(args, string(out[1:len(out)-3]))
		}
	}

	if len(args) == 0 {
		args = append(args, os.Args[1:]...)
	}

	cmd := exec.Command("goimports", append([]string{"-l=true"}, args...)...)
	out, _ := cmd.CombinedOutput()
	if len(out) != 0 {
		fmt.Printf("Run: %s\n", cmd.Args)
		fmt.Printf("%s\n", out)
		err = errors.New("gochic: failed to goimports")
	}

	return
}
