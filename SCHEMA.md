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
	"name": ${name},                 //service name
	"address": ${address},           //address (May be A record or CNAME record)
	"port": ${port}                  //port
}

Response
HTTP/1.1 200 OK

[Example]
A:
Request
POST /macedon/create HTTP/1.1
Content-Type: application/json
{
	"name": "test.example.com",
	"address": "192.168.0.2",
	"port": "8080"
}

CNAME:
Request
POST /macedon/create HTTP/1.1
Content-Type: application/json
{
	"name": "play.example.com",
	"address": "test.example.com"
}
```

Delete Record
------------
```
Request
POST /macedon/delete HTTP/1.1
Content-Type: application/json
{    
	"name": ${name},      //${name} may be a service name
	"address": ${address} //${address} is optional
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
	"name": "test.example.com",
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
	"name": ${name},      //service name
	"address": ${address} //address is optional
}

Response
HTTP/1.1 200 OK
[
	{
		"name": ${name},
		"address": ${address},
		"port": ${port}
	},
	{ ... }
]

Response record not found
Response
HTTP/1.1 204 No Content 
[Example]

Request
POST /macedon/read HTTP/1.1
Content-Type: application/json
{
	"name": "test.example.com",
}

Response
HTTP/1.1 200 OK
[
	{
		"name": "test.example.com",
		"address": "192.168.0.1",
		"port": 80
	},
	{
		"name": "test.example.com",
		"address": "192.168.0.2",
		"port": 81
	}
]
```

Status Code 
-----------

* 200 - Success
* 204 - No Content (Not found record)
* 400 - Bad request error (Arguments invalid)
* 404 - Page not found (Location incorrect)
* 500 - Internal server error (Api server internal error)
* 502 - Bad Gateway (Backend error)
