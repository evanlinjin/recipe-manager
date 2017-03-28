#ifndef WEBSOCKETPACKAGE_H
#define WEBSOCKETPACKAGE_H

#include <QJsonObject>
#include <QJsonDocument>
#include <QMessageAuthenticationCode>

class WebSocketPackage {
public:

    static QByteArray MakePackage(const QJsonObject &obj, const QByteArray &key);
    static QJsonObject ReadPackage(const QByteArray &data, const QByteArray &key);

    static const QByteArray MakeRandomBytes(const int &size);

};

#endif // WEBSOCKETPACKAGE_H
