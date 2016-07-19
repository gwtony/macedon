package message

type Record struct {
	Content     string
	Disabled    int
}

type Request struct {
	Name        string
	Type        string
	Domain_id   int
	Ttl         int
	Records     []Record
}

type ResponseData Request

type ResponseResult struct {
	Affected    int
	Data        ResponseData
}

type Response struct {
	Result	ResponseResult
}
