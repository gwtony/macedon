SCHEMA
======

Creat Record
-----------

```
A:
Request
POST /macedon/create HTTP/1.1
Content-Type: application/json
{
	"Name": ${name},                 //service name
	"Address": ${address},           //address (May be A record or CNAME record)
	"Port": ${port}                  //port
}

Response
HTTP/1.1 200 OK

[Example]
A:
Request
POST /macedon/create HTTP/1.1
Content-Type: application/json
{
	"Name": "test.example.com",
	"Address": "192.168.0.2",
	"Port": "8080"
}

CNAME:
Request
POST /macedon/create HTTP/1.1
Content-Type: application/json
{
	"Name": "play.example.com",
	"Address": "test.example.com"
}
```

Delete Record
------------
```
Request
POST /macedon/delete HTTP/1.1
Content-Type: application/json
{    
	"Name": ${name},      //${name} may be a service name
	"Address": ${address} //${address} is optional
}
Response
HTTP/1.1 200 OK

Response record not found
Response
HTTP/1.1 204 No Content

[Example]
Request
POST /macedon/delete HTTP/1.1
Content-Type: application/json
{
	"Name": "test.example.com",
}

Successful Response
Response
HTTP/1.1 200 OK
```

Read Record
-----------
```
Request
POST /macedon/read HTTP/1.1
Content-Type: application/json
{
	"Name": ${name},      //service name
	"Address": ${address} //address is optional
}

Response
HTTP/1.1 200 OK
{
	"Result": [
		{
			"Name": ${name},
			"Address": ${address},
			"Port": ${port}
		},
		{ ... }
	]
}

Response record not found
Response
HTTP/1.1 204 No Content 
[Example]

Request
POST /macedon/read HTTP/1.1
Content-Type: application/json
{
	"Name": "test.example.com",
}

Response
HTTP/1.1 200 OK
{
	"Result": [
		{
			"Name": "test.example.com",
			"Address": "192.168.0.1",
			"Port": 80
		},
		{
			"Name": "test.example.com",
			"Address": "192.168.0.2",
			"Port": 81
		}
	]
}
```

Status Code 
-----------

* 200 - Success
* 204 - No Content (Not found record)
* 400 - Bad request error (Arguments invalid)
* 404 - Page not found (Location incorrect)
* 500 - Internal server error (Api server internal error)
* 502 - Bad Gateway (Backend error)
