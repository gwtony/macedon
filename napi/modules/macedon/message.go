package macedon

/* Http message struct */
type MacedonRequest struct {
    Name    string
    Addr    string
	Ttl     int
}

type MacedonAddress struct {
	Addr string
	Ttl  int
}

type MacedonResponse struct {
	Name string
	Addrs []MacedonAddress
}


/* Etcd request */
type RecValue struct {
    Host string
    Ttl  int
}

type SubNode struct {
    Key string
    Value string
}

type Node struct {
    Key string
    Nodes []SubNode
}

type EtcdResponse struct {
    Node Node
}
