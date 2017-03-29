#ifndef PACKAGEENCRYPTOR_H
#define PACKAGEENCRYPTOR_H

#include <cryptopp/aes.h>
#include <cryptopp/randpool.h>
#include <cryptopp/modes.h>
#include <cryptopp/cryptlib.h>
#include <cryptopp/osrng.h>

#include <QObject>
#include <QDebug>
#include <iostream>

class PackageEncryptor : public QObject
{
    Q_OBJECT
public:
    explicit PackageEncryptor(QObject *parent = 0);

private:
    const int DEF_SIZE = 16;
    byte m_key[16];

public slots:
    QByteArray encryptPackage(QByteArray data);
    QByteArray decryptPackage(QByteArray data);

    void setKey(const QByteArray& encKey);
    QByteArray makeKey();
};

#endif // PACKAGEENCRYPTOR_H
