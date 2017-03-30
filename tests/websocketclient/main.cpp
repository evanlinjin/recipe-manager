#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QDebug>

#include "../../qtclient/src/websocketconnection.h"
#include "src/simplemessenger.h"

int main(int argc, char *argv[])
{
    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);
    QGuiApplication app(argc, argv);

    QQmlApplicationEngine engine;
    QQmlContext *rc = engine.rootContext();

    WebSocketConnection* ws = new WebSocketConnection();
    rc->setContextProperty("WS", ws);

    SimpleMessenger* sm = new SimpleMessenger(ws);
    rc->setContextProperty("Messenger", sm);

    engine.load(QUrl(QLatin1String("qrc:/main.qml")));
    return app.exec();
}
