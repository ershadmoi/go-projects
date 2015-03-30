package nodeinfo

type Node struct {
    Hostname string
    Portnum,Nodenum int
}

type Nodes []Node

func (node Node) HashCode() string {
    hashcode := node.Hostname + ":" + string(node.Portnum)
    return hashcode
}