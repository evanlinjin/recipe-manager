#include "session.h"

Session::Session(QObject *parent) : QObject(parent) {
    QJsonObject config;

    if (m_db.getGlobalConfig(&config) == true) {
        QJsonObject session = config.value("session").toObject();
        getSessionFromObject(session);
    }
}

bool Session::getSessionFromObject(const QJsonObject &obj) {
    if (!obj.isEmpty()) {
        m_sessionInfo.sessionID =
                obj.value("session_id").toString();
        m_sessionInfo.sessionKey =
                obj.value("session_key").toString();
        m_sessionInfo.chefID =
                obj.value("chef_id").toString();
        m_sessionInfo.chefName =
                obj.value("chef_name").toString();
        m_sessionInfo.chefEmail =
                obj.value("chef_email").toString();

        auto teamsArray = obj.value("teams").toArray();
        for (int i = 0; i < teamsArray.size(); i++) {
            m_sessionInfo.teams.append(teamsArray.at(i).toString());
        }

        emit sessionChanged();
        return true;
    }
    return false;
}

void Session::changeSession(int reqId, SessionInfo info) {
    m_sessionInfo = info;
    emit sessionChanged();
}
