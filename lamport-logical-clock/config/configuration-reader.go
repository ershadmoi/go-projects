package config

import(
    "bufio"
    "os"
    "log"
    "strings"
    "github.com/ershadmoi/go-projects/lamport-logical-clock/nodeinfo"
    "strconv"
    "github.com/ershadmoi/go-projects/lamport-logical-clock/utils"
)

/* Read the config file into a list of nodes and connections) */
func ReadConfig(filename string) (nodes map[string] nodeinfo.Node,
            nodemap map[string] nodeinfo.Nodes, connections [][]int) {
    // Open the config file
    file, err :=  os.Open(filename)

    // Check for errors
    if err != nil {
        log.Fatal("File not found : " , filename)
    }

    // Defer closing the file
    defer file.Close()

    // Initialize a 2D slice for connections
    connections = make([][]int, 10)
    for i := range connections {
        connections[i] = make([]int, 10)
    }

    // Initialize a map of nodes to store all node information when reading
    nodes = make(map[string] nodeinfo.Node)
    nodemap = make(map[string] nodeinfo.Nodes)

    // Get a new scanner
    scanner := bufio.NewScanner(bufio.NewReader(file))

    // Scan through input lines
    for scanner.Scan() {

        // Empty switch on the entry value
        switch entry := scanner.Text(); {
            // Node information case
            case strings.HasPrefix(entry, "@") :
                processNodeEntry(nodes, entry[1:])

            // Connection information case
            case strings.HasPrefix(entry, "!") :
                processConnectionEntry(connections, entry[1:])

            // Comment line case
            case strings.HasPrefix(entry, "#") :
                log.Println("Comment line : " , entry)

            default :
                log.Println("Unknown file entry : " , entry)
        }
    }

    // Build a map of a list of nodes per source
    updateNodeMap(nodes, nodemap, connections)

    return
}

/* Updates the map with source nodes to a list of destination nodes */
func updateNodeMap(nodes map[string] nodeinfo.Node,
                   nodemap map[string] nodeinfo.Nodes,
                   connections [][]int) {
    // Traverse through the entire connection array
    for i,tmp := range connections {
        for j,val := range tmp {
            // Check If there is an uplink
            if val == 1 {
                // Append the destination node
                nodemap[strconv.Itoa(i)] = append(nodemap[strconv.Itoa(i)], nodes[strconv.Itoa(j)])
            }
        }
    }
}

/* Process the node entry read from the file - Uses reference based updates */
func processNodeEntry(nodes map[string] nodeinfo.Node, entry string) {
    log.Println("Processing node entry : ", entry)
    log.Println("If you see index out of bound errors, please check your config.txt for unwanted spaces.")

    // Tokenize the line entry
    tokens := strings.Split(entry, " ")

    // if there are tokens left
    if tokens != nil {
        // We will just blindly fail if format is not met
        // With index out of bound exceptions and some error logging
        // No efforts to handle bad config cases
        log.Println("Updating node nodes: ", nodes);
        nodes[tokens[1]] = nodeinfo.Node{
            Nodenum: tokens[1], Hostname: tokens[2], Portnum: tokens[3],
        }
    }
}

/* Process the connection entry read from the file - Uses reference based updates */
func processConnectionEntry(connections [][]int, entry string) {
    log.Println("Processing connection entry : ", entry)
    tokens := strings.Split(entry, " ")
    source := utils.GetInt(tokens[1])
    log.Println("Processing connections for source node : " , source)

    for i, val := range tokens {
        if i > 1 {
            log.Println("Adding connection link : ", source , " dest " , val)
            connections[source][utils.GetInt(val)] = 1
        }
    }
}

