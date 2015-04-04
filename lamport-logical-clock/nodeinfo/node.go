package nodeinfo
import (
    "net"
    "log"
    "bytes"
    "github.com/ershadmoi/go-projects/lamport-logical-clock/utils"
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

var connectionmap = make(map[string] net.Conn)

func SetupConnections(nodenum string, nodes map[string] Node,
                nodemap map[string] Nodes, connections [][]int) {
    // Listen for all connection updates on a channel
    c := make(chan Node, 10)
    go func() {
        // Defer the close on this channel till condition breaks
        defer close(c)

        // Listen to connections until a close is called on the channel
        for val := range c {
            log.Println("Received a new connection to track to node : ", val.Nodenum)
            connectionmap[val.Nodenum] = val.Connection

            // If we have same number of connections as the number of
            // nodes we track, then break out
            if len(connectionmap) == len(nodemap[nodenum]) {
                log.Println("All nodes have been connected successfully.")
                break
            } else {
                log.Println("Connected : " , len(connectionmap), " num of nodes. Total nodes : ", len(nodemap[nodenum]))
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

func startServer(node Node, c chan Node) {
    log.Println("Starting server for node : ", node);
    listenSock, err := net.Listen("tcp", ":" + node.Portnum)
    if err != nil {
        log.Fatal("Error starting a listen operation on port : ", node.Portnum, "\n Error : ", err)
    }

    for {
        conn, err := listenSock.Accept()
        if err != nil {
            log.Println("Error accepting connection.")
            continue
        }

        // Go read data when available on this connection
        go readWhenDataIsAvailable(conn)

        // Send this new connection object over the channel
        // To track our connections
        // Need to do some lookup on our entries to figure out who we just accepted connections from
        // c <- conn
    }
}

func readWhenDataIsAvailable(conn net.Conn) {
    // Defer closing the connection
    defer conn.Close()

    // Read as long as there is data
    buffer := make([]byte, 256)
    for {
        _, err := conn.Read(buffer)
        if err != nil {
            log.Println("Error while reading from connection : ", conn.RemoteAddr())
            break
        }

        if(bytes.Equal(buffer, []byte("quit"))) {
            conn.Close()
            break
        }

        log.Println("Read data from client : ", conn.RemoteAddr(), " with value : ", string(buffer))
    }

    log.Println("Connection closed from client : ", conn.RemoteAddr())
}

func connectToNodesBeforeUs(self string, nodes Nodes, c chan Node) {
    for _, node := range nodes {
        if utils.GetInt(node.Nodenum) < utils.GetInt(self) {
            log.Println("Node ", self, " is attempting to connect to node : ", node.Nodenum)
            // Connect to to the node here
            conn, err := net.Dial("tcp", node.Hostname + ":" + node.Portnum)
            if err != nil {
                log.Println("Error when connecting to node : ", node.Nodenum, " on host : ", node.Hostname, " on portnum : ", node.Portnum)
                return
            }

            log.Println("Connected Successfully  to node : ", node.Nodenum)

            // Write our nodenum alone
            conn.Write([]byte(self))
            // Send the connected successfully connection
            // to the channel
            node.Connection = conn
            c <- node

            // Go read data when available on this connection
            go readWhenDataIsAvailable(conn)
        }
    }
}