package main

import (
  "fmt"
  "os"
  )

func main() {
  var arg_list string
  var sep string

  for i := 1; i < len(os.Args); i++ {
    arg_list += sep + os.Args[i]
    sep = " "
  }

  fmt.Println(arg_list)
}
