import QtQuick 2.7
import QtQuick.Controls 2.0
import QtQuick.Layouts 1.0

ApplicationWindow {
    visible: true
//    width: 640
//    height: 480
    width: 320
    height: 240
    minimumWidth: 320
    minimumHeight: 240
    title: qsTr("Websocket Client")

    ColumnLayout {
        anchors.fill: parent
        anchors.margins: 10

        Label {text: "Websocket URL:"}

        RowLayout {
            Layout.fillWidth: true
            TextField {
                id: urlField
                Layout.fillWidth: true
                text: "wss://localhost:8182"
                enabled: !WS.connected
                onAccepted: WS.connected ? closeCon() : openCon()
            }
            Button {
                Layout.maximumWidth: 100
                Layout.minimumWidth: 100
                text: WS.connected ? "Close" : "Open"
                onClicked: WS.connected ? closeCon() : openCon()
            }
        }

        Item {height: 5}

        Label {text: "Messenger:"}

        RowLayout {
            Layout.fillWidth: true
            TextField {
                id: msgField
                enabled: WS.connected
                Layout.fillWidth: true
                onAccepted: sendMsg()
            }
            Button {
                text: "Send"
                enabled: WS.connected
                onClicked: sendMsg()
                Layout.maximumWidth: 100
                Layout.minimumWidth: 100
            }
            Button {
                text: "Ping"
                enabled: WS.connected
                onClicked: sendPing()
                Layout.maximumWidth: 100
                Layout.minimumWidth: 100
            }
        }

        ListView {
            id: msgsView
            Layout.fillHeight: true
            Layout.fillWidth: true
            clip: true

            delegate: ItemDelegate{text: model.data}
            model: ListModel {
                id: msgsModel
            }
            ScrollBar.vertical: ScrollBar{}
        }
    }

    Component.onCompleted: {
        WS.onMsgRecieved.connect(function(obj){
            msgsModel.append(obj)
            msgsView.positionViewAtEnd()
        })
    }

    function sendMsg() {
        Messenger.sendMsg(msgField.text)
        msgField.text = ""
        msgField.forceActiveFocus()
    }

    function sendPing() {
        WS.sendPing(msgField.text)
        msgField.text = ""
        msgField.forceActiveFocus()
    }

    function openCon() {
        WS.open(urlField.text)
    }

    function closeCon() {
        WS.close()
        msgsModel.clear()
    }
}
