#include "websocketconnection.h"

WebSocketConnection::WebSocketConnection(QObject *parent) : QObject(parent) {
    m_connected = false;

    connect(&m_ws, SIGNAL(connected()),
            this, SLOT(onConnected()));

    connect(&m_ws, SIGNAL(disconnected()),
            this, SLOT(onDisconnected()));

    connect(&m_ws, SIGNAL(textMessageReceived(QString)),
            this, SIGNAL(msgRecieved(QString)));

    connect(&m_ws, SIGNAL(error(QAbstractSocket::SocketError)),
            this, SLOT(onError(QAbstractSocket::SocketError)));

    connect(&m_ws, SIGNAL(pong(quint64,QByteArray)),
            this, SLOT(onPong(quint64,QByteArray)));

    typedef void (QWebSocket:: *sslErrorsSignal)(const QList<QSslError> &);
    connect(&m_ws, static_cast<sslErrorsSignal>(&QWebSocket::sslErrors),
            this, &WebSocketConnection::onSslErrors);
}

void WebSocketConnection::onConnected() {
    m_connected = true;
    emit connectedChanged();
}

void WebSocketConnection::onDisconnected() {
    m_connected = false;
    emit connectedChanged();
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
    m_ws.sendTextMessage(v);
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
