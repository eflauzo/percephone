import requests
import json


def main():
    url = "http://localhost:8080/rpc"
    headers = {'content-type': 'application/json'}

    # Example echo method
    payload = {
        "method": "say",
        "params": ["echome!"],
        "jsonrpc": "1.0",
        "id": 0,
    }
    response = requests.post(
        url, data=json.dumps(payload), headers=headers).json()
    print response
    #assert response["result"] == "echome!"
    #assert response["jsonrpc"]
    #assert response["id"] == 0

if __name__ == "__main__":
    main()
