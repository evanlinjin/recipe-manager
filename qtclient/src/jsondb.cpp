#include "jsondb.h"

JsonDB::JsonDB(QObject *parent) : QObject(parent) {
    QDir configDir(CONFIG_DIR);
    if (!configDir.exists()) configDir.mkpath(".");
    QDir dataDir(DATA_DIR);
    if (!dataDir.exists()) dataDir.mkpath(".");
    qInfo() << "CONFIG_DIR" << CONFIG_DIR;
    qInfo() << "DATA_CIR" << DATA_DIR;
}

bool JsonDB::saveConfig(const QString &fName, const QJsonObject &obj) {
    QDir dir(CONFIG_DIR);

    QFile file(dir.absoluteFilePath(fName));
    if (!file.open(QIODevice::WriteOnly)) {
        qDebug() << "[JsonDB::saveConfig] Error:"
                 << "Unable to open file.";
        return false;
    }

    if (file.write(QJsonDocument(obj).toBinaryData()) == -1) {
        qDebug() << "[JsonDB::saveConfig] Error:"
                 << "Unable to write to file.";
        return false;
    }

    qDebug() << "[JsonDB::saveConfig] Saved:" << file.fileName();
    return true;
}

bool JsonDB::getConfig(const QString &fName, QJsonObject *obj) {
    QDir dir(CONFIG_DIR);

    QFile file(dir.absoluteFilePath(fName));
    if (!file.open(QIODevice::ReadOnly)) {
        qDebug() << "[JsonDB::getConfig] Error:"
                 << "Unable to open file.";
        return false;
    }

    *obj = QJsonDocument::fromBinaryData(file.readAll()).object();

    qDebug() << "[JsonDB::getConfig] Got:" << file.fileName();
    return true;
}

bool JsonDB::saveFile(const QString &subDir, const QString &fName, const QJsonObject &obj) {
    QDir dir(DATA_DIR);

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
    if (file.write(QJsonDocument(obj).toBinaryData()) == -1) {
        qDebug() << "[JsonDB::saveFile] Error:"
                 << "Unable to write to file.";
        return false;
    }

    qDebug() << "[JsonDB::saveFile] Saved:" << file.fileName();
    return true;
}

bool JsonDB::getFile(const QString &subDir, const QString &fName, QJsonObject *obj) {
    QDir dir(DATA_DIR);

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
    *obj = QJsonDocument::fromBinaryData(file.readAll()).object();

    qDebug() << "[JsonDB::getFile] Got:" << file.fileName();
    return true;
}

qint64 JsonDB::getCount(const QString &subDir) {
    QDir dir(DATA_DIR);
    if (dir.cd(subDir) == false)
        return -1;

    return (qint64)dir.entryList(QDir::Files).count();
}

QStringList JsonDB::getNameList(const QString &subDir) {
    QDir dir(DATA_DIR);
    if (dir.cd(subDir) == false)
        return QStringList();

    return dir.entryList(QDir::Files);
}
