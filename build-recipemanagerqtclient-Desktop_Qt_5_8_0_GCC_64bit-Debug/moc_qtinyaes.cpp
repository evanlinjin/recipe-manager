/****************************************************************************
** Meta object code from reading C++ file 'qtinyaes.h'
**
** Created by: The Qt Meta Object Compiler version 67 (Qt 5.8.0)
**
** WARNING! All changes made in this file will be lost!
*****************************************************************************/

#include "../recipemanagerqtclient/src/QTinyAes/QTinyAes/qtinyaes.h"
#include <QtCore/qbytearray.h>
#include <QtCore/qmetatype.h>
#if !defined(Q_MOC_OUTPUT_REVISION)
#error "The header file 'qtinyaes.h' doesn't include <QObject>."
#elif Q_MOC_OUTPUT_REVISION != 67
#error "This file was generated using the moc from 5.8.0. It"
#error "cannot be used with the include files from this version of Qt."
#error "(The moc has changed too much.)"
#endif

QT_BEGIN_MOC_NAMESPACE
QT_WARNING_PUSH
QT_WARNING_DISABLE_DEPRECATED
struct qt_meta_stringdata_QTinyAes_t {
    QByteArrayData data[17];
    char stringdata0[108];
};
#define QT_MOC_LITERAL(idx, ofs, len) \
    Q_STATIC_BYTE_ARRAY_DATA_HEADER_INITIALIZER_WITH_OFFSET(len, \
    qptrdiff(offsetof(qt_meta_stringdata_QTinyAes_t, stringdata0) + ofs \
        - idx * sizeof(QByteArrayData)) \
    )
static const qt_meta_stringdata_QTinyAes_t qt_meta_stringdata_QTinyAes = {
    {
QT_MOC_LITERAL(0, 0, 8), // "QTinyAes"
QT_MOC_LITERAL(1, 9, 7), // "setMode"
QT_MOC_LITERAL(2, 17, 0), // ""
QT_MOC_LITERAL(3, 18, 10), // "CipherMode"
QT_MOC_LITERAL(4, 29, 4), // "mode"
QT_MOC_LITERAL(5, 34, 6), // "setKey"
QT_MOC_LITERAL(6, 41, 3), // "key"
QT_MOC_LITERAL(7, 45, 8), // "resetKey"
QT_MOC_LITERAL(8, 54, 5), // "setIv"
QT_MOC_LITERAL(9, 60, 2), // "iv"
QT_MOC_LITERAL(10, 63, 7), // "resetIv"
QT_MOC_LITERAL(11, 71, 7), // "encrypt"
QT_MOC_LITERAL(12, 79, 5), // "plain"
QT_MOC_LITERAL(13, 85, 7), // "decrypt"
QT_MOC_LITERAL(14, 93, 6), // "cipher"
QT_MOC_LITERAL(15, 100, 3), // "CBC"
QT_MOC_LITERAL(16, 104, 3) // "ECB"

    },
    "QTinyAes\0setMode\0\0CipherMode\0mode\0"
    "setKey\0key\0resetKey\0setIv\0iv\0resetIv\0"
    "encrypt\0plain\0decrypt\0cipher\0CBC\0ECB"
};
#undef QT_MOC_LITERAL

static const uint qt_meta_data_QTinyAes[] = {

 // content:
       7,       // revision
       0,       // classname
       0,    0, // classinfo
       7,   14, // methods
       3,   66, // properties
       1,   75, // enums/sets
       0,    0, // constructors
       0,       // flags
       0,       // signalCount

 // slots: name, argc, parameters, tag, flags
       1,    1,   49,    2, 0x0a /* Public */,
       5,    1,   52,    2, 0x0a /* Public */,
       7,    0,   55,    2, 0x0a /* Public */,
       8,    1,   56,    2, 0x0a /* Public */,
      10,    0,   59,    2, 0x0a /* Public */,

 // methods: name, argc, parameters, tag, flags
      11,    1,   60,    2, 0x02 /* Public */,
      13,    1,   63,    2, 0x02 /* Public */,

 // slots: parameters
    QMetaType::Void, 0x80000000 | 3,    4,
    QMetaType::Void, QMetaType::QByteArray,    6,
    QMetaType::Void,
    QMetaType::Void, QMetaType::QByteArray,    9,
    QMetaType::Void,

 // methods: parameters
    QMetaType::QByteArray, QMetaType::QByteArray,   12,
    QMetaType::QByteArray, QMetaType::QByteArray,   14,

 // properties: name, type, flags
       4, 0x80000000 | 3, 0x0009510b,
       6, QMetaType::QByteArray, 0x00095107,
       9, QMetaType::QByteArray, 0x00095107,

 // enums: name, flags, count, data
       3, 0x0,    2,   79,

 // enum data: key, value
      15, uint(QTinyAes::CBC),
      16, uint(QTinyAes::ECB),

       0        // eod
};

