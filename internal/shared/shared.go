package shared

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/inconshreveable/mousetrap"
	"gopkg.in/yaml.v3"
	//"github.com/eiannone/keyboard"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/configShared"
)

//func Uniq[T comparable](elements []T) []T {
//	encountered := make(map[T]struct{})
//	var result []T
//
//	for _, v := range elements {
//		if _, found := encountered[v]; !found {
//			encountered[v] = struct{}{}
//			result = append(result, v)
//		}
//	}
//
//	return result
//}

func Exit(code int) {
	if mousetrap.StartedByExplorer() {
		fmt.Println("Press Enter to exit...")
		_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	}
	os.Exit(code)
}

func SetYamlComment(node *yaml.Node, path []string, comment string) bool {
	if node == nil || len(path) == 0 {
		return false
	}

	if len(path) == 1 {
		if node.Kind == yaml.MappingNode {
			for i := 0; i < len(node.Content); i += 2 {
				if node.Content[i].Value == path[0] {
					node.Content[i].HeadComment = comment
					return true
				}
			}
		}
		return false
	}

	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			if node.Content[i].Value == path[0] {
				return SetYamlComment(node.Content[i+1], path[1:], comment)
			}
		}
	}

	return false
}

func LogVerbose(format string, a ...interface{}) {
	if configShared.Verbose {
		color.HiBlack(format, a...)
	}
}

//func pressAnyKeyToContinue() {
//	fmt.Println("Press any key to continue...")
//
//	if err := keyboard.Open(); err != nil {
//		panic(err)
//	}
//	defer keyboard.Close()
//
//	_, key, err := keyboard.GetKey()
//	if err != nil {
//		panic(err)
//	}
//}
