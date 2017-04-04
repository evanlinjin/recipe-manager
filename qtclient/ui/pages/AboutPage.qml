import QtQuick 2.7
import QtQuick.Controls 2.1
import QtQuick.Layouts 1.3
import QtQuick.Controls.Material 2.1

import "../components"
import "../components/toolbars"
import "../components/labels"
import "../delegates"

Page {
    objectName: "__about__"
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
                text: "About"
            }
        }
    }
    AboutItem {
        anchors.fill: parent
    }
}
