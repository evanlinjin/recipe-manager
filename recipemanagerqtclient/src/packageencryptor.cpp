#include "packageencryptor.h"

PackageEncryptor::PackageEncryptor(QObject *parent) : QObject(parent) {
}

QByteArray PackageEncryptor::encryptPackage(const QByteArray& data) {
    return m_aes.encrypt(data)
            .toBase64(QByteArray::Base64UrlEncoding |
                      QByteArray::OmitTrailingEquals);
}

QByteArray PackageEncryptor::decryptPackage(const QByteArray& data) {
    auto out = QByteArray::fromBase64(data,
                                      QByteArray::Base64UrlEncoding |
                                      QByteArray::OmitTrailingEquals);
    return m_aes.decrypt(out);
}

void PackageEncryptor::setKey(const QByteArray& key) {
    m_aes.setKey(key);
}

void PackageEncryptor::setIv(const QByteArray &iv) {
    m_aes.setIv(iv);
}
