package version

// Help for developers:
// go build -ldflags "-X main.version=0.0.1"
// godoc -http :9090
// root@ardeshir ~/httpdoc (master) $ graphpkg -stdout net > graph
// ardeshir.org/graph to view the SVG of all dependencies 


import (
    "fmt"
    "time"
    "log"
)

func V(version string) {


 fmt.Printf("Version: %s  Univrs.io\n", version)
 fmt.Printf("Current time: %s\n", time.Now())

}

func ErrNil ( err error, mesg string) {
		if err != nil {
		log.Fatalln(mesg)
	}
}
