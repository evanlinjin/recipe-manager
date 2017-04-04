#include "websocketconnection.h"

WebSocketConnection::WebSocketConnection(QObject *parent) :
    QObject(parent), m_connectionStatus(1) {

    m_msgs = new MessageManager(this);

    m_timer = new QTimer(this);
    connect(m_timer, SIGNAL(timeout()),
            &m_ws, SLOT(ping()));

    m_checker = new QTimer(this);
    m_checker->setSingleShot(true);
    connect(&m_ws, SIGNAL(pong(quint64,QByteArray)),
            m_checker, SLOT(stop()));
    connect(m_checker, SIGNAL(timeout()),
            this, SIGNAL(networkError()));

    connect(&m_ws, SIGNAL(stateChanged(QAbstractSocket::SocketState)),
            this, SLOT(onStateChanged(QAbstractSocket::SocketState)));

    connect(&m_ws, SIGNAL(binaryMessageReceived(QByteArray)),
            this, SLOT(onReceived(QByteArray)));

    connect(&m_ws, SIGNAL(error(QAbstractSocket::SocketError)),
            this, SLOT(onError(QAbstractSocket::SocketError)));

    typedef void (QWebSocket:: *sslErrorsSignal)(const QList<QSslError> &);
    connect(&m_ws, static_cast<sslErrorsSignal>(&QWebSocket::sslErrors),
            this, &WebSocketConnection::onSslErrors);
}

QJsonObject *WebSocketConnection::sendRequestMessage(const QString &cmd, const QJsonValue &data) {
    auto msg = m_msgs->makeRequestMessage(cmd, data);
    if (msg == nullptr) {
        return msg;
    }
    send(*msg);
    return msg;
}

QJsonObject *WebSocketConnection::sendResponseMessage(const QJsonObject &reqMsg, const QJsonValue &data) {
    auto msg = m_msgs->makeResponseMessage(reqMsg, data);
    if (msg == nullptr) {
        return msg;
    }
    send(*msg);
    return msg;
}

void WebSocketConnection::send(QJsonObject &obj) {
    qInfo() << "[WebSocketConnection::send]" << obj;
    auto pkg = QJsonDocument(obj).toJson(QJsonDocument::Compact);
    m_ws.sendBinaryMessage(pkg);
}

void WebSocketConnection::onStateChanged(QAbstractSocket::SocketState state) {
    switch (state) {
    case QAbstractSocket::ConnectedState:
        m_timer->start(TIMER_INTERVAL_MS);
        // START : Claim session ---------------------------------------------->
        if (m_session->sessionID() != "")
            outgoing_claimSession();
        // END : Claim session ------------------------------------------------>

    case QAbstractSocket::UnconnectedState:
        m_msgs->resetIds();
        m_timer->stop();

    default:
    {}
    }
    m_connectionStatus = state;
    emit connectionStatusChanged();
    qDebug() << "[WebSocketConnection::onStateChanged] State:" << state;
}

/******************************************************************************/
/* INCOMING MESSAGES FROM SERVER TO CLIENT                                    */
/******************************************************************************/

void WebSocketConnection::onReceived(QByteArray data) {
    auto obj = QJsonDocument::fromJson(data).object();
    qDebug() << "[WebSocketConnection::onReceived]" << obj;

    // Convert to message struct and process.
    MSG::Message msg = MSG::obj_to_struct(obj);

    if (msg.cmd == "new_chef")
        ps_new_chef(msg);

    if (msg.cmd == "login")
        ps_login(msg);

    if (msg.cmd == "claim_session")
        ps_claimSession(msg);
}

// Processes incoming new_chef response.
bool WebSocketConnection::ps_new_chef(const MSG::Message &msg) {
    if (msg.data.isString() == false)
        qDebug() << "[WebSocketConnection::ps_new_chef]"
                 << "Invalid response from server.";

    emit responseTextMessage(msg.req->id, msg.data.toString());
    return true;
}

// Processes incoming login response.
bool WebSocketConnection::ps_login(const MSG::Message &msg) {
    // On unauthorised.
    if (msg.data.isString()) {
        emit responseTextMessage(msg.req->id, msg.data.toString());
        return true;
    }

    if (msg.data.isObject() == false)
        qDebug() << "[WebSocketConnection::ps_login]"
                 << "Wrong data type from server.";

    auto obj = msg.data.toObject();
    SessionInfo info;
    info.loadFromJsonObj(obj);

    m_session->changeSession(msg.req->id, info);
    return true;
}

// Processes incoming claim_session response.
bool WebSocketConnection::ps_claimSession(const MSG::Message &msg) {
    if (msg.data.isObject() == false)
        qDebug() << "[WebSocketConnection::ps_claimSession]"
                 << "Wrong data type from server.";

    auto obj = msg.data.toObject();

    // Remove session if not accepted.
    if (obj.value("okay") == false) {
        emit responseTextMessage(msg.req->id, obj.value("message").toString());
        m_session->clearSession();
        return false;
    }

    // Reload session info.
    SessionInfo info;
    info.loadFromJsonObj(obj.value("session").toObject());
    m_session->changeSession(msg.req->id, info);
    return true;
}

/******************************************************************************/
/* OUTGOING MESSAGES FROM CLIENT TO SERVER                                    */
/******************************************************************************/

// Sends a new_chef request. Returns the msg ID of sent message.
int WebSocketConnection::outgoing_newChef(QString email, QString password) {
    QJsonObject data;
    data.insert("email", email);
    data.insert("password", password);
    auto outMsg = sendRequestMessage("new_chef", data);
    return outMsg->value(MSG::Meta).toObject().value(MSG::ID).toInt();
}

// Sends a login request. Returns the msg ID of sent message.
int WebSocketConnection::outgoing_login(QString email, QString password) {
    QJsonObject data;
    data.insert("email", email);
    data.insert("password", password);
    auto outMsg = sendRequestMessage("login", data);
    return outMsg->value(MSG::Meta).toObject().value(MSG::ID).toInt();
}

// Sends a logout request. Returns the msg ID of sent message.
int WebSocketConnection::outgoing_logout() {
    QJsonObject data;
    auto outMsg = sendRequestMessage("logout", data);
    m_session->clearSession();
    return outMsg->value(MSG::Meta).toObject().value(MSG::ID).toInt();
}

// Sends a request to server to claim a session with session_id and session_key.
int WebSocketConnection::outgoing_claimSession() {
    QJsonObject obj;
    obj.insert("chef_id", m_session->chefID());
    obj.insert("session_id", m_session->sessionID());
    obj.insert("session_key", m_session->sessionKey());
    auto outMsg = sendRequestMessage("claim_session", obj);
    return outMsg->value(MSG::Meta).toObject().value(MSG::ID).toInt();
}
