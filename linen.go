package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

type Lines struct {
  comments , empty ,code int 
}
//go:embed VERSION.txt
var versionFS embed.FS

var progenitorDir , error = os.Getwd()
var total = make(map[string]*Lines)
var responseSent     = false    
var responseMutex    sync.Mutex 
var updateMutex      sync.Mutex
var mapMutex         sync.Mutex
var wg               sync.WaitGroup
var server *http.Server
var wantHTML = false 


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
  defer wg.Done()
  file , err := os.Open(fileDE)
  ext := filepath.Ext(file.Name())

  if err != nil{
    log.Fatal(err)
  }
  scanner := bufio.NewScanner(file)
  mapMutex.Lock()
  if _, exists := total[ext]; !exists {
		total[ext] = &Lines{}
	}
  lines := total[ext]
  mapMutex.Unlock()
  var inComment bool = false

  for scanner.Scan(){
    var text string = scanner.Text()
    inComment = analyzeLine(text, lines, inComment)
  }
  file.Close()
}

func incrementLineValue(lineType *int){
  updateMutex.Lock()
  defer updateMutex.Unlock()
  *lineType += 1 
}

func getVersion(){
  file , _:= versionFS.ReadFile("VERSION.txt")
  fmt.Printf("Version : %s\n", (string(file)))
}

func analyzeLine(line string, lines *Lines, inComment bool)bool{
  var text string = strings.TrimSpace(line)

	if strings.HasPrefix(text, "//") || strings.HasPrefix(text, "# ") || strings.HasPrefix(text, "--") {
    incrementLineValue(&lines.comments)
    return false
  } 
  if len(text) == 0 || text == "\n"{
    incrementLineValue(&lines.empty)
    return inComment
  }
  if !inComment && len(text) == 1{
    if text[0] == '#'{
      incrementLineValue(&lines.comments)
      return false
    }else{
      incrementLineValue(&lines.code)
      return false}
  }
  startBool, endBool := checkCommentSymbols(text)
  if inComment && endBool {
		incrementLineValue(&lines.comments)
		return false
	}

	if startBool && !endBool {
		incrementLineValue(&lines.comments)
		return true
	}

	if inComment || (startBool && endBool){
		incrementLineValue(&lines.comments)
		return inComment
	}


  incrementLineValue(&lines.code)
  return inComment
}

func checkCommentSymbols(text string) (bool, bool) {

	foundStart := false
	foundEnd := false

  if strings.HasPrefix(text, "/*"){
    foundStart = true
  }
  if strings.HasSuffix(text, "*/") {
    foundEnd = true
  }

	return foundStart, foundEnd
}
func readFiles(files []string)  {
  for _, element := range files{
    wg.Add(1)
    go readFile(element)
  }
  wg.Wait()
}

func dynamicHandler(w http.ResponseWriter, r *http.Request){
  responseMutex.Lock()
  defer responseMutex.Unlock()
  w.Header().Add("Content-Type", "text/html")
  w.WriteHeader(http.StatusOK)
  fmt.Fprintf(w, getHTML())
  responseSent = true
	go func() {
		time.Sleep(2 * time.Second)
		server.Shutdown(nil)
	}()
}

func getHTML() string{
  
  bar := "----------------------------------------------------<br>"
  var fullResult string
  fullResult += bar
	for ext, lines := range total {
    fullResult += fmt.Sprintf("<div style= \"color:green;\"> Extension: %s </div><br>", ext);
    fullResult += fmt.Sprintf("<div style = \"color:yellow;\"> Comments: %d</div><br>", lines.comments);
    fullResult += fmt.Sprintf("<div style = \"color:cyan;\"> Empty Lines: %d</div><br>", lines.empty);
    fullResult += fmt.Sprintf("<div style = \"color:red;\"> Code Lines: %d</div><br>", lines.code);
    fullResult += bar 
	}
  return fullResult
  
}

func printResults()  {
  fmt.Println(Blue + "Line Counting Results:" + Reset)
	fmt.Println("----------------------------------------------------")
	for ext, lines := range total {
		fmt.Printf(Green+"Extension: %s\n"+Reset, ext)
		fmt.Printf(Yellow+"Comments: %d\n"+Reset, lines.comments)
		fmt.Printf(Cyan+"Empty Lines: %d\n"+Reset, lines.empty)
		fmt.Printf(Red+"Code Lines: %d\n"+Reset, lines.code)
		fmt.Println("----------------------------------------------------")
	}
  
}

func openBrowser(url string) {

    switch runtime.GOOS {
    case "linux":
       exec.Command("xdg-open", url).Start()
    case "windows":
       exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
    case "darwin":
       exec.Command("open", url).Start()
    default:
        log.Fatalf("unsupported platform")
    }

}
func main()  {
  var versionFlag  = false
  flag.BoolVar(&wantHTML, "html", false, "Set to true if you want HTML output")
  flag.BoolVar(&versionFlag, "version", false, "Showcase the application version")
  flag.Parse()

  if versionFlag{
    getVersion()
    os.Exit(0)
  }
  var files []string
  files = getFiles(progenitorDir, files)
  readFiles(files)
  if !wantHTML{
    printResults()
    return
  }
  fmt.Println("Something")
  fs := http.FileServer(http.Dir("./static"))
  http.Handle("/", fs)

  http.HandleFunc("/dynamic", func(w http.ResponseWriter, r *http.Request){
    dynamicHandler(w, r)
  })
  server = &http.Server{Addr: ":8080"}
  go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
  openBrowser("http://localhost:8080/webUI.html")
  <-time.After(2 * time.Second) // Wait for 3 seconds before checking
	if responseSent {
		server.Shutdown(nil)
	}
}
