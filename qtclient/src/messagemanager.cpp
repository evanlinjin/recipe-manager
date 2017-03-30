#include "messagemanager.h"

MessageManager::MessageManager(QObject *parent) :
    QObject(parent), outgoingId(0), incomingId(0)
{

}

QJsonObject* MessageManager::makeRequestMessage(QString &cmd, QJsonValue &data) {
    outgoingId += 1;

    QJsonObject meta;
    meta.insert(Message::ID, outgoingId);
    meta.insert(Message::Timestamp, QDateTime::currentSecsSinceEpoch());

    QJsonObject* msg = new QJsonObject();
    msg->insert(Message::Command, cmd);
    msg->insert(Message::Type, TYPE_REQUEST);
    msg->insert(Message::Data, data);
    msg->insert(Message::Meta, meta);
    return msg;
}

QJsonObject* MessageManager::makeResponseMessage(QJsonObject &reqMsg, QJsonValue &data) {
    if (reqMsg.contains(Message::Meta) == false) {
        qInfo() << "reqMsg.Meta is nil";
        return nullptr;
    }
    auto reqId = reqMsg.value(Message::Meta).toObject().value(Message::ID).toInt(0);
    if (reqId >= 0) {
        qInfo() << "reqMsg.Meta.ID is > 0, at" << reqId;
        return nullptr;
    }
    outgoingId += 1;

    QJsonObject meta;
    meta.insert(Message::ID, outgoingId);
    meta.insert(Message::Timestamp, QDateTime::currentSecsSinceEpoch());

    QJsonObject* msg = new QJsonObject();
    msg->insert(Message::Command, reqMsg.value(Message::Command));
    msg->insert(Message::Type, TYPE_RESPONSE);
    msg->insert(Message::Meta, meta);
    msg->insert(Message::ReqMeta, reqMsg.value(Message::Meta));

    if (!data.isNull() && !data.isUndefined()) {
        msg->insert(Message::Data, data);
    }
    return msg;
}

bool MessageManager::checkIncomingMessage(QJsonObject &msg) {
    if (msg.value(Message::Meta).isUndefined()) {
        qInfo() << "msg has no meta";
        return false;
    }
    const int id = msg.value(Message::Meta).toObject().value(Message::ID).toInt(0);
    if (id >= 0) {
        qInfo() << "incoming msg has +ve id;" << id;
        return false;
    }
    if (msg.value(Message::Type).toInt(-1) == TYPE_RESPONSE) {
        if (msg.contains(Message::ReqMeta) == false) {
            qInfo() << "response msg has no meta";
            return false;
        }
        auto reqId = msg.value(Message::ReqMeta).toObject().value(Message::ID).toInt(0);
        if (reqId <= 0) {
            qInfo() << "reqMsg of incoming msg is invalid with id;" << reqId;
            return false;
        }
    }
    if (id != incomingId-1) {
        qInfo() << "[MessageManager::checkIncomingMessage] Unexpected ID."
                << "GOT:" << id
                << "EXPECTED:" << incomingId-1;
    }
    incomingId = id;
    return true;
}
