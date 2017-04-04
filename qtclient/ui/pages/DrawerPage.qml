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
    anchors.fill: parent
    Material.background: Material.color(Material.Grey, Material.Shade900)
    header: DynamicToolBar {
        component: RowLayout {
            anchors.fill: parent
            IconToolButton {
                id: backButton
                iconSource: "qrc:/ui/icons/back.png"
                ToolTip.text: "back"
                onClicked: drawer.close()
                visible: homeItem.state === "one_pane"
            }
            Spacer{visible: !backButton.visible}
            HeaderLabel {
                text: "Menu"
            }
        }
    }
    ListView {
        anchors.fill: parent
        model: VisualItemModel {
            ItemDelegate {
                text: "hello"
                width: page.width
            }
            ItemDelegate {
                text: "hello"
                width: page.width
            }
            ItemDelegate {
                text: "hello"
                width: page.width
            }
        }

    }
    footer: ListView {
        width: page.width
        height: contentHeight
        model: VisualItemModel {
            DrawerBottomDelegate {
                iconSource: "qrc:/ui/icons/settings.png"
                labelText: "Settings"
                width: page.width
                highlighted: loaderTwo.item ? loaderTwo.item.objectName === "__settings__" : false
                onClicked: openSettings()
            }
            DrawerBottomDelegate {
                iconSource: "qrc:/ui/icons/info.png"
                labelText: "About"
                width: page.width
                highlighted: loaderTwo.item ? loaderTwo.item.objectName === "__about__" : false
                onClicked: openAbout()
            }
        }
    }
}
