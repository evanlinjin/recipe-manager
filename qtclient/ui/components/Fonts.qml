import QtQuick 2.7

Item {
    id: fonts

    property int scale: 1
    property int h1: 21

    property alias ubuntuRegular: ubuntuRegular.name
    property alias ubuntuCondensed: ubuntuCondensed.name

    FontLoader {
        id: ubuntuRegular
        source: "qrc:/ui/fonts/ubuntu/Ubuntu-R.ttf"
    }
    FontLoader {
        id: ubuntuCondensed
        source: "qrc:/ui/fonts/ubuntu/Ubuntu-C.ttf"
    }
}
