#ifndef WEBSOCKETPACKAGE_H
#define WEBSOCKETPACKAGE_H

#include <QByteArray>
#include <crypto++/aes.h>

class WebSocketPackage {
public:
    static const QByteArray MakeRandomBytes(const int &size);
//    byte[1] ToStdBytes();
};

#endif // WEBSOCKETPACKAGE_H
