#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QDebug>

#include "../../qtclient/src/websocketpackage.h"

#include "ws.h"

int main(int argc, char *argv[])
{
    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);
    QGuiApplication app(argc, argv);

    QQmlApplicationEngine engine;
    QQmlContext *rc = engine.rootContext();

    QJsonObject p;
    p["name"] = "Evan";
    p["age"] = 21;
    auto k = QByteArray("secret");

    auto pac = WebSocketPackage::MakePackage(p, k);
    qDebug() << QString::fromLatin1(pac);

    auto obj = WebSocketPackage::ReadPackage(
                QByteArray("eyJuYW1lIjoiRXZhbiIsImFnZSI6MjF9.vpnih-GS0IGhQI86XvlGp2dmZhN-NIicJI2Cl7wGUnE"),
                QByteArray("secret"));
    qDebug() << QString::fromLatin1(QJsonDocument(obj).toJson(QJsonDocument::Compact));

    WS* ws = new WS();
    rc->setContextProperty("WS", ws);

    engine.load(QUrl(QLatin1String("qrc:/main.qml")));

    return app.exec();
}
