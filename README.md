# Http-Server-Test
HTTP server for proxying HTTP-requests to 3rd-party services.
# To run locally
Clone the project
```
git clone https://github.com/zhas-off/http-server-test
```
Start the server
```
go run cmd/main.go
```
## Send requests to the server with any HTTP method except "DELETE"
http://localhost:8080/
```
{
    "method": "POST",
    "url": "http://google.com",
    "headers": {
        "Authentication": "Basic bG9naW46cGFzc3dvcmQ=",
        ....
    }
}
```
