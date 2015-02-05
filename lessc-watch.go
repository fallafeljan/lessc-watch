package main

import(
  "fmt"
  "log"
  "os"
  "os/exec"
  "path"
  "path/filepath"
  "strings"
  "gopkg.in/fsnotify.v1"
)

func directoriesExist(paths ...string) (exist bool, err error) {
  for _, path := range paths {
     src, err := os.Stat(path)
     if err != nil || !src.IsDir() {
        return false, err
     }
  }

  return true, nil
}

func evaluateCompressionFlag(compress bool) string {
  if compress {
    return "-x"
  } else {
    return ""
  }
}

func compileLessFile(input string, output string, compress bool) error {
  cmd := exec.Command("lessc", evaluateCompressionFlag(compress), input, output)
  err := cmd.Start()
  if err != nil {
    return err
  }

  return cmd.Wait()
}

func getBaseName(filename string) string {
  return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func isFile(path string) (isfile bool, err error) {
  src, err := os.Stat(path)
  if err != nil || !src.Mode().IsRegular() {
    return false, err
  }

  return true, nil
}

func main() {
  args := os.Args[1:]

  if len(args) == 0 || (len(args) == 1 && args[0] == "-x") {
    fmt.Println("Usage: lessc-watch [-x] input-path [output-path]")
    os.Exit(0)
  }

  compress := false
  if args[0] == "-x" {
    compress = true
    args = args[1:]
  }

  input, err := filepath.Abs(args[0])
  if err != nil {
    log.Fatal(err)
  }

  output := input
  if len(args) > 1 {
    output, err = filepath.Abs(args[1])

    if err != nil {
      log.Fatal(err)
    }
  }

  if exist, err := directoriesExist(input, output); !exist || err != nil {
    if !exist {
      log.Fatal("Error: Input and/or output directories do not exist.")
    } else if err != nil {
      log.Fatal(err)
    }
  }

  watcher, err := fsnotify.NewWatcher()
  if err != nil {
      log.Fatal(err)
  }
  defer watcher.Close()

  done := make(chan bool)
  go func() {
    for {
      select {
        case event := <-watcher.Events:
          if event.Op&fsnotify.Write == fsnotify.Write {
            filePath := event.Name
            if isfile, err := isFile(filePath); err == nil && isfile && filepath.Ext(filePath) == ".less" {
              in := filePath
              name := path.Base(filePath)
              base := getBaseName(name)
              out := output + "/" + base + ".css"

              if err := compileLessFile(in, out, compress); err == nil {
                log.Printf("Compiled %s\n", name)
              } else {
                log.Print(err)
              }
            }
          }
      }
    }
  }()

  err = watcher.Add(input)
  if err != nil {
    log.Fatal(err)
  }

  <-done
}