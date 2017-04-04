import QtQuick 2.7
import QtQuick.Controls 2.1
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
    property int maxWidth: 1024
    signal changedHeader(string txt);

    header: DynamicToolBar {
        component: RowLayout {
            IconToolButton {
                iconSource: stack.depth > 1 ?
                                "qrc:/ui/icons/back.png" : "qrc:/ui/icons/info.png"
                ToolTip.text: stack.depth > 1 ?
                                  "Back" : "About"
                onClicked: stack.depth > 1 ?
                               stack.pop() : stack.push(aboutItem)
            }
            HeaderLabel {
                id: headerLabel
                Component.onCompleted: page.onChangedHeader.connect(function(txt){
                    text = txt
                })
            }
        }
        background: Item{}
    }

    StackView {
        id: stack
        anchors.fill: parent
        initialItem: loginItem
        onCurrentItemChanged: changedHeader(currentItem.objectName)
        Component.onCompleted: changedHeader(currentItem.objectName)
    }

    Component {
        id: loginItem
        Item {
            objectName: "Login"
            Pane {
                id: pane
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
                    anchors.margins: 10
                    Image {
                        id: icon
                        source: "qrc:/ui/icons/recipemanager.png"
                        Layout.alignment: Layout.Center
                        Layout.maximumWidth: parent.width/2
                        Layout.maximumHeight: width
                        Layout.preferredHeight: 100
                        Layout.preferredWidth: 100
                        opacity: enabled ? 1 : 0.3
                        enabled: WebSocket.connectionStatus === 3
                    }
                    TextField {
                        id: urlField
                        placeholderText: "URL"
                        Layout.fillWidth: true
                        inputMethodHints: Qt.ImhNoAutoUppercase | Qt.ImhUrlCharactersOnly
                        onAccepted: {emailField.forceActiveFocus(); Session.url = text}
                        onFocusChanged: Session.url = text
                        Component.onCompleted: text = Session.url
                    }
                    TextField {
                        id: emailField
                        placeholderText: "Email"
                        Layout.fillWidth: true
                        inputMethodHints: Qt.ImhEmailCharactersOnly | Qt.ImhLowercaseOnly | Qt.ImhNoAutoUppercase
                        onAccepted: passwordField.forceActiveFocus()
                        onTextChanged: checkSubmitOkay()
                        visible: WebSocket.connectionStatus === 3
                    }
                    TextField {
                        id: passwordField
                        placeholderText: "Password"
                        Layout.fillWidth: true
                        echoMode: TextInput.Password
                        onAccepted: loginButton.enabled ?
                                        send() : {}
                        onTextChanged: checkSubmitOkay()
                        visible: WebSocket.connectionStatus === 3
                    }
                    Button {
                        id: loginButton
                        Layout.fillWidth: true
                        text: "Login"
                        enabled: false
                        onClicked: send()
                    }
                    Button {
                        Layout.fillWidth: true
                        text: "Join a server"
                        flat: true
                        onClicked: stack.push(createAccountItem)
                    }
                    Item {
                        Layout.fillHeight: true
                        Layout.fillWidth: true
                    }
                }
            }
            BusyIndicator {
                id: busySpinner
                anchors.centerIn: parent
                running: false
            }
            states: [
                State {
                    name: "processing"
                    PropertyChanges {
                        target: pane
                        enabled: false
                    }
                    PropertyChanges {
                        target: busySpinner
                        running: true
                    }
                }
            ]
            function checkSubmitOkay() {
                loginButton.enabled =
                        /[^\s@]+@[^\s@]+\.[^\s@]+/.test(emailField.text) &&
                        passwordField.text.length >= 6
            }
            function send() {
                outgoingReqId = WebSocket.outgoing_login(emailField.text,
                                                         passwordField.text)
                state = "processing"
            }
        }
    }

    Component {
        id: createAccountItem
        Item {
            objectName: "Join a Server"
            Pane {
                id: pane
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
                    anchors.margins: 10
                    TextField {
                        id: urlField
                        placeholderText: "URL"
                        Layout.fillWidth: true
                        inputMethodHints: Qt.ImhNoAutoUppercase | Qt.ImhUrlCharactersOnly
                        onAccepted: emailField.forceActiveFocus()
                        onFocusChanged: Session.url = text
                        Component.onCompleted: text = Session.url
                    }
                    Label {
                        id: invalidEmail
                        Layout.fillWidth: true
                        horizontalAlignment: Text.AlignRight
                        text: "invalid email"
                        color: Material.color(Material.Red, Material.Shade600)
                        visible: false
                    }
                    TextField {
                        id: emailField
                        placeholderText: "Email"
                        Layout.fillWidth: true
                        inputMethodHints: Qt.ImhEmailCharactersOnly | Qt.ImhLowercaseOnly | Qt.ImhNoAutoUppercase
                        onAccepted: passwordField.forceActiveFocus()
                        onTextChanged: checkEmailValid()
                    }
                    Label {
                        id: passwordTooShort
                        Layout.fillWidth: true
                        horizontalAlignment: Text.AlignRight
                        text: "password too short"
                        color: Material.color(Material.Red, Material.Shade600)
                        visible: false
                    }
                    TextField {
                        id: passwordField
                        placeholderText: "Password"
                        Layout.fillWidth: true
                        echoMode: TextInput.Password
                        onAccepted: confirmPasswordField.forceActiveFocus()
                        onTextChanged: checkPasswordLength()
                    }
                    Label {
                        id: passwordsDontMatch
                        Layout.fillWidth: true
                        horizontalAlignment: Text.AlignRight
                        text: "passwords don't match"
                        color: Material.color(Material.Red, Material.Shade600)
                        visible: false
                    }
                    TextField {
                        id: confirmPasswordField
                        placeholderText: "Confirm Password"
                        Layout.fillWidth: true
                        echoMode: TextInput.Password
                        onAccepted: send()
                        onTextChanged: checkMatchingPasswords()
                    }
                    Button {
                        id: submit
                        text: "Submit"
                        Layout.fillWidth: true
                        visible: WebSocket.connectionStatus === 3
                        enabled: false
                        onClicked: send()
                    }
                    Button {
                        id: cancel
                        text: "Cancel"
                        Layout.fillWidth: true
                        flat: true
                        onClicked: stack.pop()
                    }
                }
            }
            BusyIndicator {
                id: busySpinner
                anchors.centerIn: parent
                running: false
            }
            states: [
                State {
                    name: "processing"
                    PropertyChanges {
                        target: pane
                        enabled: false
                    }
                    PropertyChanges {
                        target: busySpinner
                        running: true
                    }
                }
            ]
            function checkEmailValid() {
                invalidEmail.visible =
                        !/[^\s@]+@[^\s@]+\.[^\s@]+/.test(emailField.text)
                checkOkaySubmit()
            }
            function checkPasswordLength() {
                passwordTooShort.visible =
                        passwordField.text.length < 6
                checkOkaySubmit()
            }
            function checkMatchingPasswords() {
                passwordsDontMatch.visible =
                        passwordField.text !== confirmPasswordField.text
                checkOkaySubmit()
            }
            function checkOkaySubmit() {
                submit.enabled = !invalidEmail.visible &&
                        !passwordTooShort.visible &&
                        !passwordsDontMatch.visible &&
                        emailField.text.length !== 0 &&
                        passwordField.text.length !== 0 &&
                        confirmPasswordField.text.length !== 0
            }
            function send() {
                outgoingReqId = WebSocket.outgoing_newChef(emailField.text,
                                                           passwordField.text)
                state = "processing"
            }
        }
    }

    Dialog {
        id: msgDialog
        x: (parent.width - width)/2
        y: (parent.height - height)/2
        parent: ApplicationWindow.overlay
        title: ("Message")
        modal: true
        Label {
            id: msgDialogBody
            anchors.centerIn: parent
            elide: Text.ElideRight
            wrapMode: Text.Wrap
            width: parent.width
            onTextChanged: console.log("text changed to: ", text)
        }
        standardButtons: Dialog.Ok
        function setTextAndOpen(txt) {msgDialogBody.text = txt; open()}
    }

    Component{id: aboutItem; AboutItem{objectName: "About"}}


    property int outgoingReqId: -1

    function response(reqId, txtMsg) {
//        if (reqId !== outgoingReqId)
//            return;
        while (stack.depth > 1)
            stack.pop()
        msgDialog.setTextAndOpen(txtMsg)
        stack.clear()
        stack.push(loginItem)
    }

    Component.onCompleted: WebSocket.onResponseTextMessage.connect(response)

}
