#include "encryptor.h"

Encryptor::Encryptor(QObject *parent) : QObject(parent) {
    memset(m_key, 0, DEF_SIZE);
}

QByteArray Encryptor::encrypt(QByteArray data) {
    return QByteArray(data)
            .toBase64(QByteArray::Base64UrlEncoding | QByteArray::OmitTrailingEquals);
}

QByteArray Encryptor::decrypt(QByteArray data) {
    return QByteArray::fromBase64(
                data, QByteArray::Base64UrlEncoding | QByteArray::OmitTrailingEquals);
}

void Encryptor::setKey(const QByteArray &encKey) {
    qDebug() << "[Encryptor::setKey] GOT KEY:" << encKey;
    QByteArray key = QByteArray::fromBase64(
                encKey, QByteArray::Base64UrlEncoding | QByteArray::OmitTrailingEquals);
    key = QByteArray::fromBase64(
                key, QByteArray::Base64UrlEncoding | QByteArray::OmitTrailingEquals);

    if (key.length() != DEF_SIZE) {
        qDebug() << "[Encryptor::setKey] Invalid key length."
                 << "GOT:" << key.length()
                 << "EXPECTED:" << DEF_SIZE;
    }

    const char* c = key.constData();
    for (qint64 i = 0; i < DEF_SIZE; i++) {
        m_key[i] = (unsigned char)c[i];
    }
}

QByteArray Encryptor::makeKey() {
    unsigned char randBytes[DEF_SIZE];

    QByteArray key(DEF_SIZE, Qt::Uninitialized);
    for(int i = 0; i < DEF_SIZE; i++)
        key[i] = (char)qrand();

    return key.toBase64(QByteArray::Base64UrlEncoding | QByteArray::OmitTrailingEquals);
}

QByteArray Encryptor::makeTimestampedKey() {
    return QByteArray::number(QDateTime::currentMSecsSinceEpoch()).
            append(makeKey());
}

void Encryptor::resetKey() {
    memset(m_key, 0, DEF_SIZE);
}
