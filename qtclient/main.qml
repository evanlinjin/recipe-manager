import QtQuick 2.7
import QtQuick.Controls 2.0
import QtQuick.Layouts 1.0
import "ui"
import "ui/components"
import "ui/pages"

ApplicationWindow {
    id: main
    visible: true
    width: 720
    height: 520
//    minimumWidth: 420
//    minimumHeight: 480
    title: qsTr("Recipe Manager")

    property int maxWidth: 1024
//    property string wsUrl: "ws://localhost:8080/ws"
    property string wsUrl: "ws://192.168.20.29:8080/ws"

    Loader {
        id: mainLoader
        anchors.fill: parent
        sourceComponent: welcomePage

        states: [
            State {
                name: "loggedIn"
                PropertyChanges {
                    target: mainLoader
                    sourceComponent: test
                }
            }
        ]

        state: Session.sessionID === "" ?
                   "" : "loggedIn"
    }

    Component{id: welcomePage; WelcomePage{}}

    Component{
        id: test
        Item {
            anchors.fill: parent
            Label {
                anchors.centerIn: parent
                text: "HELLO WORLD!\nYOU ARE LOGGED IN!!!"
            }
        }
    }

    Component.onCompleted: {
        WebSocket.open(wsUrl)
    }
}
