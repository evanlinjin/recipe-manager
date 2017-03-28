#include "ws.h"

WS::WS(QObject *parent) : QObject(parent)
{
    m_connected = false;

    connect(&m_ws, SIGNAL(connected()), this, SLOT(onConnected()));
    connect(&m_ws, SIGNAL(disconnected()), this, SLOT(onDisconnected()));
    connect(&m_ws, SIGNAL(textMessageReceived(QString)), this, SIGNAL(msgRecieved(QString)));
    connect(&m_ws, SIGNAL(error(QAbstractSocket::SocketError)), this, SLOT(onError(QAbstractSocket::SocketError)));
    connect(&m_ws, SIGNAL(pong(quint64,QByteArray)), this, SLOT(onPong(quint64,QByteArray)));

    typedef void (QWebSocket:: *sslErrorsSignal)(const QList<QSslError> &);
    connect(&m_ws, static_cast<sslErrorsSignal>(&QWebSocket::sslErrors), this, &WS::onSslErrors);

}

void WS::onConnected() {
    m_connected = true;
    emit connectedChanged();
}

void WS::onDisconnected() {
    m_connected = false;
    emit connectedChanged();
}

void WS::onPong(quint64 i, QByteArray a) {
    qInfo() << "I got a pong message: [" << i << "] [" << a << "]";
}

void WS::open(QString v) {
    m_ws.open(QUrl(v));
}

void WS::close() {
    m_ws.close();
}

void WS::sendMsg(QString v) {
    m_ws.sendTextMessage(v);
}

void WS::sendPing(QString v) {
    m_ws.ping(v.toLatin1());
}

void WS::onError(QAbstractSocket::SocketError e) {
    qErrnoWarning(e, "Websocket error");
}

void WS::onSslErrors(const QList<QSslError> &errors) {
    Q_UNUSED(errors);
    m_ws.ignoreSslErrors();
}
