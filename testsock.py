from __future__ import print_function
import websocket

import time

if __name__ == "__main__":
    websocket.enableTrace(True)
    ws = websocket.create_connection("ws://localhost:8080/control")
    print('{"type":"GET_CONFIG"}')

    time.sleep(1.0)

    ws.send('{"type":"GET_CONFIG"}\n')
    print("Sent")

    ws.send('{"type":"GET_CONFIG"}\n')
    print("Sent")

    print ("...")

    print("Receiving...")
    result = ws.recv()
    print("Received '%s'" % result)
    ws.close()
