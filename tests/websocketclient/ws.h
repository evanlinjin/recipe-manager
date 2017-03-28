#ifndef WS_H
#define WS_H

#include <QObject>
#include <QWebSocket>
#include <QDebug>

class WS : public QObject
{
    Q_OBJECT
    Q_PROPERTY(bool connected READ connected NOTIFY connectedChanged)
public:
    explicit WS(QObject *parent = 0);

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

#endif // WS_H
