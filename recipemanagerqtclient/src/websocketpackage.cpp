#include "websocketpackage.h"

const QByteArray WebSocketPackage::MakeRandomBytes(const int &size) {
    QByteArray data(size, Qt::Uninitialized);
        for(int i = 0; i < size; i++)
            data[i] = (char)qrand();
//            data[i] = (char)0x00;
        return data;
}
