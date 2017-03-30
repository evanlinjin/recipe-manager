#ifndef WEBSOCKETCONNECTION_H
#define WEBSOCKETCONNECTION_H

#include <QObject>
#include <QWebSocket>
#include <QDebug>
#include <QString>
#include <QStringList>

#include "package.h"
#include "encryptor.h"
#include "messagemanager.h"

class WebSocketConnection : public QObject
{
    Q_OBJECT
    Q_PROPERTY(bool connected READ connected NOTIFY connectedChanged)
public:
    explicit WebSocketConnection(QObject *parent = 0);

    bool connected() const {return m_connected;}

    bool sendRequestMessage(const QString &cmd, const QJsonValue &data);
    bool sendResponseMessage(const QJsonObject &reqMsg, const QJsonValue &data);

private:
    QWebSocket m_ws;
    Encryptor* m_enc;
    MessageManager* m_msgs;
    bool m_connected;

    void send(QJsonObject &obj);

signals:
    void connectedChanged();
    void msgRecieved(QJsonObject);

private slots:
    void onConnected();
    void onDisconnected();
    void onReceived(QString data);
    void onPong(quint64, QByteArray);

    void onError(QAbstractSocket::SocketError);
    void onSslErrors(const QList<QSslError> &errors);

public slots:
    void open(QString v);
    void close();

    void sendPing(QString v);
};

#endif // WEBSOCKETCONNECTION_H
