#ifndef WEBSOCKETCONNECTION_H
#define WEBSOCKETCONNECTION_H

#include <QObject>
#include <QWebSocket>
#include <QDebug>
#include <QString>
#include <QStringList>
#include <QTimer>

#include "package.h"
#include "encryptor.h"
#include "messagemanager.h"

class WebSocketConnection : public QObject
{
    Q_OBJECT
    // (0):not connected. (3):connected. (other):connecting.
    Q_PROPERTY(int connectionStatus READ connectionStatus NOTIFY connectionStatusChanged)
public:
    explicit WebSocketConnection(QObject *parent = 0);

    int connectionStatus() const {return m_connectionStatus;}

    QJsonObject *sendRequestMessage(const QString &cmd, const QJsonValue &data);
    QJsonObject *sendResponseMessage(const QJsonObject &reqMsg, const QJsonValue &data);

private:
    QWebSocket m_ws;
    Encryptor* m_enc;
    MessageManager* m_msgs;

    void send(QJsonObject &obj);
    void process(const QJsonObject &obj);

    bool ps_handshake(const MSG::Message &msg);

    QTimer *m_timer, *m_checker;
    int m_connectionStatus;
    bool m_gotPong;
    const qint64 TIMER_INTERVAL_MS = 10 * 1000;
    const qint64 CHECKER_INTERVAL_MS = 5 * 1000;

signals:
    void connectionStatusChanged();
    void msgRecieved(QJsonObject);
    void networkError();

private slots:
    void onStateChanged(QAbstractSocket::SocketState);
    void onReceived(QString data);

    void onError(QAbstractSocket::SocketError e) {
        qErrnoWarning(e, "[WebSocketConnection] Error:");
    }

    void onSslErrors(const QList<QSslError> &errors) {
        foreach(QSslError e, errors) {
            qInfo() << "[WebSocketConnection] SSL Error:" << e.errorString();
        }
        m_ws.ignoreSslErrors();
    }

public slots:
    void open(QString v) {m_ws.open(QUrl(v));}
    void close() {m_ws.close();}
};

#endif // WEBSOCKETCONNECTION_H
