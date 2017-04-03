#ifndef ENCRYPTOR_H
#define ENCRYPTOR_H

#include <QObject>
#include <QDebug>
#include <QDateTime>

#include "package.h"

class Encryptor : public QObject
{
    Q_OBJECT
public:
    explicit Encryptor(QObject *parent = 0);
    static const qint64 DEF_SIZE = 16;

private:
    unsigned char m_key[16];

public slots:
    QByteArray encrypt(QByteArray data);
    QByteArray decrypt(QByteArray data);

    void setKey(const QByteArray& encKey);
    QByteArray makeKey();
    QByteArray makeTimestampedKey();
    void resetKey();
};

#endif // ENCRYPTOR_H
