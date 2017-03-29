#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QDebug>

//#include "../../qtclient/src/package.h"
//#include "../../qtclient/src/encryptor.h"

#include "../../qtclient/src/websocketconnection.h"

int main(int argc, char *argv[])
{
    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);
    QGuiApplication app(argc, argv);

    QQmlApplicationEngine engine;
    QQmlContext *rc = engine.rootContext();

    WebSocketConnection* ws = new WebSocketConnection();
    rc->setContextProperty("WS", ws);

    /* BEGIN TEST */
//    QJsonObject obj;
//    obj["name"] = "Evan Lin";
//    obj["age"] = 21;

//    Encryptor enc;
//    qDebug() << "ENCRYPTED:" << enc.encrypt(Package::MakePackage(obj, ""));
//    return 0;
    /* BEGIN TEST */



    engine.load(QUrl(QLatin1String("qrc:/main.qml")));
    return app.exec();
}
