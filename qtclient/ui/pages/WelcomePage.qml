import QtQuick 2.7
import QtQuick.Controls 2.0
import QtQuick.Layouts 1.0
import QtQuick.Controls.Material 2.1
import QtGraphicalEffects 1.0

import "../components"
import "../components/labels"
import "../components/toolbars"

Page {
    id: page
    property int maxContainerWidth: 380

    header: DynamicToolBar {
        ToolButton {
            Image {
                id: name
                source: stack.depth > 1 ?
                            "qrc:/ui/icons/back.png" : "qrc:/ui/icons/contents.png"
                anchors.centerIn: parent
            }
            onClicked: stack.depth > 1 ?
                           stack.pop() : stack.push(aboutItem)
        }
        background: Item{}
    }

    background: Image {
        id: bgImg
        anchors.fill: parent
        source: "qrc:/ui/backgrounds/veges.png"
        fillMode: Image.PreserveAspectCrop
        Rectangle {
            anchors.fill: parent
            opacity: 0.9
            color: Material.background
        }
    }

    StackView {
        id: stack
        anchors.fill: parent
        initialItem: loginItem
    }

    Component {
        id: loginItem
        Item {
            Pane {
                anchors.centerIn: parent
                width: parent.width > maxContainerWidth ?
                           maxContainerWidth : parent.width
                Layout.margins: 20
                background: Item{}
                ColumnLayout {
                    id: container
                    anchors.fill: parent
                    Image {
                        id: icon
                        source: "qrc:/ui/icons/recipemanager.png"
                        Layout.alignment: Layout.Center
                        Layout.maximumWidth: parent.width/2
                        Layout.maximumHeight: width
                    }
                    TextField {
                        placeholderText: "Email"
                        Layout.fillWidth: true
                    }
                    TextField {
                        placeholderText: "Password"
                        Layout.fillWidth: true
                    }
                    Button {
                        Layout.fillWidth: true
                        text: "Login"
                    }
                    Button {
                        Layout.fillWidth: true
                        text: "Create Account"
                        flat: true
                    }
                    Button {
                        Layout.fillWidth: true
                        text: "Advanced Settings"
                        flat: true
                    }

                    Item {
                        Layout.fillHeight: true
                        Layout.fillWidth: true
                    }
                }
            }
        }
    }

    Component{id: aboutItem; AboutItem{}}

    Label {
        text: "image from pexels.com"
        opacity: 0.7
        anchors.bottom: parent.bottom
        anchors.right: parent.right
        anchors.margins: 10
        horizontalAlignment: Text.AlignRight
        font.pixelSize: 10
    }

}
