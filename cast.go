package main

import "fmt"
import "os/exec"
import "os"
import "bufio"
import "sync"

const color_red = "\033[0;31m"
const color_clear = "\033[0m"

func printOut(prefix string, x string) {
  fmt.Printf("[%s] %s\n", prefix, x)
}

func printErrString(prefix string, x string) {
  fmt.Printf(color_red + "[%s] %s\n" + color_clear, prefix, x)
}

func printErr(prefix string, x error) {
  printErrString(prefix, x.Error())
}

func cmdToString(x exec.Cmd) string {
  s := x.Args[0] 
  for i := 1; i < len(x.Args); i++ {
    s += " '" + x.Args[i] + "'"
  }
  return s
}

func ssh(target string, command string, globalWaitGroup *sync.WaitGroup) {
  cmd := exec.Command("/usr/bin/ssh", target, command)
  outPipe, _ := cmd.StdoutPipe()
  errPipe, _ := cmd.StderrPipe()

  printOut("cast", cmdToString(*cmd))
  err := cmd.Start()

  if err != nil {
    printErr(target, err)
    os.Exit(1)
  }

  outScanner := bufio.NewScanner(outPipe)
  errScanner := bufio.NewScanner(errPipe)

  var wg sync.WaitGroup
  wg.Add(2)

  go func(out bufio.Scanner) {
    for(out.Scan()) {
      printOut(target, out.Text())
    }
    wg.Done()
  }(*outScanner)

  go func(err bufio.Scanner) {
    for (err.Scan()) {
      printErrString(target, err.Text())
    }
    wg.Done()
  }(*errScanner)

  wg.Wait()
  err = cmd.Wait()

  if err != nil {
    printErr(target, err)
  }

  globalWaitGroup.Done()
}

func main() {
  var waitGroup sync.WaitGroup
  waitGroup.Add(1)
  go ssh("oxblood3", "ls -la /", &waitGroup)
  waitGroup.Wait()
}


