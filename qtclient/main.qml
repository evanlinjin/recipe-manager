import QtQuick 2.7
import QtQuick.Controls 2.0
import QtQuick.Layouts 1.0
import "ui/components"
import "ui/pages"

ApplicationWindow {
    id: main
    visible: true
    width: 640
    height: 480
    minimumWidth: 320
    minimumHeight: 480
    title: qsTr("Recipe Manager")

    property string wsUrl: "ws://localhost:8080"

    Loader {
        id: mainLoader
        anchors.fill: parent
        sourceComponent: welcomePage
    }



    Component{id: welcomePage; WelcomePage{}}

    Component.onCompleted: {
        WebSocket.open("")
    }
}
