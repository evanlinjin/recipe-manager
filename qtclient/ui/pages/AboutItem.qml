import QtQuick 2.7
import QtQuick.Controls 2.0
import QtQuick.Layouts 1.0
import QtQuick.Controls.Material 2.1

Item {
    ListView {
        anchors.fill: parent
        model: VisualItemModel {
            ColumnLayout {
                anchors.fill: parent
                spacing: 40
                Pane {
                    Layout.alignment: Layout.Center
                    Material.elevation: 10
                    Row {
                        spacing: 20
                        Image {
                            source: "qrc:/ui/icons/recipemanager.png"
                            width: 100
                            height: 100
                        }
                        Label {
                            text: "<h3>Recipe Manager</h3>
Manage your recipes.<br>
Source code avaliable on <a href=\"https://github.com/evanlinjin/recipe-manager\">Github</a>."
                            height: 100
                            verticalAlignment: Text.AlignVCenter
                            onLinkActivated: Qt.openUrlExternally(link)
                        }
                        Item {height: 1; width: 1}
                    }
                }

                Row {
                    Layout.alignment: Layout.Center
                    spacing: 20
                    Label {
                        text:"<h3>Built with Qt</h3>
A cross-platform framework.<br>
Source code avaliable <a href=\"http://code.qt.io/cgit/\">here</a>."
                        height: 70
                        verticalAlignment: Text.AlignVCenter
                        onLinkActivated: Qt.openUrlExternally(link)
                    }
                    Image {
                        source: "qrc:/ui/icons/qt-logo.png"
                        width: 100
                        height: 70
                    }
                }

            }
        }
    }
}

// http://code.qt.io/cgit/
// https://github.com/evanlinjin/recipe-manager
