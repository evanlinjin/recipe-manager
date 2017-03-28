#ifndef PACKAGEENCRYPTOR_H
#define PACKAGEENCRYPTOR_H

#include "QTinyAes/QTinyAes/qtinyaes.h"

#include <QObject>

class PackageEncryptor : public QObject
{
    Q_OBJECT
public:
    explicit PackageEncryptor(QObject *parent = 0);

private:
    QTinyAes m_aes;

public slots:
    QByteArray encryptPackage(const QByteArray& data);
    QByteArray decryptPackage(const QByteArray& data);

    void setKey(const QByteArray& key);
    void setIv(const QByteArray& iv);
};

#endif // PACKAGEENCRYPTOR_H
