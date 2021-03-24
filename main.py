# -*- coding: utf-8 -*-

################################################################################
## Form generated from reading UI file 'form.ui'
##
## Created by: Qt User Interface Compiler version 5.15.2
##
## WARNING! All changes made in this file will be lost when recompiling UI file!
################################################################################
import sys
from multiprocessing import Process

from PyQt5 import QtWidgets, QtCore
from PyQt5.QtCore import QThread
from PyQt5.QtWidgets import QApplication, QMainWindow

from qws_client import Client


class WorkThread(QThread):
    def run(self):
        print("dd")


class Ui_Home(object):

    def setupUi(self, Home):
        Home.setObjectName("Home")
        Home.resize(800, 600)
        self.verticalLayoutWidget = QtWidgets.QWidget(Home)
        self.verticalLayoutWidget.setGeometry(QtCore.QRect(9, 9, 781, 591))
        self.verticalLayoutWidget.setObjectName("verticalLayoutWidget")
        self.verticalLayout = QtWidgets.QVBoxLayout(self.verticalLayoutWidget)
        self.verticalLayout.setSizeConstraint(QtWidgets.QLayout.SetDefaultConstraint)
        self.verticalLayout.setContentsMargins(8, 8, 8, 8)
        self.verticalLayout.setObjectName("verticalLayout")
        self.listView = QtWidgets.QListView(self.verticalLayoutWidget)
        self.listView.setObjectName("listView")
        self.verticalLayout.addWidget(self.listView)
        self.textEdit = QtWidgets.QTextEdit(self.verticalLayoutWidget)
        self.textEdit.setObjectName("textEdit")
        self.verticalLayout.addWidget(self.textEdit)
        self.horizontalLayout = QtWidgets.QHBoxLayout()
        self.horizontalLayout.setObjectName("horizontalLayout")
        self.connectBtn = QtWidgets.QPushButton(self.verticalLayoutWidget)
        self.connectBtn.setObjectName("connectBtn")
        self.connectBtn.clicked.connect(self.connectBtnPresssed)
        self.horizontalLayout.addWidget(self.connectBtn)
        self.disconnectBtn = QtWidgets.QPushButton(self.verticalLayoutWidget)
        self.disconnectBtn.setObjectName("disconnectBtn")
        self.disconnectBtn.clicked.connect(self.disconnectBtnPresssed)
        self.horizontalLayout.addWidget(self.disconnectBtn)
        self.verticalLayout.addLayout(self.horizontalLayout)

        self.retranslateUi(Home)
        QtCore.QMetaObject.connectSlotsByName(Home)

    def retranslateUi(self, Home):
        _translate = QtCore.QCoreApplication.translate
        Home.setWindowTitle(_translate("Home", "Home"))
        self.connectBtn.setText(_translate("Home", "Connect"))
        self.disconnectBtn.setText(_translate("Home", "Disconnect"))

    def connectBtnPresssed(self):
        print("hello")


    def disconnectBtnPresssed(self):
        print("fuck")


if __name__ == '__main__':

    app = QApplication([])
    windows = QMainWindow()
    widget = Ui_Home()
    widget.setupUi(windows)
    windows.show()

    client = Client(app)
    sys.exit(app.exec_())
