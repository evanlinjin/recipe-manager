#ifndef PACKAGE_H
#define PACKAGE_H

#include <cryptopp/osrng.h>

#include <QJsonObject>
#include <QJsonDocument>
#include <QMessageAuthenticationCode>
#include <QDebug>

class Package {
public:

    static QByteArray MakePackage(const QJsonObject &obj, const QByteArray &key);
    static QJsonObject ReadPackage(const QByteArray &data, const QByteArray &key);

    static const QByteArray MakeRandomBytes(const int &size);

};

#endif //PACKAGE_H
