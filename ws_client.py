import json
import uuid
import time

import websocket
from win10toast import ToastNotifier
from multiprocessing import Process

try:
    import thread
except ImportError:
    import _thread as thread

client_id = {}


def set_client_id(id):
    client_id["id"] = id


def get_client_id():
    return str(client_id["id"])


wsApp = ''


def connect():
    websocket.enableTrace(True)
    # ws = websocket.WebSocketApp("ws://213.59.119.106:8089/ws/health",
    wsApp = websocket.WebSocketApp("ws://socket.art2cat.com/ws/health",
                                   on_open=on_open,
                                   on_message=on_message,
                                   on_error=on_error,
                                   on_close=on_close)

    wsApp.run_forever()


def send_message(message):
    messaged = json.dumps(message, sort_keys=True)
    wsApp.send(messaged)


def on_message(ws, message):
    print(message)
    payload = json.loads(message)
    type_ = payload["type"]
    if "showNotice" == type_:
        p = Process(target=show_toast, args=(payload["title"], payload["message"]))
        p.start()
        p.join()
    elif "heartbeat" == type_:
        payload["clientId"] = get_client_id()
        print(payload)
        messaged = json.dumps(payload, sort_keys=True)
        ws.send(messaged)


def on_error(ws, error):
    print(error)


def on_close(ws):
    print("### closed ###")
    print("reconnecting......................")
    connect()


def on_open(ws):
    def run(*args):
        id_ = uuid.uuid4()
        set_client_id(id_)
        timestamp = time.time()
        message = json.dumps({"type": "handshake", "clientId": str(id_), "timestamp": timestamp}, sort_keys=True)
        ws.send(message)

    thread.start_new_thread(run, ())


def show_toast(title, message):
    toaster = ToastNotifier()
    toaster.show_toast(title, message,
                       icon_path=None,
                       duration=5)
