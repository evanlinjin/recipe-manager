#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>

#include "src/websocketconnection.h"

#include "src/package.h"

int main(int argc, char *argv[])
{
    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);
    QGuiApplication app(argc, argv);


    QQmlApplicationEngine engine;
    QQmlContext *rc = engine.rootContext();

    WebSocketConnection* ws = new WebSocketConnection();
    rc->setContextProperty("WebSocket", ws);

    engine.load(QUrl(QLatin1String("qrc:/main.qml")));

    return app.exec();
}
