import QtQuick 2.7
import QtQuick.Controls 2.0
import QtQuick.Layouts 1.0
import "ui"
import "ui/components"
import "ui/pages"

ApplicationWindow {
    id: mainWindow
    visible: true
    width: 720
    height: 520
    minimumWidth: 420
    minimumHeight: 480
    title: qsTr("Recipe Manager")

    property int maxWidth: 1024

    Loader {
        id: mainLoader
        anchors.fill: parent
        sourceComponent: welcomePage

        states: [
            State {
                name: "loggedIn"
                PropertyChanges {
                    target: mainLoader
                    sourceComponent: homeItem
                }
            }
        ]

        state: Session.sessionID === "" ?
                   "" : "loggedIn"
    }

    Component{id: welcomePage; WelcomePage{}}
    Component{id: homeItem; HomeItem{}}

    Component.onCompleted: {
        WebSocket.open(Session.url)
        Session.onUrlChanged.connect(function(){
            console.log(Session.url)
            WebSocket.open(Session.url)
        })
        Session.onSessionChanged.connect(function(){
            console.log("session changed:", Session.sessionID)
            mainLoader.state = Session.sessionID === "" ?
                        "" : "loggedIn"
        })
    }
}
