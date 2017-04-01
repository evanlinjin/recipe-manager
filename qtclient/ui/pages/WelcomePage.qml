import QtQuick 2.7
import QtQuick.Controls 2.0
import QtQuick.Layouts 1.0
import QtQuick.Controls.Material 2.1
import QtGraphicalEffects 1.0

import "../"
import "../components"
import "../components/labels"
import "../components/toolbars"

Page {
    id: page
    property int maxContainerWidth: 380

    header: DynamicToolBar {
        maxWidth: page.width
        component: RowLayout {
            ToolButton {
                Layout.fillHeight: true
                Image {
                    id: name
                    source: stack.depth > 1 ?
                                "qrc:/ui/icons/back.png" : "qrc:/ui/icons/info.png"
                    anchors.centerIn: parent
                }
                onClicked: stack.depth > 1 ?
                               stack.pop() : stack.push(aboutItem)
            }

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

    footer: StatusBar{}

    Component {
        id: loginItem
        Item {
            Pane {
                anchors.centerIn: parent
                anchors.verticalCenterOffset: -30
                width: parent.width > maxContainerWidth ?
                           maxContainerWidth : parent.width
                Layout.margins: 20
                background: Item{}
                ColumnLayout {
                    id: container
                    anchors.left: parent.left
                    anchors.right: parent.right
                    Image {
                        id: icon
                        source: "qrc:/ui/icons/recipemanager.png"
                        Layout.alignment: Layout.Center
                        Layout.maximumWidth: parent.width/2
                        Layout.maximumHeight: width
                        opacity: enabled ? 1 : 0.3
                        enabled: WebSocket.connected
                    }
                    TextField {
                        id: emailField
                        placeholderText: "Email"
                        Layout.fillWidth: true
                        enabled: WebSocket.connected
                        inputMethodHints: Qt.ImhEmailCharactersOnly | Qt.ImhLowercaseOnly | Qt.ImhNoAutoUppercase
                        onAccepted: passwordField.forceActiveFocus()
                    }
                    TextField {
                        id: passwordField
                        placeholderText: "Password"
                        Layout.fillWidth: true
                        enabled: WebSocket.connected
                        echoMode: TextInput.Password
                        onAccepted: {}
                    }
                    Button {
                        id: loginButton
                        Layout.fillWidth: true
                        text: "Login"
                        enabled: WebSocket.connected
                    }
                    Button {
                        Layout.fillWidth: true
                        text: "Create Account"
                        flat: true
                        enabled: WebSocket.connected
                    }
                    Button {
                        Layout.fillWidth: true
                        text: "Advanced Settings"
                        flat: true
                        onClicked: stack.push(changeURLItem)
                    }
                    Item {
                        Layout.fillHeight: true
                        Layout.fillWidth: true
                    }
                    Component.onCompleted: emailField.forceActiveFocus()
                }
            }
        }
    }

    Component{
        id: changeURLItem
        Item {
            Pane {
                anchors.centerIn: parent
                anchors.verticalCenterOffset: -30
                width: parent.width > maxContainerWidth ?
                           maxContainerWidth : parent.width
                Layout.margins: 20
                background: Item{}
                ColumnLayout {
                    id: container
                    anchors.left: parent.left
                    anchors.right: parent.right
                    Label {
                        Layout.fillWidth: true
                        text: "Change Connection URL:"
                    }
                    TextField {
                        id: urlInput
                        Layout.fillWidth: true
                        text: wsUrl
                        onAccepted: {
                            WebSocket.close()
                            WebSocket.open(urlInput.text)
                            stack.pop()
                        }
                        Component.onCompleted: {
                            selectAll()
                            forceActiveFocus()
                        }
                    }
                    Button {
                        id: changeButton
                        Layout.fillWidth: true
                        text: "Change"
                        onClicked: {
                            WebSocket.close()
                            WebSocket.open(urlInput.text)
                            stack.pop()
                        }
                    }
                    Button {
                        Layout.fillWidth: true
                        flat: true
                        text: "Cancel"
                        onClicked: stack.pop()
                    }
                }
            }
        }
    }

    Component{id: aboutItem; AboutItem{}}

    Label {
        text: "image from pexels.com"
        opacity: 0.5
        anchors.bottom: parent.bottom
        anchors.right: parent.right
        anchors.margins: 10
        horizontalAlignment: Text.AlignRight
        font.pixelSize: 10
    }

}
