package main

import (
   "fmt"
   "time"
   "math/rand"
   "io"
   "flag"
   "os"
   "log"
   u "github.com/ardeshir/version"
)

var ( 
 debug bool = false
 version string = "0.0.0"
 )
 
 

func main() {
    
    fmt.Println("Basic channel:")
    
    c := make(chan string)
    
    
    go boring("Boring!", c)
    
    for i := 0; i<3; i++ {
        fmt.Printf("You say: %q\n", <-c) // Receive expression is just value.
    }
  
    fmt.Println("Sending Channels:")
    
    cin1 := boringChan("Srvc1")
    cin2 := boringChan("Srvc2")
    
    
    for i := 0; i<6; i++ {
        fmt.Printf("Chan1: %q\n", <-cin1) 
        fmt.Printf("Chan2: %q\n", <-cin2) 
        
    }
    
      fmt.Println("Using FanIn:")  
      
      chan3 :=  fanIn(boringChan("Fn1_Srvc0"), boringChan("Fn1_Srvc1"))
      
       for i := 0; i<10; i++ {
        fmt.Println(<-chan3) // Receive fanIn channels 
       }
       
       
      chan4 :=  fanIn2(boringChan("Fn2_Srvc2"), boringChan("Fn2_Srvc3"))
      
      fmt.Println("Using Select with FanIn 2:")  
      
      for i := 0; i<10; i++ {
        fmt.Println(<-chan4) // Receive fanIn channels 
       }
    
    
      fmt.Println("main exits...")
    
    /*************************************/
    /*   THIS IS ALL JUST IO TESTING     */
    /*************************************/
     filename    := flag.String("file", defaultFile(), "Name of 1st file to use with createFile() ")
     // text        := flag.String("text", "This default text will be printed\n", "Some text goes here")
     // newfilename := flag.String("newfile", defaultFile2(), "Name of 2nd file use with createFile() ")
     flag.Parse()
    
     deleteFile(*filename)
     // delete file used the checkExistance before deleting.
     // checkExistence(*filename) 
     // createFile(*filename)
     // deleteFile(*filename)
     // writeToFile(*filename, *text)
     // renameFile(*filename, "newfiletest.txt")
     // copyFile(*filename, *newfilename)
     // renameFile(*newfilename, "newfiletest2.txt")
     // deleteFile(*newfilename)
     // writeToFile(*filename, *text)
      
  if debugTrue() {
    u.V(version)
  }

}

/******* CHAN FUNCTION ********/

func boring(msg string, c chan string) {
    for i := 0; ; i++ {
        c <- fmt.Sprintf("%s %d", msg, i )  // Expression to be sent can be any suitable value.
        time.Sleep(time.Duration( rand.Intn(1e3) ) * time.Millisecond)
    }
}

func boringChan(msg string) <-chan string { // Returns receive-only channel of strings
  c := make(chan string)
  
  go func() {  // We launch the goroutine from inside the function
        for i := 0; ; i++ {
            c <- fmt.Sprintf("%s %d", msg, i)
            time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
        }
  }()

  return c // Return the channel to the caller
}

// Multiplexing

func fan(input1 <- chan string) <- chan string {
     c := make(chan string)
     timeout := time.After(2 * time.Second)
      go func() {
       for {
           select {
               
               case s := <-input1 : c <- s
                     // fmt.Println(s)
               case <-timeout:
                     fmt.Println("Joe is not talking, bye!")          
               return    
           }
       }
      }()
    return c
}


func fanIn(input1, input2 <-chan string) <-chan string {
    c := make(chan string)
    go func() { for { c <- <- input1 } } ()
    go func() { for { c <- <- input2 } } ()
    
    return c
}

func fanIn2(cin1, cin2 <-chan string) <-chan string {
    c := make(chan string)
    
    go func() { 
        
      for {
        select {
           case v := <-cin1: c <- v
                 // fmt.Printf("Select rcvd %v from cin1\n", v)
           case v := <-cin2: c <- v
                 // fmt.Printf("Select rcvd %v from cin2\n", v)
        }
      } 
    }()
    
    return c
}

/******************************************/
/*   IO FUNTIONS **************************/

func writeToFile(filename string, text string) {
    f, err := os.Create(filename)
    u.ErrNil(err, "wtf: Can't create new file")
    
    defer f.Close()
    
    if _, err := f.Write([]byte(text)); err != nil {
        log.Fatalln(err)
    }
    
    log.Println("wtf: Done\n")
}

// Function to copy a file
func copyFile(filename string, newfilename string) {
    of, err := os.Open(filename)
    u.ErrNil(err, "cpf: Can't Open old file")
    defer of.Close()
    
    nf, err2 := os.Create(newfilename)
    u.ErrNil(err2, "cpf: Can't create new file")
    defer nf.Close()
    
    bw, err3 := io.Copy(nf, of)
    u.ErrNil(err3, "cpf: Can't copy from old to new")
    log.Printf("cpf: Bytes written: %d\n", bw)
    
    if err4 := nf.Sync(); err4 != nil {
        log.Fatalln(err4)
    }
    log.Printf("cpf: Done copying from %s to %s\n", filename, newfilename)

}


// Function to rename a file
func renameFile (filename string, newname string) {
        checkExistence(filename)
        err := os.Rename(filename, newname)
        u.ErrNil(err, "rnf: File was corrupted")
        
    fi, err2 := os.Stat(newname)
    if err2 != nil {
        if os.IsNotExist(err2) {
           log.Fatalf("rnf: File: %s does not exist", fi.Name)
        } 
    } 
    
    log.Printf("rnf: Exists, last modified %v\n", fi.ModTime())
    // deleteFile(newname)
}


// Function to create a file 
func createFile( filename string) {
    f, err := os.Create(filename)
    u.ErrNil(err, "cf: Unable to create file")
    defer f.Close()
    log.Printf("cf: Created %s\n", f.Name())
}

// Function to delete a file 
func deleteFile( filename string) {
     // createFile()  // use the createFile first, then delete file
     checkExistence(filename)
     err := os.Remove(filename)
     u.ErrNil(err, "df: Unable to remove testfile1.txt")
     log.Printf("df: Deleted %s", filename)
}

func checkExistence( filename string) {
    
    fi, err := os.Stat(filename)
    if err != nil {
        if os.IsNotExist(err) {
           log.Printf("cke: Test File: does not exist")
           
           createFile(filename)
           // fi, err = os.Stat(filename)
           return 
        } 
    } 
    
    log.Printf("cke: Exists, last modified %v\n", fi.ModTime())
    // deleteFile(filename)
}


// Function to check env variable DEFAULT_DEBUG bool
func debugTrue() bool {
     
     if os.Getenv("DEFAULT_DEBUG") != "" {
        return true
     }  
     return false 
}

// Function to check env variable DEFAULT_FILE to get
func defaultFile() string {
    if os.Getenv("DEFAULT_FILE") != "" {
        return os.Getenv("DEFAULT_FILE")
    }
    return "testfile.txt"
}

func defaultFile2() string {
    if os.Getenv("DEFAULT_FILE2") != "" {
        return os.Getenv("DEFAULT_FILE2")
    }
    return "newtestfile.txt"
}
