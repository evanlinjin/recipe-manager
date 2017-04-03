#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QStandardPaths>
#include <QFontDatabase>
#include <QDebug>

#include "src/websocketconnection.h"
#include "src/session.h"
#include "src/jsondb.h"

int main(int argc, char *argv[])
{
    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);
    QGuiApplication app(argc, argv);

    QGuiApplication::setApplicationName("recipemanager");
    QGuiApplication::setApplicationDisplayName("Recipe Manager");
    QGuiApplication::setOrganizationName("recipemanager.io");
    QGuiApplication::setOrganizationDomain("recipemanager.io");

    // Configure Appearance.
    QFontDatabase::addApplicationFont("qrc:/ui/fonts/ubuntu/Ubuntu-R.ttf");
    QGuiApplication::setFont(QFont("ubuntu", 12));
    QGuiApplication::setDesktopSettingsAware(false);

    QQmlApplicationEngine engine;
    QQmlContext *rc = engine.rootContext();

    Session* session = new Session();
    rc->setContextProperty("Session", session);

    WebSocketConnection* ws = new WebSocketConnection();
    rc->setContextProperty("WebSocket", ws);

    engine.load(QUrl(QLatin1String("qrc:/main.qml")));

    return app.exec();
}
