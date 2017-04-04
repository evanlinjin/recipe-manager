#ifndef JSONDB_H
#define JSONDB_H

#include <QObject>
#include <QStandardPaths>
#include <QDir>
#include <QFile>
#include <QJsonDocument>
#include <QJsonObject>
#include <QDebug>

class JsonDB : public QObject
{
    Q_OBJECT
public:
    explicit JsonDB(QObject *parent = 0);

    const QString CONFIG_DIR = QStandardPaths::writableLocation(QStandardPaths::AppConfigLocation);
    const QString DATA_DIR = QStandardPaths::writableLocation(QStandardPaths::AppLocalDataLocation);

    bool saveConfig(const QString &fName, const QJsonObject &obj);
    bool getConfig(const QString &fname, QJsonObject* obj);

    bool saveFile(const QString &subDir, const QString &fName, const QJsonObject &obj);
    bool getFile(const QString &subDir, const QString &fName, QJsonObject *obj);

    qint64 getCount(const QString &subDir);
    QStringList getNameList(const QString &subDir);

private:

signals:

public slots:
};

#endif // JSONDB_H
