#include "websocketconnection.h"

WebSocketConnection::WebSocketConnection(QObject *parent) :
    QObject(parent), m_connectionStatus(1) {

    m_enc = new Encryptor(this);
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

    connect(&m_ws, SIGNAL(textMessageReceived(QString)),
            this, SLOT(onReceived(QString)));

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
    // Make random Signature.
    auto signature = m_enc->makeKey();

    // Make Data into Package with Signature.
    auto package = Package::MakePackage(obj, signature);

    // Encrypt Data and Signature.
    auto encSignature = m_enc->encrypt(signature);
    auto encPackage = m_enc->encrypt(package);

    // Join with dot.
    QByteArray out;
    out.append(encSignature).append('.').append(encPackage);

    m_ws.sendTextMessage(QString::fromLatin1(out));
}

void WebSocketConnection::onStateChanged(QAbstractSocket::SocketState state) {
    switch (state) {
    case QAbstractSocket::ConnectedState:
        m_timer->start(TIMER_INTERVAL_MS);

    case QAbstractSocket::UnconnectedState:
        m_enc->resetKey();
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

void WebSocketConnection::onReceived(QString data) {

    // Split msg into Signature and Data.
    QStringList split = data.split('.', QString::SkipEmptyParts);
    if (split.length() != 2) {
        qDebug("Expected %d parts, got %d", 2, split.length());
        return;
    }

    auto encSignature = split.at(0).toLatin1();
    auto encPackage = split.at(1).toLatin1();

    // Decrypt Signature and Data.
    auto signature = m_enc->decrypt(encSignature);
    auto package = m_enc->decrypt(encPackage);

    // Vertify Data with Signature.
    auto obj = Package::ReadPackage(package, signature);
    qInfo() << "[WebSocketConnection::onReceived]" << obj;
    if (m_msgs->checkIncomingMessage(obj) == false) {
        return;
    }

    // Convert to message struct and process.
    MSG::Message msg = MSG::obj_to_struct(obj);

    if (msg.cmd == "handshake")
        ps_handshake(msg);

    if (msg.cmd == "new_chef")
        ps_new_chef(msg);

    if (msg.cmd == "login")
        ps_login(msg);
}

// Processes incoming handshake request.
bool WebSocketConnection::ps_handshake(const MSG::Message &msg) {
    if (msg.typ != TYPE_REQUEST) {
        m_ws.close(QWebSocketProtocol::CloseCodeWrongDatatype,
                   "handshake is not request type");
        return false;
    }
    if (msg.data.isString() == false) {
        m_ws.close(QWebSocketProtocol::CloseCodeWrongDatatype,
                   "data does not contain string key");
        return false;
    }
    auto obj = sendResponseMessage(MSG::struct_to_obj(msg), true);
    if (obj == nullptr) {
        qDebug() << "[WebSocketConnection::ps_handshake]"
                 << "Failed to create response; got nullptr.";
        return false;
    }
    m_enc->setKey(msg.data.toString().toLatin1());
    return true;
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
    if (msg.data.isObject() == false)
        qDebug() << "[WebSocketConnection::ps_login]"
                 << "Wrong data type from server.";

    auto obj = msg.data.toObject();
    SessionInfo info;

    if (obj.value("okay") == true) {
        obj = obj.value("session").toObject();
        info.sessionID = obj.value("session_id").toString();
        info.sessionKey = obj.value("session_key").toString();
        info.chefID = obj.value("chef_id").toString();
        info.chefName = obj.value("chef_name").toString();
        info.chefEmail = obj.value("chef_email").toString();
        QJsonArray teams = obj.value("teams").toArray();
        for (int i = 0; i < teams.size(); i++) {
            info.teams.append(teams.at(i).toString());
        }
    } else {
        emit responseTextMessage(msg.req->id, obj.value("message").toString());
    }

    emit changeSession(msg.req->id, info);
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
