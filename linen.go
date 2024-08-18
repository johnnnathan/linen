package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)


type Lines struct {
  comments , empty ,code int 
}

var progenitorDir , error = os.Getwd()
var total = make(map[string]*Lines)


func getFiles(directory string, files []string) []string{
  filesCurDir, _ := os.ReadDir(directory)
  for _, element := range filesCurDir {
		fullPath := filepath.Join(directory, element.Name())
		if element.IsDir() {
			files = getFiles(fullPath, files)
		} else {
			files = append(files, fullPath)
		}
	}
  return files
}

func readFile(fileDE string){
  file , err := os.Open(fileDE)
  ext := filepath.Ext(file.Name())

  if err != nil{
    log.Fatal(err)
  }
  scanner := bufio.NewScanner(file)
  if _, exists := total[ext]; !exists {
		total[ext] = &Lines{}
	}
  lines := total[ext]
  var inComment bool = false

  for scanner.Scan(){
    var text string = scanner.Text()
    fmt.Printf("scanner: %v %d\n", text, len(text))
    analyzeLine(text, lines, inComment)
  }
  file.Close()
}

func analyzeLine(line string, lines *Lines, inComment bool)bool{
  var text string = strings.TrimSpace(line)

  if len(text) == 0 || text == "\n"{
    lines.empty += 1
    return inComment
  }
  if !inComment && len(text) == 1{
    if text[0] == '#'{
      lines.comments += 1; return false
    }else{
      lines.code += 1; return false}
  }

  char1 := text[0]
  char2 := text[1]
  if !inComment && char1 != '#' && (char1 != '/' && char2 != '/') && (char1 != '-' && char2 != '-'){
    lines.code += 1 
    return false
  } 
  startBool, endBool := checkCommentSymbols(text)
  if !inComment && startBool && endBool{
    lines.comments += 1 
    return false
  }
  if inComment && endBool{
    lines.comments += 1 
    return false
  }
  if !inComment && startBool{
    lines.comments += 1 
    return true
  } 
  return inComment
}

func checkCommentSymbols(text string) (bool, bool) {
	startPattern := regexp.MustCompile(`/\*`)
	endPattern := regexp.MustCompile(`\*/`)

	foundStart := false
	foundEnd := false

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if startPattern.MatchString(line) {
			foundStart = true
		}
		if endPattern.MatchString(line) {
			foundEnd = true
		}
	}

	return foundStart, foundEnd
}
func readFiles(files []string)  {
  for _, element := range files{
    readFile(element)
  }
}

func main()  {
  var files []string
  fmt.Println(len("\n"))
  files = getFiles(progenitorDir, files)
  for _, element := range files{
    fmt.Println(element)
  }
  readFiles(files)
  for ext, lines := range total {
		fmt.Printf("Extension: %s\n", ext)
		fmt.Printf("Comments: %d\n", lines.comments)
		fmt.Printf("Empty Lines: %d\n", lines.empty)
		fmt.Printf("Code Lines: %d\n", lines.code)
	}
}
