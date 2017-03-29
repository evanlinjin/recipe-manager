#ifndef PACKAGEENCRYPTOR_H
#define ENCRYPTOR_H

#include <cryptopp/aes.h>
#include <cryptopp/randpool.h>
#include <cryptopp/modes.h>
#include <cryptopp/cryptlib.h>
#include <cryptopp/osrng.h>

#include <QObject>
#include <QDebug>
//#include <iostream>

#include "package.h"

class Encryptor : public QObject
{
    Q_OBJECT
public:
    explicit Encryptor(QObject *parent = 0);
    static const int DEF_SIZE = 16;

private:
    byte m_key[16];

public slots:
    QByteArray encrypt(QByteArray data);
    QByteArray decrypt(QByteArray data);

    void setKey(const QByteArray& encKey);
    QByteArray makeKey();
//    QByteArray getKey();
};

#endif // ENCRYPTOR_H
