package main
import(
    "github.com/ershadmoi/go-projects/lamport-logical-clock/config"
    "github.com/ershadmoi/go-projects/lamport-logical-clock/nodeinfo"
    "flag"
    "fmt"
)

func main() {
    // Get the nodenum of this process from commandline
    nodenumPtr := flag.String("nodenum", "0", "The node number of this process")
    flag.Parse()

    nodes, nodemap, connections :=  config.ReadConfig("config/config.txt")
    go nodeinfo.SetupConnections(*nodenumPtr, nodes, nodemap, connections )

    // For now blocking on user input
    // So that main thread doesnt die
    var input string
    fmt.Scanln(&input)

    // Let's start some random simulation of send/receive events to update our clocks
}
