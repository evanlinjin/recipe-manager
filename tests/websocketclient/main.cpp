#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>

#include "ws.h"

int main(int argc, char *argv[])
{
    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);
    QGuiApplication app(argc, argv);

    QQmlApplicationEngine engine;
    QQmlContext *rc = engine.rootContext();

    WS* ws = new WS();
    rc->setContextProperty("WS", ws);

    engine.load(QUrl(QLatin1String("qrc:/main.qml")));

    return app.exec();
}
