package processing

import (
    "ana/module"
    "strings"
    "time"
    "fmt"
    "gopkg.in/mgo.v2/bson"
    "strconv"
    "gopkg.in/mgo.v2"
)

var fields = []string{"open", "close", "high", "low", "volume", "mean", "size"}
var steps = []int{5, 10, 30, 60, 90}

func Do() {

    db := module.MongoSession.DB("dongsi")

    // =============== 2016.12.29 ================

    names, _ := db.CollectionNames()
    validNames := []string{}
    for _, name := range names {
        if !strings.Contains(name, "code_") {
            continue
        }
        validNames = append(validNames, name)
    }

    startTime := time.Now()
    handle(db, validNames, startTime)

    fmt.Println("done")

}

func interface2float(item interface{}) float64 {
    switch item.(type) {
    case int:
        tem := item.(int)
        return float64(tem)
    case float64:
        return item.(float64)
    default:
        return float64(0.0)
    }
}

func createPoint(statics map[string]float64, dataItem bson.M, preDataItem bson.M) (point bson.M) {

    point = bson.M{}

    for _, step := range steps {
        stepStr := strconv.Itoa(step)

        for k, v := range dataItem {
            point[k] = v
        }

        for _, f := range fields {
            if f != "size" {
                k := f + "_ma_" + stepStr
                //pointV := point[f].(float64)
                pointV := interface2float(point[f])
                pointKV := (pointV + statics[k]) / (statics["size_ma_" + stepStr] + 1)
                point[k] = pointKV

                //fmt.Println(k, pointV, statics[k], statics["size_ma_" + stepStr], pointKV)

                if preDataItem != nil {
                    preV := interface2float(preDataItem[k])
                    point[k + "_changeRate"] = (pointKV - preV) / preV
                }

            }
        }
    }

    // 特殊指标
    // 乖离率
    closeValue := interface2float(point["close"])
    close_ma_5 := interface2float(point["close_ma_5"])
    close_ma_10 := interface2float(point["close_ma_10"])

    bias_5 := (closeValue - close_ma_5) / close_ma_5;
    bias_10 := (closeValue - close_ma_10) / close_ma_10;
    point["bias_5"] = bias_5
    point["bias_10"] = bias_10


    if preDataItem != nil {
        point["isExtended"] = true
    }

    return
}

func reinitBucketStatics(statics map[string]float64, bucket []bson.M, bucketItem bson.M) (map[string]float64, []bson.M) {
    bucket = append(bucket, bucketItem)
    bucketSize := len(bucket)
    for _, step := range steps {
        stepStr := strconv.Itoa(step)
        oneItem := bucket[bucketSize - step]
        for _, f := range fields {
            v := interface2float(oneItem[f])
            statics[f + "_ma_" + stepStr] -= v
        }
    }
    bucket = bucket[1:]
    return statics, bucket
}

func initBucketStatics(bucket []bson.M) (statics map[string]float64) {
    statics = initSumCache()
    size := len(bucket)

    //maxs := map[string]float64{}
    //mins := map[string]float64{}

    for idx, item := range bucket {
        for _, f := range fields {
            v := interface2float(item[f])
            for _, step := range steps {
                if (size - idx < step) {
                    stepStr := strconv.Itoa(step)
                    statics[f + "_ma_" + stepStr] += v
                }
            }
        }
    }
    return
}
func initSumCache() map[string]float64 {
    sc := make(map[string]float64)
    for _, f := range fields {
        for _, s := range steps {
            sStr := strconv.Itoa(s)
            key := f + "_ma_" + sStr
            sc[key] = 0
        }
    }
    return sc
}
func createInitBucketCondition(date string, limit int) []bson.M {
    groupCondition := bson.M{
        "_id": "$date",
        "size": bson.M{
            "$sum": 1,
        },
    }

    for _, field := range fields {
        if field != "size" {
            groupCondition[field] = bson.M{
                "$sum": "$" + field,
            }
        }
    }

    conditions := []bson.M{
        bson.M{
            "$match": bson.M{
                "date": bson.M{
                    "$lt": date,
                },
            },
        },
        bson.M{
            "$group": groupCondition,
        },
        bson.M{
            "$sort": bson.M{
                "_id": 1,
            },
        },
        bson.M{
            "$limit": limit,
        },
    }

    return conditions
}

func handle(db *mgo.Database, names []string, startTime time.Time) {
    // init basic ratios
    namesSize := len(names)
    for nameIdx, name := range names {
        fmt.Println(name, "(", nameIdx, "/", namesSize, ")", time.Since(startTime).String())
        var dataSet []bson.M
        db.C(name).Find(bson.M{}).Sort("ts").All(&dataSet)

        var bucket []bson.M

        // 根据date去重，计算应该从哪个位置开始处理数据
        dateMap := make(map[string]bool)
        dateCount := 0
        dateCountLast := 0

        // 各个field的sum
        //sumCache := initSumCache()

        var theLastPoint bson.M = nil

        var bucketStatics map[string]float64
        for _, dataItem := range dataSet {

            currentDate := dataItem["date"].(string)

            if _, ok := dateMap[currentDate]; !ok {
                dateMap[currentDate] = true
                dateCount += 1
            }

            if dateCount >= 90 {
                if dateCount == 90 {
                    // 刚跑到90的时候，初始化一下需要用到的统计数据
                    conditions := createInitBucketCondition(currentDate, 90)
                    db.C(name).Pipe(conditions).AllowDiskUse().All(&bucket)
                    bucketStatics = initBucketStatics(bucket)
                } else if (dateCountLast != dateCount) {
                    // 每跳到新的一天，要load一下前一天的数据
                    conditions := createInitBucketCondition(currentDate, 1)
                    var bucketItem bson.M
                    db.C(name).Pipe(conditions).AllowDiskUse().One(&bucketItem)
                    bucketStatics, bucket = reinitBucketStatics(bucketStatics, bucket, bucketItem)
                }

                //if dateCount == 91 {
                //    return
                //}

                //fmt.Println(bucketStatics)
                //return
                // 获取过去14天的最高最低
                //bucketSize := len(bucket)
                //before14Date := bucket[bucketSize - 14]["_id"].(string)
                //currentDate := dataItem["date"].(string)

                //fmt.Println(maxAndMin)
                //
                //return


                point := createPoint(bucketStatics, dataItem, theLastPoint)

                theLastPoint = point
                //// 特殊ratio
                //var maxAndMin bson.M
                //db.C(name).Pipe([]bson.M{
                //    bson.M{
                //        "$match": bson.M{
                //            "date": bson.M{
                //                "$gte": before14Date,
                //                "$lt": currentDate,
                //            },
                //        },
                //    },
                //    bson.M{
                //        "$group": bson.M{
                //            "_id": nil,
                //            "maxHigh": bson.M{
                //                "$max": "$high",
                //            },
                //            "minLow": bson.M{
                //                "$min": "$low",
                //            },
                //        },
                //    },
                //}).AllowDiskUse().One(&maxAndMin)
                //
                ////point["max_high_14"] = maxAndMin["maxHigh"]
                ////point["min_low_14"] = maxAndMin["minLow"]
                //maxHigh := interface2float(maxAndMin["maxHigh"])
                //minLow := interface2float(maxAndMin["minLow"])
                //closeV := interface2float(point["close"])
                //point["WR"] = (maxHigh - closeV) / (maxHigh - minLow)


                db.C(name).UpdateId(point["_id"], point)
            }

            dateCountLast = dateCount

        }

        db.C(name).RemoveAll(bson.M{"isExtended":nil})

    }
}

