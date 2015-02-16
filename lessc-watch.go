package main

import(
  "bytes"
  "errors"
  "fmt"
  "log"
  "os"
  "os/exec"
  "path"
  "path/filepath"
  "strings"
  "gopkg.in/fsnotify.v1"
)

type Settings struct {
  indir, outdir string
  compress bool
}

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

func compileLessFile(input string, output string, compress bool) (err error, errstr string) {
  cmd := exec.Command("lessc", evaluateCompressionFlag(compress), input, output)

  var stderr bytes.Buffer
  cmd.Stderr = &stderr

  err = cmd.Start()
  if err != nil {
    return err, ""
  }

  err = cmd.Wait()
  if (err != nil) {
    return err, stderr.String()
  } else {
    return err, ""
  }
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

func watch(watcher *fsnotify.Watcher, settings *Settings) {
  for {
    select {
      case event := <-watcher.Events:
        if event.Op&fsnotify.Write == fsnotify.Write {
          filePath := event.Name
          if isfile, err := isFile(filePath); err == nil && isfile && filepath.Ext(filePath) == ".less" {
            in := filePath
            name := path.Base(filePath)
            base := getBaseName(name)
            out := settings.outdir + "/" + base + ".css"

            if err, errstr := compileLessFile(in, out, settings.compress); err == nil {
              log.Printf("Compiled %s\n", name)
            } else {
              log.Printf("An error occured (%s):\n%s", fmt.Sprint(err), errstr)
            }
          }
        }

      case err := <-watcher.Errors:
        log.Printf("Error: %s\n", err)
    }
  }
}

func getargs(args []string) (s *Settings, err error) {
  compress := false
  if args[0] == "-x" {
    compress = true
    args = args[1:]
  }

  indir, err := filepath.Abs(args[0])
  if err != nil {
    return nil, err
  }

  outdir := indir
  if len(args) > 1 {
    outdir, err = filepath.Abs(args[1])

    if err != nil {
      return nil, err
    }
  }

  if exist, err := directoriesExist(indir, outdir); !exist || err != nil {
    if !exist {
      err = errors.New("Error: Input and/or output directories do not exist.")
    } else {}

    return nil, err
  }

  return &Settings{
    indir: indir,
    outdir: outdir,
    compress: compress }, nil
}

func main() {
  args := os.Args[1:]
  if len(args) == 0 || (len(args) == 1 && args[0] == "-x") {
    fmt.Println("Usage: lessc-watch [-x] input-path [output-path]")
    os.Exit(0)
  }

  s, err := getargs(args)
  if (err != nil) {
    log.Fatal(err)
  }

  watcher, err := fsnotify.NewWatcher()
  if err != nil {
      log.Fatal(err)
  }
  defer watcher.Close()

  done := make(chan bool)
  go watch(watcher, s)

  err = watcher.Add(s.indir)
  if err != nil {
    log.Fatal(err)
  }

  <-done
}