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
    m_msgs = new MessageManager(this);
}

bool WebSocketConnection::sendRequestMessage(const QString &cmd, const QJsonValue &data) {
    auto msg = m_msgs->makeRequestMessage(cmd, data);
    if (msg == nullptr) {
        return false;
    }
    sendMsg(*msg);
    return true;
}

bool WebSocketConnection::sendResponseMessage(const QJsonObject &reqMsg, const QJsonValue &data) {
    auto msg = m_msgs->makeResponseMessage(reqMsg, data);
    if (msg == nullptr) {
        return false;
    }
    sendMsg(*msg);
    return true;
}

void WebSocketConnection::sendMsg(QJsonObject &obj) {
    // Make random Signature.
    auto signature = m_enc->makeKey();

    // Make Data into Package with Signature.
    auto package = Package::MakePackage(obj, signature);

    // Encrypt Data and Signature.
    auto encSignature = m_enc->encrypt(signature);
    auto encPackage = m_enc->encrypt(package);

    // Join with dot.
    QByteArray out;
    out.append(encSignature).append('.').append(encPackage);

    qDebug() << "SENDING:" << QString::fromLatin1(out);
    m_ws.sendTextMessage(QString::fromLatin1(out));
}

void WebSocketConnection::onConnected() {
    m_connected = true;
    emit connectedChanged();
}

void WebSocketConnection::onDisconnected() {
    m_connected = false;
    emit connectedChanged();
}

void WebSocketConnection::onMsg(QString data) {
    // Split msg into Signature and Data.
    QStringList split = data.split('.', QString::SkipEmptyParts);
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
    auto msg = Package::ReadPackage(package, signature);
    if (m_msgs->checkIncomingMessage(msg) == false) {
        return;
    }
    emit msgRecieved(msg);
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

void WebSocketConnection::sendPing(QString v) {
    m_ws.ping(v.toLatin1());
}

void WebSocketConnection::onError(QAbstractSocket::SocketError e) {
    qErrnoWarning(e, "[WebSocketConnection] Error:");
}

void WebSocketConnection::onSslErrors(const QList<QSslError> &errors) {
    Q_UNUSED(errors);
    foreach(QSslError e, errors) {
        qInfo() << "[WebSocketConnection] SSL Error:" << e.errorString();
    }
    m_ws.ignoreSslErrors();
}
