from PyQt5 import QtCore, QtWebSockets
from PyQt5.QtCore import QUrl


class Client(QtCore.QObject):
    def __init__(self, parent):
        super().__init__(parent)

        self.client =  QtWebSockets.QWebSocket("",QtWebSockets.QWebSocketProtocol.Version13,None)
        self.client.error.connect(self.error)

        self.client.open(QUrl("ws://socket.art2cat.com/ws/health"))
        self.client.textMessageReceived(message=)
        
        self.client.pong.connect(self.onMessage)

    def do_ping(self):
        print("client: do_ping")
        self.client.ping(b"foo")

    def send_message(self, ):
        print("client: send_message")
        self.client.sendTextMessage("asd")

    def onMessage(self, elapsedTime, payload):
        print("onPong - time: {} ; payload: {}".format(elapsedTime, payload))

    def error(self, error_code):
        print("error code: {}".format(error_code))
        print(self.client.errorString())

    def close(self):
        self.client.close()