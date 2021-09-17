package main

import (
	"flag"
	"fmt"
	"os"
)

// delete all localy pushed branches
func deletePushed() {
	//needs to execute git fetch
}

// delete all local branches
func deleteAll() {
	branchesPath := ".git/refs/heads"
	files, err := os.ReadDir(branchesPath)

	if err != nil {
		print("You are not in a git repository...\n")
	}

	//not tested
	for _, file := range files {
		filePath := fmt.Sprint("%s/%s", branchesPath, file.Name())
		if file.Name() != "master" && os.Remove(filePath) == nil {
			print(fmt.Sprint("Cold not remove %s\n", file.Name()))
		}
	}

}

func main() {
	print("git-extras...\n")
	var toDeletePushed, toDeleteAll bool
	deletePushedFlagName, deleteAllFlagName := "delete-pushed", "delete-all"

	flag.BoolVar(&toDeletePushed, deletePushedFlagName, false, "delete local branches that are available on origin")
	flag.BoolVar(&toDeleteAll, deleteAllFlagName, false, "delete all local branches")
	flag.Parse()

	if toDeleteAll && toDeletePushed {
		print(fmt.Sprintf("Can not provide both %s and %s flags.", deletePushedFlagName, deleteAllFlagName))
		return
	}

	deletePushed()
	//if toDeleteAll {
	//	deleteAll()
	//} else if toDeletePushed {
	//	deletePushed()
	//}

}
