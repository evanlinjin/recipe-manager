#include "websocketconnection.h"

WebSocketConnection::WebSocketConnection(QObject *parent) :
    QObject(parent), m_connected(false) {

    m_enc = new Encryptor(this);
    m_msgs = new MessageManager(this);

    m_timer = new QTimer(this);
    connect(m_timer, SIGNAL(timeout()), &m_ws, SLOT(ping()));

    m_checker = new QTimer(this);
    m_checker->setSingleShot(true);
    connect(&m_ws, SIGNAL(pong(quint64,QByteArray)), m_checker, SLOT(stop()));
    connect(m_checker, SIGNAL(timeout()), this, SIGNAL(networkError()));

    connect(&m_ws, SIGNAL(connected()), this, SLOT(onConnected()));
    connect(&m_ws, SIGNAL(disconnected()), this, SLOT(onDisconnected()));

    connect(&m_ws, SIGNAL(pong(quint64,QByteArray)), this, SLOT(onPong()));

    connect(&m_ws, SIGNAL(textMessageReceived(QString)),
            this, SLOT(onReceived(QString)));

    connect(&m_ws, SIGNAL(error(QAbstractSocket::SocketError)),
            this, SLOT(onError(QAbstractSocket::SocketError)));

    typedef void (QWebSocket:: *sslErrorsSignal)(const QList<QSslError> &);
    connect(&m_ws, static_cast<sslErrorsSignal>(&QWebSocket::sslErrors),
            this, &WebSocketConnection::onSslErrors);
}

QJsonObject *WebSocketConnection::sendRequestMessage(const QString &cmd, const QJsonValue &data) {
    auto msg = m_msgs->makeRequestMessage(cmd, data);
    if (msg == nullptr) {
        return msg;
    }
    send(*msg);
    return msg;
}

QJsonObject *WebSocketConnection::sendResponseMessage(const QJsonObject &reqMsg, const QJsonValue &data) {
    auto msg = m_msgs->makeResponseMessage(reqMsg, data);
    if (msg == nullptr) {
        return msg;
    }
    send(*msg);
    return msg;
}

void WebSocketConnection::send(QJsonObject &obj) {
    qInfo() << "[WebSocketConnection::send] (plain)" << obj;
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

    qInfo() << "[WebSocketConnection::send] (ciphr)" << out;
    m_ws.sendTextMessage(QString::fromLatin1(out));
}

void WebSocketConnection::onConnected() {
    m_connected = true;
    emit connectedChanged();
    m_timer->start(TIMER_INTERVAL_MS);
}

void WebSocketConnection::onDisconnected() {
    m_connected = false;
    emit connectedChanged();
    m_msgs->resetIds();
    m_timer->stop();
}

void WebSocketConnection::onReceived(QString data) {
    qInfo() << "[WebSocketConnection::onReceived] (ciphr)" << data;

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
    qInfo() << "[WebSocketConnection::onReceived] (plain)" << msg;
    if (m_msgs->checkIncomingMessage(msg) == false) {
        return;
    }
    emit msgRecieved(msg);
}

