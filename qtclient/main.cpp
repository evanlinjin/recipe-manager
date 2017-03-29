#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QDebug>

#include <iostream>
#include <cryptopp/aes.h>

#include "src/package.h"

int main(int argc, char *argv[])
{
    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);
    QGuiApplication app(argc, argv);


    QQmlApplicationEngine engine;
    engine.load(QUrl(QLatin1String("qrc:/main.qml")));

    return app.exec();
}
