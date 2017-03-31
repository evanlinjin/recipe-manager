#include "messagemanager.h"

MessageManager::MessageManager(QObject *parent) :
    QObject(parent), outgoingId(0), incomingId(0)
{

}

QJsonObject* MessageManager::makeRequestMessage(const QString &cmd, const QJsonValue &data) {
    outgoingId += 1;

    QJsonObject meta;
    meta.insert(MSG::ID, outgoingId);
    meta.insert(MSG::Timestamp, QDateTime::currentSecsSinceEpoch());

    QJsonObject* msg = new QJsonObject();
    msg->insert(MSG::Command, cmd);
    msg->insert(MSG::Type, TYPE_REQUEST);
    msg->insert(MSG::Data, data);
    msg->insert(MSG::Meta, meta);
    return msg;
}

QJsonObject* MessageManager::makeResponseMessage(const QJsonObject &reqMsg, const QJsonValue &data) {
    if (reqMsg.contains(MSG::Meta) == false) {
        qInfo() << "[MessageManager::makeResponseMessage] Error:"
                << "reqMsg.Meta is nil";
        return nullptr;
    }
    auto reqId = reqMsg.value(MSG::Meta).toObject().value(MSG::ID).toInt(0);
    if (reqId >= 0) {
        qInfo() << "[MessageManager::makeResponseMessage] Error:"
                << "reqMsg.Meta.ID is > 0, at" << reqId;
        return nullptr;
    }
    outgoingId += 1;

    QJsonObject meta;
    meta.insert(MSG::ID, outgoingId);
    meta.insert(MSG::Timestamp, QDateTime::currentSecsSinceEpoch());

    QJsonObject* msg = new QJsonObject();
    msg->insert(MSG::Command, reqMsg.value(MSG::Command));
    msg->insert(MSG::Type, TYPE_RESPONSE);
    msg->insert(MSG::Meta, meta);
    msg->insert(MSG::ReqMeta, reqMsg.value(MSG::Meta));

    if (!data.isNull() && !data.isUndefined()) {
        msg->insert(MSG::Data, data);
    }
    return msg;
}

bool MessageManager::checkIncomingMessage(QJsonObject &msg) {
    if (msg.value(MSG::Meta).isUndefined()) {
        qInfo() << "[MessageManager::checkIncomingMessage] Error:"
                << "msg has no meta";
        return false;
    }
    const qint64 id = msg.value(MSG::Meta).toObject().value(MSG::ID).toInt(0);
    if (id >= 0) {
        qInfo() << "[MessageManager::checkIncomingMessage] Error:"
                << "incoming msg has +ve id;" << id;
        return false;
    }
    if (msg.value(MSG::Type).toInt(-1) == TYPE_RESPONSE) {
        if (msg.contains(MSG::ReqMeta) == false) {
            qInfo() << "[MessageManager::checkIncomingMessage]"
                    << "response msg has no meta";
            return false;
        }
        auto reqId = msg.value(MSG::ReqMeta).toObject().value(MSG::ID).toInt(0);
        if (reqId <= 0) {
            qInfo() << "[MessageManager::checkIncomingMessage]"
                    << "reqMsg of incoming msg is invalid with id;" << reqId;
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

MSG::Message MSG::obj_to_struct(const QJsonObject &obj) {
    MSG::Message msg;
    msg.meta = msg.req = nullptr;

    msg.cmd = obj.value(MSG::Command).toString();
    msg.typ = obj.value(MSG::Type).toInt(0);
    msg.data = obj.value(MSG::Data);

    auto objMeta = obj.value(MSG::Meta).toObject();
    msg.meta = new MSG::MessageMeta;
    msg.meta->id = objMeta.value(MSG::ID).toInt();
    msg.meta->ts = objMeta.value(MSG::Timestamp).toInt();

    if (msg.typ == TYPE_RESPONSE) {
        auto objReq = obj.value(MSG::ReqMeta).toObject();
        msg.req = new MSG::MessageMeta;
        msg.req->id = objReq.value(MSG::ID).toInt();
        msg.req->ts = objReq.value(MSG::Timestamp).toInt();
    }
    return msg;
}

QJsonObject MSG::struct_to_obj(const MSG::Message &msg) {
    QJsonObject obj;
    obj.insert(MSG::Command, msg.cmd);
    obj.insert(MSG::Type, msg.typ);
    obj.insert(MSG::Data, msg.data);

    QJsonObject meta;
    meta.insert(MSG::ID, msg.meta->id);
    meta.insert(MSG::Timestamp, msg.meta->ts);
    obj.insert(MSG::Meta, meta);

    if (msg.typ == TYPE_RESPONSE) {
        QJsonObject req;
        req.insert(MSG::ID, msg.req->id);
        req.insert(MSG::Timestamp, msg.req->ts);
        obj.insert(MSG::ReqMeta, req);
    }
    return obj;
}
