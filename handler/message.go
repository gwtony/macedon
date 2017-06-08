package handler

// MacedonRequest is Macedon request
type MacedonRequest struct {
    Name    string
    Address string
    Ttl     int
	Token   string
}

// MacedonUpdateRequest is Macedon request
type MacedonUpdateRequest struct {
    Name    string
    Address string
	Old     string
    Ttl     int
	Token   string
}

type MacedonResponseRecord struct {
    Name    string
    Address string
    Ttl     int
}

// MacedonResponse is Macedon response
type MacedonResponse struct {
    Result []MacedonResponseRecord
}


// RecValue is Etcd request
type RecValue struct {
    Host string
    Ttl  int
}

// Node is Etcd node
type Node struct {
    Key   string
    Value string
    Dir   bool
    Nodes []Node
}

// EtcdResponse is Etcd response
type EtcdResponse struct {
    Node Node
}

// ServerRequest is Server request
type ServerRequest struct {
    Address string
	Token   string
}

// ServerResponseRecord is Server response record
type ServerResponseRecord struct {
    Address string
}

// ServerResponse is Server response
type ServerResponse struct {
    Result []ServerResponseRecord
}
