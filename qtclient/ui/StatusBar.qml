import QtQuick 2.7
import QtQuick.Controls 2.0
import QtQuick.Layouts 1.0
import QtQuick.Controls.Material 2.1

ToolBar {
    id: statusBar
    Material.elevation: 0
    Material.primary: Material.accent
    RowLayout {
        anchors.fill: parent
        anchors.leftMargin: 10
        anchors.rightMargin: 10
        spacing: 10

        Item {
            Layout.fillWidth: true
            Layout.fillHeight: true
        }

        Label {
            id: messageText
            Layout.fillHeight: true
            verticalAlignment: Text.AlignVCenter
            text: "disconnected."
        }

        Button {
            id: button
            Layout.fillHeight: true
            Material.elevation: 1
            Material.accent: Material.Orange
            text: "try again"
            onClicked: WebSocket.open(wsUrl)
        }

        BusyIndicator {
            id: spinner
            Layout.fillHeight: true
            Layout.minimumWidth: height
            Layout.maximumWidth: height
            visible: false
            Material.accent: "white"
        }
    }

    states: [
        State {
            name: "disconnected"
            when: WebSocket.connected === false
        },
        State {
            name: "connecting"
            PropertyChanges {
                target: messageText
                text: "connecting..."
            }
            PropertyChanges {
                target: button
                visible: false
            }
            PropertyChanges {
                target: spinner
                visible: true
                running: true
            }
        },
        State {
            name: "connected"
            when: WebSocket.connected === true
            PropertyChanges {
                target: statusBar
                height: 0
            }
            PropertyChanges {
                target: messageText
                visible: false
            }
        }
    ]

    Behavior on height {NumberAnimation{duration: 220}}

    state: "connecting"
}
