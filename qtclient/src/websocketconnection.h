#ifndef WEBSOCKETCONNECTION_H
#define WEBSOCKETCONNECTION_H

#include <QObject>
#include <QWebSocket>
#include <QDebug>
#include <QString>
#include <QStringList>
#include <QJsonArray>
#include <QTimer>

#include "messagemanager.h"
#include "session.h"

class WebSocketConnection : public QObject
{
    Q_OBJECT
    // (0):not connected. (3):connected. (other):connecting.
    Q_PROPERTY(int connectionStatus READ connectionStatus NOTIFY connectionStatusChanged)
public:
    explicit WebSocketConnection(QObject *parent = 0);
    void setSession(Session* s) {m_session = s;}

    int connectionStatus() const {return m_connectionStatus;}

    QJsonObject *sendRequestMessage(const QString &cmd, const QJsonValue &data);
    QJsonObject *sendResponseMessage(const QJsonObject &reqMsg, const QJsonValue &data);

private:
    QWebSocket m_ws;
    Session* m_session;
    MessageManager* m_msgs;

    QTimer *m_timer, *m_checker;
    int m_connectionStatus;
    bool m_gotPong;
    const qint64 TIMER_INTERVAL_MS = 10 * 1000;
    const qint64 CHECKER_INTERVAL_MS = 5 * 1000;

signals:
    void connectionStatusChanged();
    void networkError();

private slots:
    void onStateChanged(QAbstractSocket::SocketState);
    void onReceived(QByteArray data);

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
    void open(QString v) {m_ws.close(); m_ws.open(QUrl(v));}
    void close() {m_ws.close();}

    /**************************************************************************/
    /* INCOMING/OUTGOING STUFF                                                */
    /**************************************************************************/

private:
    void send(QJsonObject &obj);

    bool ps_new_chef(const MSG::Message &msg);
    bool ps_login(const MSG::Message &msg);
    bool ps_claimSession(const MSG::Message &msg);

public slots:
    int outgoing_newChef(QString email, QString password);
    int outgoing_login(QString email, QString password);
    int outgoing_logout();
    int outgoing_claimSession();

signals:
    void responseTextMessage(int reqId, QString textMsg);
};

#endif // WEBSOCKETCONNECTION_H
