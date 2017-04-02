#include "websocketconnection.h"

// Sends a new_chef request. Returns the msg ID of sent message.
int WebSocketConnection::outgoing_newChef(QString email, QString password) {
    QJsonObject data;
    data.insert("email", email);
    data.insert("password", password);
    auto outMsg = sendRequestMessage("new_chef", data);
    return outMsg->value(MSG::Meta).toObject().value(MSG::ID).toInt();
}
