#include "websocketconnection.h"

WebSocketConnection::WebSocketConnection(QObject *parent) : QObject(parent) {
    m_connected = false;

    connect(&m_ws, SIGNAL(connected()),
            this, SLOT(onConnected()));

    connect(&m_ws, SIGNAL(disconnected()),
            this, SLOT(onDisconnected()));

    connect(&m_ws, SIGNAL(textMessageReceived(QString)),
            this, SLOT(onMsg(QString)));

    connect(&m_ws, SIGNAL(error(QAbstractSocket::SocketError)),
            this, SLOT(onError(QAbstractSocket::SocketError)));

    connect(&m_ws, SIGNAL(pong(quint64,QByteArray)),
            this, SLOT(onPong(quint64,QByteArray)));

    typedef void (QWebSocket:: *sslErrorsSignal)(const QList<QSslError> &);
    connect(&m_ws, static_cast<sslErrorsSignal>(&QWebSocket::sslErrors),
            this, &WebSocketConnection::onSslErrors);
    m_enc = new Encryptor(this);
}

void WebSocketConnection::onConnected() {
    m_connected = true;
    emit connectedChanged();
}

void WebSocketConnection::onDisconnected() {
    m_connected = false;
    emit connectedChanged();
}

void WebSocketConnection::onMsg(QString msg) {
    // Split msg into Signature and Data.
    QStringList split = msg.split('.', QString::SkipEmptyParts);
    if (split.length() != 2) {
        qDebug("Expected %d parts, got %d", 2, split.length());
        return;
    }

    auto encSignature = split.at(0).toLatin1();
    auto encPackage = split.at(1).toLatin1();

    // Decrypt Signature and Data.
    auto signature = m_enc->decrypt(encSignature);
    auto package = m_enc->decrypt(encPackage);

    // Vertify Data with Signature.
    QJsonObject obj = Package::ReadPackage(package, signature);
    emit msgRecieved(obj["msg"].toString());
}

void WebSocketConnection::onPong(quint64 i, QByteArray a) {
    qInfo() << "I got a pong message: [" << i << "] [" << a << "]";
}

void WebSocketConnection::open(QString v) {
    m_ws.open(QUrl(v));
}

void WebSocketConnection::close() {
    m_ws.close();
}

void WebSocketConnection::sendMsg(QString v) {
    // Make random Signature.
    auto signature = m_enc->makeKey();
    qDebug() << "RAW SIGNATURE:" << signature;

    // Make Data into Package with Signature.
    QJsonObject obj;
    obj["msg"] = v;
    auto package = Package::MakePackage(obj, signature);
    qDebug() << "RAW PACKAGE:" << package;

    // Encrypt Data and Signature.
    auto encSignature = m_enc->encrypt(signature);
    qDebug() << "ENC SIGNATURE:" << encSignature;
    auto encPackage = m_enc->encrypt(package);
    qDebug() << "ENC PACKAGE:" << encPackage;

    // Join with dot.
    QByteArray out;
    out.append(encSignature).append('.').append(encPackage);

    qDebug() << "SENDING:" << QString::fromLatin1(out);
    m_ws.sendTextMessage(QString::fromLatin1(out));
}

void WebSocketConnection::sendPing(QString v) {
    m_ws.ping(v.toLatin1());
}

void WebSocketConnection::onError(QAbstractSocket::SocketError e) {
    qErrnoWarning(e, "Websocket error");
}

void WebSocketConnection::onSslErrors(const QList<QSslError> &errors) {
    Q_UNUSED(errors);
    foreach(QSslError e, errors) {
        qInfo() << "SSL error:" << e.errorString();
    }
    m_ws.ignoreSslErrors();
}
