#include "session.h"

Session::Session(QObject *parent) : QObject(parent) {
    loadUrl();
    loadSessionInfo();
}

bool Session::loadUrl() {
    QJsonObject obj;
    if (m_db.getConfig("url", &obj) == false) {
        return false;
    }
    if (!obj.isEmpty()) {
        m_url = obj.value("url").toString();
        emit urlChanged();
        return true;
    }
    return false;
}

bool Session::loadSessionInfo() {
    QJsonObject obj;
    if (m_db.getConfig("session", &obj) == false) {
        m_sessionInfo = SessionInfo();
        emit sessionChanged();
        return false;
    }
    if (!obj.isEmpty()) {
        m_sessionInfo.loadFromJsonObj(obj);
        emit sessionChanged();
        return true;
    }
    return false;
}

bool Session::saveSessionInfo() {
    return m_db.saveConfig("session", m_sessionInfo.toJsonObj());
}

void Session::changeSession(int reqId, SessionInfo info) {
    Q_UNUSED(reqId)
    m_sessionInfo = info;
    saveSessionInfo();
    emit sessionChanged();
}

void Session::clearSession() {
    QJsonObject obj;
    m_db.saveConfig("session", obj);
    m_sessionInfo = SessionInfo();
    emit sessionChanged();
}
