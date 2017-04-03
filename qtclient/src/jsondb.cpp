#include "jsondb.h"

JsonDB::JsonDB(QObject *parent) : QObject(parent) {

}

bool JsonDB::setSession(const QString &name) {
    m_session = name;
    if (QDir(CONFIG_DIR).mkpath(m_session) == false)
        return false;
    if (QDir(DATA_DIR).mkpath(m_session) == false)
        return false;
    return true;
}

bool JsonDB::saveGlobalConfig(const QJsonObject &obj) {
    QDir dir(CONFIG_DIR);
    QFile file(dir.absoluteFilePath("config"));
    if (!file.open(QIODevice::WriteOnly)) {
        qDebug() << "[JsonDB::saveGlobalConfig] Error:"
                 << "Unable to open file.";
        return false;
    }

    auto data = m_enc.encrypt(QJsonDocument(obj).toBinaryData());
    if (file.write(data) == -1) {
        qDebug() << "[JsonDB::saveGlobalConfig] Error:"
                 << "Unable to write to file.";
        return false;
    }

    qDebug() << "[JsonDB::saveGlobalConfig] Saved:" << file.fileName();
    return true;
}

bool JsonDB::getGlobalConfig(QJsonObject* obj) {
    QDir dir(CONFIG_DIR);
    QFile file(dir.absoluteFilePath("config"));
    if (!file.open(QIODevice::ReadOnly)) {
        qDebug() << "[JsonDB::getConfig] Error:"
                 << "Unable to open file.";
        return false;
    }

    auto data = m_enc.decrypt(file.readAll());
    *obj = QJsonDocument::fromBinaryData(data).object();

    qDebug() << "[JsonDB::getConfig] Got:" << file.fileName();
    return true;
}

bool JsonDB::saveConfig(const QJsonObject &obj) {
    QDir dir(CONFIG_DIR);
    if (dir.cd(m_session) == false) {
        qDebug() << "[JsonDB::saveConfig] Error:"
                 << "Unable to open session directory.";
        return false;
    }

    QFile file(dir.absoluteFilePath("config"));
    if (!file.open(QIODevice::WriteOnly)) {
        qDebug() << "[JsonDB::saveConfig] Error:"
                 << "Unable to open file.";
        return false;
    }

    auto data = m_enc.encrypt(QJsonDocument(obj).toBinaryData());
    if (file.write(data) == -1) {
        qDebug() << "[JsonDB::saveConfig] Error:"
                 << "Unable to write to file.";
        return false;
    }

    qDebug() << "[JsonDB::saveConfig] Saved:" << file.fileName();
    return true;
}

bool JsonDB::getConfig(QJsonObject *obj) {
    QDir dir(CONFIG_DIR);
    if (dir.cd(m_session) == false) {
        qDebug() << "[JsonDB::getConfig] Error:"
                 << "Unable to open session directory.";
        return false;
    }

    QFile file(dir.absoluteFilePath("config"));
    if (!file.open(QIODevice::ReadOnly)) {
        qDebug() << "[JsonDB::getConfig] Error:"
                 << "Unable to open file.";
        return false;
    }

    auto data = m_enc.decrypt(file.readAll());
    *obj = QJsonDocument::fromBinaryData(data).object();

    qDebug() << "[JsonDB::getConfig] Got:" << file.fileName();
    return true;
}

bool JsonDB::saveFile(const QString &subDir, const QString &fName, const QJsonObject &obj) {
    QDir dir(DATA_DIR);

    // Open session dir.
    if (dir.cd(m_session) == false) {
        qDebug() << "[JsonDB::saveFile] Error:"
                 << "Unable to open session dir.";
        return false;
    }

    // Ensure subDir exists.
    if (dir.mkpath(subDir) == false) {
        qDebug() << "[JsonDB::saveFile] Error:"
                 << "Unable to create subDir" << subDir;
        return false;
    }

    // Open subDir.
    if (dir.cd(subDir) == false) {
        qDebug() << "[JsonDB::saveFile] Error:"
                 << "Unable to open subDir" << subDir;
        return false;
    }

    // Open file for writing.
    QFile file(dir.absoluteFilePath(fName));
    if (!file.open(QIODevice::WriteOnly)) {
        qDebug() << "[JsonDB::saveFile] Error:"
                 << "Unable to open file for writing.";
        return false;
    }

    // Write to file.
    auto data = m_enc.encrypt(QJsonDocument(obj).toBinaryData());
    if (file.write(data) == -1) {
        qDebug() << "[JsonDB::saveFile] Error:"
                 << "Unable to write to file.";
        return false;
    }

    qDebug() << "[JsonDB::saveFile] Saved:" << file.fileName();
    return true;
}

bool JsonDB::getFile(const QString &subDir, const QString &fName, QJsonObject *obj) {
    QDir dir(DATA_DIR);

    // Open session dir.
    if (dir.cd(m_session) == false) {
        qDebug() << "[JsonDB::getFile] Error:"
                 << "Unable to open session dir.";
        return false;
    }

    // Ensure subDir exists.
    if (dir.mkpath(subDir) == false) {
        qDebug() << "[JsonDB::getFile] Error:"
                 << "Unable to create subDir" << subDir;
        return false;
    }

    // Open subDir.
    if (dir.cd(subDir) == false) {
        qDebug() << "[JsonDB::getFile] Error:"
                 << "Unable to open subDir" << subDir;
        return false;
    }

    // Open file for reading.
    QFile file(dir.absoluteFilePath(fName));
    if (!file.open(QIODevice::ReadOnly)) {
        qDebug() << "[JsonDB::getFile] Error:"
                 << "Unable to open file for reading.";
        return false;
    }

    // Read file.
    auto data = m_enc.decrypt(file.readAll());
    *obj = QJsonDocument::fromBinaryData(data).object();

    qDebug() << "[JsonDB::getFile] Got:" << file.fileName();
    return true;
}

qint64 JsonDB::getCount(const QString &subDir) {
    QDir dir(DATA_DIR);
    if (dir.cd(m_session) == false)
        return -1;
    if (dir.cd(subDir) == false)
        return -1;

    return (qint64)dir.entryList(QDir::Files).count();
}

QStringList JsonDB::getNameList(const QString &subDir) {
    QDir dir(DATA_DIR);
    if (dir.cd(m_session) == false)
        return QStringList();
    if (dir.cd(subDir) == false)
        return QStringList();

    return dir.entryList(QDir::Files);
}
