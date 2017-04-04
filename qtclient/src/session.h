#ifndef SESSION_H
#define SESSION_H

#include <QObject>
#include <QJsonArray>

#include "jsondb.h"

struct SessionInfo {
    QString sessionID;
    QString sessionKey;
    QString chefID;
    QString chefName;
    QString chefEmail;

    QJsonObject toJsonObj() {
        QJsonObject obj;
        obj.insert("session_id", sessionID);
        obj.insert("session_key", sessionKey);
        obj.insert("chef_id", chefID);
        obj.insert("chef_name", chefName);
        obj.insert("chef_email", chefEmail);
        return obj;
    }

    void loadFromJsonObj(const QJsonObject &obj) {
        sessionID = obj.value("session_id").toString();
        sessionKey = obj.value("session_key").toString();
        chefID = obj.value("chef_id").toString();
        chefName = obj.value("chef_name").toString();
        chefEmail = obj.value("chef_email").toString();
    }
};

class Session : public QObject
{
    Q_OBJECT
    Q_PROPERTY(QString url READ url WRITE setUrl NOTIFY urlChanged)
    Q_PROPERTY(QString sessionID READ sessionID NOTIFY sessionChanged)
    Q_PROPERTY(QString sessionKey READ sessionKey NOTIFY sessionChanged)
    Q_PROPERTY(QString chefID READ chefID NOTIFY sessionChanged)
    Q_PROPERTY(QString chefName READ chefName NOTIFY sessionChanged)
    Q_PROPERTY(QString chefEmail READ chefEmail NOTIFY sessionChanged)

public:
    explicit Session(QObject *parent = 0);

    QString url() const {return m_url;}
    QString sessionID() const {return m_sessionInfo.sessionID;}
    QString sessionKey() const {return m_sessionInfo.sessionKey;}
    QString chefID() const {return m_sessionInfo.chefID;}
    QString chefName() const {return m_sessionInfo.chefName;}
    QString chefEmail() const {return m_sessionInfo.chefEmail;}

    void setUrl(const QString &v) {
        if (m_url == v) return;
        m_url = v;
        QJsonObject obj;
        obj.insert("url", v);
        m_db.saveConfig("url", obj);
        emit urlChanged();
    }

private:
    JsonDB m_db;
    SessionInfo m_sessionInfo;
    QString m_url;

    bool loadUrl();
    bool loadSessionInfo();
    bool saveSessionInfo();

signals:
    void sessionChanged();
    void urlChanged();

public slots:
    void changeSession(int reqId, SessionInfo info);
    void clearSession();
};

#endif // SESSION_H
