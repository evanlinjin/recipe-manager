#ifndef WEBSOCKETCONNECTION_H
#define WEBSOCKETCONNECTION_H

#include <QObject>
#include <QWebSocket>
#include <QDebug>

class WebSocketConnection : public QObject
{
    Q_OBJECT
    Q_PROPERTY(bool connected READ connected NOTIFY connectedChanged)
public:
    explicit WebSocketConnection(QObject *parent = 0);

    bool connected() const {return m_connected;}

private:
    QWebSocket m_ws;
    bool m_connected;

signals:
    void connectedChanged();
    void msgRecieved(QString);

private slots:
    void onConnected();
    void onDisconnected();
    void onPong(quint64, QByteArray);

    void onError(QAbstractSocket::SocketError);
    void onSslErrors(const QList<QSslError> &errors);

public slots:
    void open(QString v);
    void close();
    void sendMsg(QString v);
    void sendPing(QString v);
};

#endif // WEBSOCKETCONNECTION_H
