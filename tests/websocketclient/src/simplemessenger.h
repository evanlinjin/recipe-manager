#ifndef SIMPLEMESSENGER_H
#define SIMPLEMESSENGER_H

#include <QObject>

#include "../../../qtclient/src/websocketconnection.h"

class SimpleMessenger : public QObject
{
    Q_OBJECT
public:
    explicit SimpleMessenger(WebSocketConnection* wsc, QObject *parent = 0);

private:
    WebSocketConnection* m_wsc;

signals:

public slots:
    void sendMsg(QString msg);

};

#endif // SIMPLEMESSENGER_H
