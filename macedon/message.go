package macedon

/* Http message struct */
type Request struct {
	Name        string
	Address     string
	Port        int
}

type Response []Request

/* Consul message struct */
type ConsulRequest struct {
	ID      string
	Name    string
	Address string
	Tags    []string
	Port    int
}

type CResponse struct {
	Node           string
	Address        string
	ServiceID      string
	ServiceName    string
	ServiceAddress string
	ServicePort    int
	ServiceTags    []string
}

type ConsulResponse []CResponse

