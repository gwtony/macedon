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
	"name": ${name},                 //domain name or ptr-ip(ptr ip of "192.168.0.1" is "1.0.168.192.in-addr.arpa")
	"type": ${type},                 //must be "A", "CNAME", "PTR"        
	"domain_id": ${domain_id}        //domain id number
	"ttl": ${TTL},                   //TTL
	"records": [
		{ "content" : ${content} }   //domain name or ip
		{ ... },
		{ ... }
	]
}

Response
HTTP/1.1 200 OK
Content-Type: application/json
{
	"result": { "affected": ${ affected objects number }}
}

[Example]
A:
Request
POST /macedon/create HTTP/1.1
Content-Type: application/json
{
	"name": "test.example.com",
	"type": "A",
	"domain_id": 1,
	"ttl": 30,
	"records": [ 
		{ "content" : "192.168.0.1" },
		{ "content" : "192.168.0.2" }           
	]
}

CNAME:
Request
POST /macedon/create HTTP/1.1
Content-Type: application/json
{
	"name": "play.example.com",
	"type": "CNAME",
	"domain_id": 1,
	"ttl": 30,
	"records": [      
		{ "content" : "test.example.com" }
	]
}

PTR:
Request
POST /macedon/create HTTP/1.1
Content-Type: application/json
{
	"name": "1.0.168.192.in-addr.arpa",
	"type": "PTR",
	"domain_id": 1,
	"ttl": 30,
	"records": [      
		{ "content" : "test.example.com" }
	]
}

Response
HTTP/1.1 200 OK Content-Type: application/json
{ 
	"result": { "affected": 1}
}
```

Delete Record
------------
```
Request
POST /macedon/delete HTTP/1.1
Content-Type: application/json
{    
	"name": ${name},      //${name} may be a domain name or ip or ptr-ip
	"type": ${type},      //${type} must be "A" or "CNAME" or "PTR"
	"records": [          //"records" is optional
		{ "content": { ... } },
		{ ... }
	]  
}
Response
HTTP/1.1 200 OK
Content-Type: application/json 
{
	"result": { "affected": ${ affected objects number } }
}

Response record not found
Response
HTTP/1.1 204 No Content

[Example]
Request
POST /macedon/delete HTTP/1.1
Content-Type: application/json
{
	"name": "test.example.com",
	"type": "A"
}

Successful Response
Response
HTTP/1.1 200 OK
{   "result": { "affected": 1 } }
```
Update Record
-------------
```
Request
POST /macedon/update HTTP/1.1
Content-Type: application/json
{
	"name": ${name},
	"type": "A",
	"records": [
		{
			"content": ${ip},
			"disabled": ${disabled}    //${disabled} is 1 or 0
		},
		{ ... }
	]
}


Response will return the old record
Response
HTTP/1.1 200 OK
{
	"result": {
		"affected": ${affected objects number},
		"data": {
			"name": ${name},
			"type": "A",
			"records": [
				{
					"content": ${ip},
					"disabled": ${disabled_old}  //old disabled state
				}
			]
		}
	}
}

Response record not found
Response
HTTP/1.1 204 No Content 
Example

Request
POST /macedon/update HTTP/1.1
Content-Type: application/json
{
	"name": "test.example.com",
	"type": "A",
	"records": [
		{
			"content": "192.168.0.1",
			"disabled": 1
		}
	]
}

Successful Response
HTTP/1.1 200 OK
{   "result": {
		"affected": 1, 
		"data": [{
			"name": "test.example.com",
			"type": "A",
			"records": [
				{
				"content": "192.168.0.1",
				"disabled": 0
				}
			]
		}]
	}
}
```
Read Record
-----------
```
Request
POST /macedon/read HTTP/1.1
Content-Type: application/json
{
	"name": ${name},  //domain name or ip
	"type": ${type}   //A or CNAME
}
}
Response
HTTP/1.1 200 OK
{
	"result": {
		"affected": ${affected objects number}, 
		"data": {
			"name": ${name},
			"type": ${type},
			"domain_id": ${domain_id},
			"ttl": ${ttl},
			"records": [
				{ 
					"content": ${content},   //ip or domain name
					"disabled": ${state}     //0 or 1
				},
				{ ... }
			]
		}
	}
}

Response record not found
Response
HTTP/1.1 204 No Content 
[Example]

Request
POST /macedon/dig HTTP/1.1
Content-Type: application/json
{
	"name": "test.example.com",
	"type": "A"
}

Response
HTTP/1.1 200 OK
{
	"result": {
		"data": {
			"name": "test.example.com",
			"type": "A",
			"domain_id": 1,
			"ttl": 86400,
			"records": [
				{
					"content": "192.168.0.1",
					"disabled": 0
				}
			]
		}
		"affected": 1
	}
}
```
Status Code 
-----------

* 200 - Success
* 204 - No Content
* 400 - Bad request error
* 500 - Internal server error
