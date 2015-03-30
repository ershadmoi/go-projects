package nodeinfo
import (
    "net"
    "log"
    "fmt"
)

type Node struct {
    Hostname string
    Portnum,Nodenum string
}

type Nodes []Node

func (node Node) HashCode() string {
    hashcode := node.Hostname + ":" + string(node.Portnum)
    return hashcode
}

func SetupConnections(nodenum string, nodes map[string] Node,
                nodemap map[string] Nodes, connections [][]int) {
    // Begin by starting a listen socket
    go startServer(nodes[nodenum])
    // We will then connect to any other node in our nodemap
    // which preceeds us
}

func startServer(node Node) {
    fmt.Println("Starting server on port : " , node.Portnum);
    listenSock, err := net.Listen("tcp", ":" + node.Portnum)
    if err != nil {
        log.Fatal("Error starting a listen operation on port : ", node.Portnum, "\n Error : ", err)
    }

    for {
        conn, err := listenSock.Accept()
        if err != nil {
            fmt.Println("Error accepting connection.")
            continue
        }

        // Handle this connection if accept is successful
        go handleServerConnection(conn)
    }
}

func handleServerConnection(conn net.Conn) {

}