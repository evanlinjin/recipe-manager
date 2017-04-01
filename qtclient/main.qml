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
    minimumWidth: 420
    minimumHeight: 480
    title: qsTr("Recipe Manager")

    property int maxWidth: 1024
    property string wsUrl: "ws://localhost:8080"

    Loader {
        id: mainLoader
        anchors.fill: parent
        sourceComponent: welcomePage
    }

    Component{id: welcomePage; WelcomePage{}}

    Component.onCompleted: {
        WebSocket.open(wsUrl)
    }
}
