package stock

import (
    "gopkg.in/mgo.v2/bson"
    "reflect"
)

type PointField struct {
    Value    float64
    Change   float64
    Position float64
}

type HistoryStockPoint struct {
    Id           bson.ObjectId `bson:"_id"`
    //SrcRecordId                  bson.ObjectId `bson:"srcRecordId"`
    Ts           float64
    Date         string
    Amplitude    float64

    //PosInYear                    float64 `bson:"posInYear"`
    //OpenThanPreClose             float64 `bson:"openThanPreClose"`
    //CloseThanOpen                float64 `bson:"closeThanOpen"`

    Open         PointField
    Open_ma_5    PointField
    Open_ma_10   PointField
    Open_ma_30   PointField
    Open_ma_60   PointField
    Open_ma_90   PointField

    Close        PointField
    Close_ma_5   PointField
    Close_ma_10  PointField
    Close_ma_30  PointField
    Close_ma_60  PointField
    Close_ma_90  PointField

    High         PointField
    High_ma_5    PointField
    High_ma_10   PointField
    High_ma_30   PointField
    High_ma_60   PointField
    High_ma_90   PointField

    Low          PointField
    Low_ma_5     PointField
    Low_ma_10    PointField
    Low_ma_30    PointField
    Low_ma_60    PointField
    Low_ma_90    PointField

    Avg          PointField
    Avg_ma_5     PointField
    Avg_ma_10    PointField
    Avg_ma_30    PointField
    Avg_ma_60    PointField
    Avg_ma_90    PointField

    Volume       PointField
    Volume_ma_5  PointField
    Volume_ma_10 PointField
    Volume_ma_30 PointField
    Volume_ma_60 PointField
    Volume_ma_90 PointField
}

func (me *HistoryStockPoint) GetValueByFieldName(n string) (value interface{}, isValid bool) {
    //fmt.Println(me, n)
    r := reflect.ValueOf(me)
    f := reflect.Indirect(r).FieldByName(n)
    isValid = f.IsValid()
    value = f.Interface()
    return
}

func (me *HistoryStockPoint) SetFloatByFieldName(n string, value float64) (ok bool) {
    r := reflect.ValueOf(me)
    f := r.Elem().FieldByName(n)
    if f.IsValid() && f.CanSet() && f.Kind() == reflect.Float64 {
        f.SetFloat(value)
    }
    return
}
