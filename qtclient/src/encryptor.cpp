#include "encryptor.h"

Encryptor::Encryptor(QObject *parent) : QObject(parent) {
    memset(m_key, 0, DEF_SIZE);
}

QByteArray Encryptor::encrypt(QByteArray data) {

    // Pad data with 0's. TODO: OPTIMIZE!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
    while (data.length() % DEF_SIZE != 0) {
        data.append(" ");
    }

    // Get PlainText Data.
    char plainText[data.length()];
    const char* c = data.constData();
    for (int i = 0; i < data.length(); i++) {
        plainText[i] = c[i];
    }

    // Get Message Length.
    int msgLen = (sizeof(plainText)/sizeof(*plainText));

    // Make random iv.
    CryptoPP::AutoSeededRandomPool rnd;
    byte iv[DEF_SIZE];
    rnd.GenerateBlock(iv, DEF_SIZE);

    // Encrypt.
    char cipherText[data.length()];
    CryptoPP::CBC_Mode<CryptoPP::AES>::Encryption e(m_key, DEF_SIZE, iv);
    e.ProcessData((byte*)cipherText, (byte*)plainText, msgLen);

    QByteArray ivPart;
    ivPart.reserve(DEF_SIZE);
    for (int i = 0; i < DEF_SIZE; i++) {
        ivPart[i] = (char)iv[i];
    }

    return QByteArray(cipherText, msgLen)
            .prepend(ivPart)
            .toBase64(QByteArray::Base64UrlEncoding | QByteArray::OmitTrailingEquals);
}

QByteArray Encryptor::decrypt(QByteArray data) {
    QByteArray rawData = QByteArray::fromBase64(
                data, QByteArray::Base64UrlEncoding | QByteArray::OmitTrailingEquals);

    // Split into iv\cipherText parts.
    if (rawData.length() < DEF_SIZE) {
        qDebug("cipherText too short at %d, expecting >= %d", rawData.length(), DEF_SIZE);
        return QByteArray();
    }

    const char* c = rawData.constData();
    const int msgLen = rawData.length() - DEF_SIZE;

    // Get iv part.
    byte iv[DEF_SIZE];
    for (int i = 0; i < DEF_SIZE; i++) {
        iv[i] = (byte)c[i];
    }

    // Get cipher part.
    char cipherText[msgLen];
    for (int i = 0; i < msgLen; i++) {
        cipherText[i] = (byte)c[i+DEF_SIZE];
    }

    // Check if cipher part satisfies block size.
    if (msgLen % DEF_SIZE != 0) {
        qDebug("cipherText should be a multiple of %d, got %d", DEF_SIZE, msgLen);
        return QByteArray();
    }

    // Decrypt.
    char plainText[msgLen];
    CryptoPP::CBC_Mode<CryptoPP::AES>::Decryption e(m_key, DEF_SIZE, iv);
    e.ProcessData((byte*)plainText, (byte*)cipherText, msgLen);

    return QByteArray(plainText, msgLen).trimmed();
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
    for (int i = 0; i < DEF_SIZE; i++) {
        m_key[i] = (byte)c[i];
    }
}

QByteArray Encryptor::makeKey() {
    QByteArray key;
    key.reserve(DEF_SIZE);

    CryptoPP::AutoSeededRandomPool rnd;
    rnd.GenerateBlock((byte*)key.data(), DEF_SIZE);

    return key.toBase64(QByteArray::Base64UrlEncoding | QByteArray::OmitTrailingEquals);
}

//QByteArray Encryptor::getKey() {
//    QByteArray key;
//    key.reserve(DEF_SIZE);
//    for (int i = 0; i < DEF_SIZE; i++) {
//        key[i] = (char)key[i];
//    }
//}
