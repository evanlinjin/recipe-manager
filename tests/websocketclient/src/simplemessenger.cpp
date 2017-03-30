#include "simplemessenger.h"

SimpleMessenger::SimpleMessenger(WebSocketConnection *wsc, QObject *parent) :
    QObject(parent), m_wsc(wsc) {}

void SimpleMessenger::sendMsg(QString msg) {
    m_wsc->sendRequestMessage("undefined", msg);
}
