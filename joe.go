package main

import (
   "fmt"
   "time"
   "os"
   "math/rand"
   u "github.com/ardeshir/version"
)

var (
 debug bool = false
 version string = "0.0.1"
)


func main() {

    c := boring("Joe")

    timeout := time.After(5 * time.Second)
    for {
        
       select {
          case s := <-c : 
            fmt.Println(s)
          case <- timeout: 
            fmt.Println("You talk too much")
            return 
       } 

   }


  if debugTrue() {
    u.V(version)
  }

}

func boring(msg string) <-chan string { // Returns receive-only channel of strings

  c := make(chan string)
  
  go func() {  // We launch the goroutine from inside the function
        for i := 0; ; i++ {
            c <- fmt.Sprintf("%s %d", msg, i)
            time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
        }
  }()

  return c // Return the channel to the caller
}

// Function to check env variable DEFAULT_DEBUG bool
func debugTrue() bool {
     
     if os.Getenv("DEFAULT_DEBUG") != "" {
        return true
     }  
     return false 
}