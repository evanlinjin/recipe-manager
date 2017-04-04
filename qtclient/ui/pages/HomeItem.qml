import QtQuick 2.7
import QtQuick.Controls 2.1
import QtQuick.Layouts 1.3
import QtQuick.Controls.Material 2.1

import "../"
import "../components"
import "../components/toolbars"
import "../components/labels"

Item {
    id: homeItem
    anchors.fill: parent

    property int drawerMinWidth: 65
    property int drawerMaxWidth: 220
    property int leftPaneMinWidth: 280
    property int transitionWidth: 700

    RowLayout {
        anchors.fill: parent
        anchors.bottomMargin: statusBar.height
        spacing: 0
        Loader {
            id: loaderOne
            sourceComponent: drawerComponent
            Layout.fillHeight: true
            Layout.preferredWidth: drawerMaxWidth
            Layout.maximumWidth: (loaderThree.status === Loader.Null) ||
                                 (loaderThree.status === Loader.Error) ?
                                     Layout.preferredWidth : drawerMinWidth
        }
        Loader {
            id: loaderTwo
//            sourceComponent: leftComponent
            Layout.fillHeight: true
            Layout.fillWidth: true
            Layout.minimumWidth: leftPaneMinWidth
            Layout.maximumWidth: (loaderThree.status === Loader.Null) ||
                                 (loaderThree.status === Loader.Error) ?
                                     -1 : leftPaneMinWidth
        }
        Loader {
            id: loaderThree
            Layout.fillHeight: true
            Layout.fillWidth: true
            Layout.minimumWidth: 400
            visible: (loaderThree.status !== Loader.Null) &&
                     (loaderThree.status !== Loader.Error)
        }
    }
    StatusBar {
        id: statusBar
        anchors.left: parent.left
        anchors.right: parent.right
        anchors.bottom: parent.bottom
    }

    states: [
        State {
            name: "one_pane"
            when: homeItem.width <= transitionWidth
            PropertyChanges {
                target: loaderOne
                Layout.maximumWidth: 0
                sourceComponent: undefined
                visible: false
            }
            PropertyChanges {
                target: loaderTwo
                visible: (loaderThree.status === Loader.Null) ||
                         (loaderThree.status === Loader.Error)

            }
            PropertyChanges {
                target: loaderThree
                visible: (loaderThree.status !== Loader.Null) &&
                         (loaderThree.status !== Loader.Error)
            }
            PropertyChanges {
                target: drawer
                width: drawerMaxWidth*1.2
                height: homeItem.height
            }
            PropertyChanges {
                target: drawerLoader
                sourceComponent: drawerComponent
            }
        }
    ]

    Drawer {
        id: drawer
        width: 0
        height: 0
        Loader {
            id: drawerLoader
            anchors.fill: parent
        }
        Material.elevation: 0
    }

    Component{id: drawerComponent; DrawerPage{}}
    Component{id: settingsComponent; SettingsPage{}}
    Component{id: aboutComponent; AboutPage{}}

    Component{
        id: rightComponent;
        Rectangle{
            anchors.fill: parent
            color: "black"
            MouseArea {
                anchors.fill: parent
                onClicked: closeRightPane()
            }
        }
    }

    function closeRightPane() {
        loaderThree.sourceComponent = undefined
        drawer.close()
    }

    function openSettings() {
        loaderTwo.sourceComponent = settingsComponent
        drawer.close()
    }

    function openAbout() {
        loaderTwo.sourceComponent = aboutComponent
        drawer.close()
    }
}
