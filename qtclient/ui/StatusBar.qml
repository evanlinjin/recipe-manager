import QtQuick 2.7
import QtQuick.Controls 2.0
import QtQuick.Layouts 1.0
import QtQuick.Controls.Material 2.1

ToolBar {
    id: statusBar
    Material.elevation: 0
    Material.primary: Material.accent
    height: 40
    RowLayout {
        anchors.centerIn: parent
        height: parent.height
        spacing: 10

        BusyIndicator {
            id: spinner
            Layout.fillHeight: true
            Layout.minimumWidth: height
            Layout.maximumWidth: height
            height: 30
            visible: false
            Material.accent: "white"
        }

        Label {
            id: messageText
            Layout.fillHeight: true
            verticalAlignment: Text.AlignVCenter
            text: "disconnected."
            font.bold: true
        }

        Button {
            id: button
            Layout.fillHeight: true
            Material.elevation: 1
            Material.accent: Material.Orange
            text: "try again"
            visible: false
            onClicked: WebSocket.open(Session.url)
        }


    }

    states: [
        State {
            name: "unknown"
            PropertyChanges {
                target: statusBar
                height: 0
            }
            PropertyChanges {
                target: messageText
                visible: false
            }
        },
        State {
            name: "disconnected"
            when: WebSocket.connectionStatus === 0
            PropertyChanges {
                target: button
                visible: true
            }
        },
        State {
            name: "connecting"
            when: WebSocket.connectionStatus === 2
            PropertyChanges {
                target: messageText
                text: "attempting to establish a connection..."
            }
            PropertyChanges {
                target: spinner
                visible: true
                running: true
            }
        },
        State {
            name: "connected"
            when: WebSocket.connectionStatus === 3
            PropertyChanges {
                target: statusBar
                height: 0
                Material.primary: Material.color(Material.Green, Material.Shade600)
            }
            PropertyChanges {
                target: messageText
                visible: false
            }
        }
    ]

    Behavior on height {NumberAnimation{duration: 220}}

    onStateChanged: switch(state) {
                    case "connected":
                        ToolTip.show("Connected!", 3000)
                        break
                    case "disconnected":
                        ToolTip.show("Disconnected.", 3000)
                        break
                    }

    state: "unknown"
}
