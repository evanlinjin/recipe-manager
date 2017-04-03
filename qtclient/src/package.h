#ifndef PACKAGE_H
#define PACKAGE_H

//#ifndef Q_OS_LINUX
////#include <cryptopp/osrng.h>
//#endif //Q_OS_LINUX

//#ifndef Q_OS_ANDROID
//#include "third_party/cryptopp/osrng.h"
//#endif //Q_OS_ANDROID

#include <QJsonObject>
#include <QJsonDocument>
#include <QMessageAuthenticationCode>
#include <QDebug>

class Package {
public:
    static QByteArray MakePackage(const QJsonObject &obj, const QByteArray &key);
    static QJsonObject ReadPackage(const QByteArray &data, const QByteArray &key);
};

#endif //PACKAGE_H
