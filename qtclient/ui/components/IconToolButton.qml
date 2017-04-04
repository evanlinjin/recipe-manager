import QtQuick 2.7
import QtQuick.Controls 2.1
import QtQuick.Layouts 1.3
import QtQuick.Controls.Material 2.1

ToolButton {
    property alias iconSource: icon.source
    ToolTip.text: ""
    onClicked: {}

    anchors.bottom: parent.bottom
    anchors.top: parent.top
    Layout.fillHeight: true
    ToolTip.visible: hovered

    Image {
        id: icon
        anchors.centerIn: parent
    }
}
