package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os/exec"
	"strings"
)

func deleteBranches(onlyMerged bool) {
	branchesToDelete, err := getAllBranchesExceptCurrent()

	if err != nil {
		print("error...")
		return
	}

	deleteAlways := false
	var inputDelete string
	for _, branch := range branchesToDelete {
		if !deleteAlways {
			validInputRange := [4]string{"y", "Y", "n", "N"}

			for do := true; do && !deleteAlways; do = !contains(validInputRange[:], inputDelete) {
				fmt.Printf("Do you want to delete this branch: %q (y/n), if you want to repeat this action for all other branches press (Y/N): ", branch)
				fmt.Scanln(&inputDelete)
			}

			if inputDelete == "N" {
				fmt.Print("Exiting...")
				return
			} else if inputDelete == "n" {
				fmt.Printf("Skiping %q branch", branch)
				continue
			}

			deleteAlways = inputDelete == "Y"
		}
		deleteBranch(branch, onlyMerged)
	}
}

func deleteBranch(branchName string, onlyMerged bool) {
	var cmdRemove *exec.Cmd
	if onlyMerged {
		cmdRemove = exec.Command("git", "branch", "-d", branchName)
	} else {
		cmdRemove = exec.Command("git", "branch", "-D", branchName)
	}
	_ = cmdRemove.Run()
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func removeCurrentBranch(branches []string) ([]string, error) {
	indexToRemove := -1
	for i, val := range branches {
		if val == "*" {
			indexToRemove = i
		}
	}

	if indexToRemove == -1 {
		return make([]string, 0), errors.New("pattern not found")
	}

	indexToRemoveAlso := indexToRemove + 2
	return append(branches[0:indexToRemove], branches[indexToRemoveAlso:]...), nil
}

func getAllBranchesExceptCurrent() ([]string, error) {
	cmd := exec.Command("git", "branch")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return make([]string, 0), err
	}
	return removeCurrentBranch(strings.Fields(out.String()))
}

func main() {
	var toDeleteMerged, toDeleteAll bool
	deleteMergedFlagName, deleteAllFlagName := "delete-merged", "delete-all"

	flag.BoolVar(&toDeleteMerged, deleteMergedFlagName, false, "delete local branches that have been pushed and merged")
	flag.BoolVar(&toDeleteAll, deleteAllFlagName, false, "delete all local branches")
	flag.Parse()

	if toDeleteAll && toDeleteMerged {
		print(fmt.Sprintf("Can not provide both %s and %s flags.", deleteMergedFlagName, deleteAllFlagName))
		return
	}

	deleteBranches(toDeleteMerged)
}