void QTinyAes::qt_static_metacall(QObject *_o, QMetaObject::Call _c, int _id, void **_a)
{
    if (_c == QMetaObject::InvokeMetaMethod) {
        QTinyAes *_t = static_cast<QTinyAes *>(_o);
        Q_UNUSED(_t)
        switch (_id) {
        case 0: _t->setMode((*reinterpret_cast< CipherMode(*)>(_a[1]))); break;
        case 1: _t->setKey((*reinterpret_cast< QByteArray(*)>(_a[1]))); break;
        case 2: _t->resetKey(); break;
        case 3: _t->setIv((*reinterpret_cast< QByteArray(*)>(_a[1]))); break;
        case 4: _t->resetIv(); break;
        case 5: { QByteArray _r = _t->encrypt((*reinterpret_cast< QByteArray(*)>(_a[1])));
            if (_a[0]) *reinterpret_cast< QByteArray*>(_a[0]) = _r; }  break;
        case 6: { QByteArray _r = _t->decrypt((*reinterpret_cast< QByteArray(*)>(_a[1])));
            if (_a[0]) *reinterpret_cast< QByteArray*>(_a[0]) = _r; }  break;
        default: ;
        }
    }
#ifndef QT_NO_PROPERTIES
    else if (_c == QMetaObject::ReadProperty) {
        QTinyAes *_t = static_cast<QTinyAes *>(_o);
        Q_UNUSED(_t)
        void *_v = _a[0];
        switch (_id) {
        case 0: *reinterpret_cast< CipherMode*>(_v) = _t->mode(); break;
        case 1: *reinterpret_cast< QByteArray*>(_v) = _t->key(); break;
        case 2: *reinterpret_cast< QByteArray*>(_v) = _t->iv(); break;
        default: break;
        }
    } else if (_c == QMetaObject::WriteProperty) {
        QTinyAes *_t = static_cast<QTinyAes *>(_o);
        Q_UNUSED(_t)
        void *_v = _a[0];
        switch (_id) {
        case 0: _t->setMode(*reinterpret_cast< CipherMode*>(_v)); break;
        case 1: _t->setKey(*reinterpret_cast< QByteArray*>(_v)); break;
        case 2: _t->setIv(*reinterpret_cast< QByteArray*>(_v)); break;
        default: break;
        }
    } else if (_c == QMetaObject::ResetProperty) {
        QTinyAes *_t = static_cast<QTinyAes *>(_o);
        Q_UNUSED(_t)
        switch (_id) {
        case 1: _t->resetKey(); break;
        case 2: _t->resetIv(); break;
        default: break;
        }
    }
#endif // QT_NO_PROPERTIES
}

const QMetaObject QTinyAes::staticMetaObject = {
    { &QObject::staticMetaObject, qt_meta_stringdata_QTinyAes.data,
      qt_meta_data_QTinyAes,  qt_static_metacall, Q_NULLPTR, Q_NULLPTR}
};


const QMetaObject *QTinyAes::metaObject() const
{
    return QObject::d_ptr->metaObject ? QObject::d_ptr->dynamicMetaObject() : &staticMetaObject;
}

void *QTinyAes::qt_metacast(const char *_clname)
{
    if (!_clname) return Q_NULLPTR;
    if (!strcmp(_clname, qt_meta_stringdata_QTinyAes.stringdata0))
        return static_cast<void*>(const_cast< QTinyAes*>(this));
    return QObject::qt_metacast(_clname);
}

int QTinyAes::qt_metacall(QMetaObject::Call _c, int _id, void **_a)
{
    _id = QObject::qt_metacall(_c, _id, _a);
    if (_id < 0)
        return _id;
    if (_c == QMetaObject::InvokeMetaMethod) {
        if (_id < 7)
            qt_static_metacall(this, _c, _id, _a);
        _id -= 7;
    } else if (_c == QMetaObject::RegisterMethodArgumentMetaType) {
        if (_id < 7)
            *reinterpret_cast<int*>(_a[0]) = -1;
        _id -= 7;
    }
#ifndef QT_NO_PROPERTIES
   else if (_c == QMetaObject::ReadProperty || _c == QMetaObject::WriteProperty
            || _c == QMetaObject::ResetProperty || _c == QMetaObject::RegisterPropertyMetaType) {
        qt_static_metacall(this, _c, _id, _a);
        _id -= 3;
    } else if (_c == QMetaObject::QueryPropertyDesignable) {
        _id -= 3;
    } else if (_c == QMetaObject::QueryPropertyScriptable) {
        _id -= 3;
    } else if (_c == QMetaObject::QueryPropertyStored) {
        _id -= 3;
    } else if (_c == QMetaObject::QueryPropertyEditable) {
        _id -= 3;
    } else if (_c == QMetaObject::QueryPropertyUser) {
        _id -= 3;
    }
#endif // QT_NO_PROPERTIES
    return _id;
}
QT_WARNING_POP
QT_END_MOC_NAMESPACE
