import QtQuick 2.7
import QtQuick.Controls 2.1
import QtQuick.Layouts 1.3
import QtQuick.Controls.Material 2.1

import "../components"
import "../components/toolbars"
import "../components/labels"

ItemDelegate {
    id: item
    property string iconSource: "qrc:/ui/icons/settings.png"
    property string labelText: "Settings"
    Loader {
        anchors.fill: parent
        sourceComponent: width > 100 ?
                             fullLayout : smallLayout
    }

    Component {
        id: fullLayout
        RowLayout {
            anchors.fill: parent
            anchors.leftMargin: 20
            Image {
                id: icon
                source: iconSource
                anchors.verticalCenter: parent.verticalCenter
            }
            Spacer{}
            Label {
                id: label
                text: labelText
                anchors.verticalCenter: parent.verticalCenter
                Layout.fillWidth: true
            }
            visible: width > 100
        }
    }

    Component {
        id: smallLayout
        Item {
            Image {
                id: icon
                source: iconSource
                anchors.centerIn: parent
            }
            visible: width <= 100
        }
    }
}
