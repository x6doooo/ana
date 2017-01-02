package stock

import "gopkg.in/mgo.v2/bson"

type ExtendPoint struct {
    Id                      bson.ObjectId `bson:"_id"`
    Ts                      float64
    Date                    string
    Bias_5                  float64
    Bias_10                 float64

    IsValid                 bool
    IsExtended              bool

    High                    float64
    HighChangeRate          float64
    High_ma_5               float64
    High_ma_5_changeRate    float64
    High_ma_10              float64
    High_ma_10_changeRate   float64
    High_ma_30              float64
    High_ma_30_changeRate   float64
    High_ma_60              float64
    High_ma_60_changeRate   float64
    High_ma_90              float64
    High_ma_90_changeRate   float64

    Low                     float64
    LowChangeRate           float64
    Low_ma_5                float64
    Low_ma_5_changeRate     float64
    Low_ma_10               float64
    Low_ma_10_changeRate    float64
    Low_ma_30               float64
    Low_ma_30_changeRate    float64
    Low_ma_60               float64
    Low_ma_60_changeRate    float64
    Low_ma_90               float64
    Low_ma_90_changeRate    float64

    Volume                  float64
    VolumeChangeRate        float64
    Volume_ma_5             float64
    Volume_ma_5_changeRate  float64
    Volume_ma_10            float64
    Volume_ma_10_changeRate float64
    Volume_ma_30            float64
    Volume_ma_30_changeRate float64
    Volume_ma_60            float64
    Volume_ma_60_changeRate float64
    Volume_ma_90            float64
    Volume_ma_90_changeRate float64

    Open                    float64
    OpenChangeRate          float64
    Open_ma_5               float64
    Open_ma_5_changeRate    float64
    Open_ma_10              float64
    Open_ma_10_changeRate   float64
    Open_ma_30              float64
    Open_ma_30_changeRate   float64
    Open_ma_60              float64
    Open_ma_60_changeRate   float64
    Open_ma_90              float64
    Open_ma_90_changeRate   float64

    Close                   float64
    CloseChangeRate         float64
    Close_ma_5              float64
    Close_ma_5_changeRate   float64
    Close_ma_10             float64
    Close_ma_10_changeRate  float64
    Close_ma_30             float64
    Close_ma_30_changeRate  float64
    Close_ma_60             float64
    Close_ma_60_changeRate  float64
    Close_ma_90             float64
    Close_ma_90_changeRate  float64

    Mean                    float64
    MeanChangeRate          float64
    Mean_ma_5               float64
    Mean_ma_5_changeRate    float64
    Mean_ma_10              float64
    Mean_ma_10_changeRate   float64
    Mean_ma_30              float64
    Mean_ma_30_changeRate   float64
    Mean_ma_60              float64
    Mean_ma_60_changeRate   float64
    Mean_ma_90              float64
    Mean_ma_90_changeRate   float64
}
