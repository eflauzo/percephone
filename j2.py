import pyjsonrpc

http_client = pyjsonrpc.HttpClient(
    url = "http://127.0.0.1:8080/rpc",
)
print http_client.call("HelloService.Say", {"Who":"123"})
