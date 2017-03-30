#ifndef MESSAGEMANAGER_H
#define MESSAGEMANAGER_H

#include <QObject>
#include <QJsonValue>
#include <QJsonObject>
#include <QDateTime>
#include <QDebug>

namespace Message
{
    const char Command[] = "cmd";
    const char Type[] = "type";
    const char Data[] = "data";
    const char Meta[] = "meta";
    const char ReqMeta[] = "req";

    const char ID[] = "id";
    const char Timestamp[] = "ts";
}

class MessageManager : public QObject
{
    Q_OBJECT
public:
    explicit MessageManager(QObject *parent = 0);

    const int TYPE_REQUEST = 0;
    const int TYPE_RESPONSE = 1;

    QJsonObject* makeRequestMessage(QString &cmd, QJsonValue &data);
    QJsonObject* makeResponseMessage(QJsonObject &reqMsg, QJsonValue &data);
    bool checkIncomingMessage(QJsonObject &msg);

private:
    int outgoingId;
    int incomingId;

signals:

public slots:
};

#endif // MESSAGEMANAGER_H
