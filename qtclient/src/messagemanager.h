#ifndef MESSAGEMANAGER_H
#define MESSAGEMANAGER_H

#include <array>

#include <QObject>
#include <QJsonValue>
#include <QJsonObject>
#include <QDateTime>
#include <QDebug>

#define TYPE_REQUEST 0
#define TYPE_RESPONSE 1

namespace MSG
{
const char Command[] = "cmd";
const char Type[] = "typ";
const char Data[] = "data";
const char Meta[] = "meta";
const char ReqMeta[] = "req";

const char ID[] = "id";
const char Timestamp[] = "ts";

struct MessageMeta {
    qint64 id;
    qint64 ts;
};

struct Message {
    QString cmd;
    qint64 typ;
    QJsonValue data;
    MessageMeta* meta;
    MessageMeta* req;
};

Message obj_to_struct(const QJsonObject &obj);
QJsonObject struct_to_obj(const Message &msg);
}

class MessageManager : public QObject
{
    Q_OBJECT
public:
    explicit MessageManager(QObject *parent = 0);

    QJsonObject* makeRequestMessage(const QString &cmd, const QJsonValue &data);
    QJsonObject* makeResponseMessage(const QJsonObject &reqMsg, const QJsonValue &data);
    bool checkIncomingMessage(QJsonObject &msg);

    void resetIds() {outgoingId = incomingId = 0;}

private:
    qint64 outgoingId;
    qint64 incomingId;

signals:

public slots:
};

#endif // MESSAGEMANAGER_H
