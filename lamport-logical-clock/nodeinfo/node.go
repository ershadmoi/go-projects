package nodeinfo
import (
    "net"
    "log"
    "fmt"
    "bytes"
)

type Node struct {
    Hostname string
    Portnum,Nodenum string
    Connection net.Conn
}

type Nodes []Node

func (node Node) HashCode() string {
    hashcode := node.Hostname + ":" + string(node.Portnum)
    return hashcode
}

var connectionmap = make(map[string] Node)

func SetupConnections(nodenum string, nodes map[string] Node,
                nodemap map[string] Nodes, connections [][]int) {
    // Listen for all connection updates on a channel
    c := make(chan Node, 10)
    go func() {
        // Defer the close on this channel till condition breaks
        defer close(c)

        // Listen to connections until a close is called on the channel
        for val := range c {
            connectionmap[val.Nodenum] = val.Connection

            // If we have same number of connections as the number of
            // nodes we track, then break out
            if len(connectionmap) == len(nodes[nodenum]) {
                break
            }
        }
    }()

    // Write data to all our connections whenever we put anything into the
    // broadcast channel
    broadcastChan := make(chan string, 500)
    go func() {
        for data := range broadcastChan {
            for _, val := range nodemap[nodenum] {
                log.Println("Sending data : ", val , " to node : ", val.Nodenum , " on connection : ", val.Connection)
                val.Connection.Write([]byte(data))
            }
        }
    } ()

    // Begin by starting a listen socket
    // Connections accepted by the server will be sent on channel c
    go startServer(nodes[nodenum], c)

    // We will then connect to any other node in our nodemap
    // which preceeds us
    go connectToNodesBeforeUs(nodenum, nodemap[nodenum], c)
}

func startServer(node Node, c chan net.Conn) {
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

        // Go read data when available on this connection
        go readWhenDataIsAvailable(conn)

        // Send this new connection object over the channel
        // To track our connections
        c <- conn
    }
}

func readWhenDataIsAvailable(conn net.Conn) {
    // Read as long as there is data
    buffer := make([]byte, 256)
    for conn.Read(buffer) {
        if(bytes.Equal(buffer, []byte("quit"))) {
            conn.Close()
            break
        }

        log.Println("Read data from client : ", conn, " with value : ", string(buffer))
    }

    log.Println("Connection closed from client : ", conn)
}

func connectToNodesBeforeUs(self int, nodes Nodes, c chan Node) {
    for _, node := range Nodes {
        if node.Nodenum < self {
            go func() {
                // Connect to to the node here
                conn, err := net.Dial("tcp", node.Hostname + ":" + node.Portnum)
                if err != nil {
                    log.Println("Error when connecting to node : ", node.Nodenum, " on host : ", node.Hostname, " on portnum : ", node.Portnum)
                }

                // Send the connected successfully connection
                // to the channel
                node.Connection = conn
                c <- node

                // Go read data when available on this connection
                go readWhenDataIsAvailable(conn)
            } ()
        }
    }
}