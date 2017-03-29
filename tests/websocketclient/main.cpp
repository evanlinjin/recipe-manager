#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QDebug>

#include "../../qtclient/src/websocketpackage.h"
#include "../../qtclient/src/packageencryptor.h"

#include "ws.h"

int main(int argc, char *argv[])
{
    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);
    QGuiApplication app(argc, argv);

    QQmlApplicationEngine engine;
    QQmlContext *rc = engine.rootContext();

//    QJsonObject p;
//    p["name"] = "Evan";
//    p["age"] = 21;
//    auto k = QByteArray("secret");

//    auto pac = WebSocketPackage::MakePackage(p, k);
//    qDebug() << QString::fromLatin1(pac);

//    auto obj = WebSocketPackage::ReadPackage(
//                QByteArray("eyJuYW1lIjoiRXZhbiIsImFnZSI6MjF9.vpnih-GS0IGhQI86XvlGp2dmZhN-NIicJI2Cl7wGUnE"),
//                QByteArray("secret"));
//    qDebug() << QString::fromLatin1(QJsonDocument(obj).toJson(QJsonDocument::Compact));

//    auto key = WebSocketPackage::MakeRandomBytes(16);
//    auto iv1 = WebSocketPackage::MakeRandomBytes(16);

    PackageEncryptor pe;

    pe.setKey("ggw-eExRUaA37UeKIQpP8Q");
    qDebug() << pe.decryptPackage("ZKIC2Cq1B2CeBxekhQFz0O8c0rFB8PHpG_KkKP826DVHlfngj6fLNM_Kc7TrgIaJ");
//    qDebug() << "[KEY]" << pe.makeKey();
//    qDebug() << "[CIPHER]" << pe.encryptPackage("This is a test whether I can read this from golang?");

//    qDebug() << pe.decryptPackage("NBC_Zbs05MqvhomYFP15b-");

    WS* ws = new WS();
    rc->setContextProperty("WS", ws);

    engine.load(QUrl(QLatin1String("qrc:/main.qml")));

    return app.exec();
}
