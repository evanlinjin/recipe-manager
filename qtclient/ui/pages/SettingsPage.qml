import QtQuick 2.7
import QtQuick.Controls 2.1
import QtQuick.Layouts 1.3
import QtQuick.Controls.Material 2.1

import "../components"
import "../components/toolbars"
import "../components/labels"
import "../delegates"

Page {
    id: page
    objectName: "__settings__"
    header: DynamicToolBar {
        component: RowLayout {
            IconToolButton {
                id: menuButton
                iconSource: "qrc:/ui/icons/contents.png"
                ToolTip.text: "Menu"
                visible: !loaderOne.visible
                onClicked: drawer.open()
            }
            Spacer{visible: !menuButton.visible}
            HeaderLabel {
                text: "Settings"
            }
        }
    }

    ColumnLayout {
        id: container
        anchors.fill: parent
        ListView {
            id: listView
            Layout.fillHeight: true
            Layout.fillWidth: true
            Layout.alignment: Layout.Center
            Layout.maximumWidth: maxWidth
            model: VisualItemModel {
                ItemDelegate {
                    text: "Log out"
                    width: listView.width
                    onClicked: WebSocket.outgoing_logout()
                }
            }
        }
    }
}
