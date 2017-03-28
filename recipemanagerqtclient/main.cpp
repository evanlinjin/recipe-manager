#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QDebug>

#include <iostream>
#include <cryptopp/aes.h>

#include "src/websocketpackage.h"

int main(int argc, char *argv[])
{
    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);
    QGuiApplication app(argc, argv);



//    byte key[CryptoPP::AES::DEFAULT_KEYLENGTH];
//    byte iv[CryptoPP::AES::IV_LENGTH];
    const int count = 16;

    QByteArray key = WebSocketPackage::MakeRandomBytes(count);
    QByteArray iv = WebSocketPackage::MakeRandomBytes(count);

    byte key2[key.length()];
    for (int i = 0; i < key.length(); i++) {
        key2[i] = (byte)key.constData()[i];
    }

    std::cout << key2 << std::endl;

    QQmlApplicationEngine engine;
    engine.load(QUrl(QLatin1String("qrc:/main.qml")));

    return app.exec();
}
