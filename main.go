package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type sortFilename []string
type checkValidFile func(string) bool

var num = flag.String("num", "100", "number of gen name")

func (s sortFilename) Less(i, j int) bool {
	iNum, _ := strconv.Atoi(s[i])
	jNum, _ := strconv.Atoi(s[j])
	return iNum < jNum
}

func (s sortFilename) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortFilename) Len() int {
	return len(s)
}

func main() {
	jsFiles, _ := ioutil.ReadDir("./src/")
	jsFilePaths := walkAndGetFilePath(jsFiles, "src", isValidJsFile)
	resFiles, _ := ioutil.ReadDir("./res/")
	pngFilePaths := walkAndGetFilePath(resFiles, "res", isValidPngFile)
	fontFilePaths := walkAndGetFilePath(resFiles, "res", isValidFontFile)

	// project json
	projectJsonFileContent, err := ioutil.ReadFile("./project.json")
	if err != nil {
		panic(err)
	}
	var data map[string]interface{}
	err = json.Unmarshal(projectJsonFileContent, &data)
	if err != nil {
		panic(err)
	}

	// get the order out first
	jsOrderList := GetStringSliceAtPath(data, "jsListOrder")

	for _, otherFilePath := range jsFilePaths {
		if !ContainsByString(jsOrderList, otherFilePath) {
			jsOrderList = append(jsOrderList, otherFilePath)
		}
	}

	data["jsList"] = jsOrderList
	newProjectJsonFileContent, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("./project.json", newProjectJsonFileContent, 0644)
	if err != nil {
		panic(err)
	}

	// res file
	coreContent := ""
	for _, filePath := range pngFilePaths {
		index := strings.LastIndex(filePath, "/")
		name := filePath[index+1 : len(filePath)-4] // cut .png
		coreContent = fmt.Sprintf("%s\t%s : \"%s\",\n", coreContent, name, filePath)
	}

	coreFontContent := ""
	for _, filePath := range fontFilePaths {
		index := strings.LastIndex(filePath, "/")
		name := filePath[index+1 : len(filePath)-4] // cut .ttf
		coreFontContent = fmt.Sprintf(`
%s		
g_resources.push({
    type:"font",
    name:"%s",
    srcs:["%s"]
});
			`, coreFontContent, name, filePath)
	}

	resFileContent := fmt.Sprintf(`
var res = {
%s
};

var g_resources = [];
for (var i in res) {
    g_resources.push(res[i]);
}
%s
	`, coreContent, coreFontContent)

	err = ioutil.WriteFile("./src/resource.js", []byte(resFileContent), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Cog success!")
}

func walkAndGetFilePath(filesInfo []os.FileInfo, path string, isValidFunc checkValidFile) []string {
	paths := []string{}
	for _, f := range filesInfo {
		if f.IsDir() {
			subPath := fmt.Sprintf("%s/%s", path, f.Name())
			subFilesInfo, _ := ioutil.ReadDir(fmt.Sprintf("./%s", subPath))
			newPaths := walkAndGetFilePath(subFilesInfo, subPath, isValidFunc)
			paths = append(paths, newPaths...)
		} else {
			if isValidFunc(f.Name()) {
				newPath := fmt.Sprintf("%s/%s", path, f.Name())
				paths = append(paths, newPath)
			}
		}
	}
	return paths
}

func isValidJsFile(name string) bool {
	length := len(name)
	return name[length-3:] == ".js"
}

func isValidPngFile(name string) bool {
	length := len(name)
	return name[length-4:] == ".png"
}

func isValidFontFile(name string) bool {
	length := len(name)
	return (name[length-4:] == ".ttf")
}
