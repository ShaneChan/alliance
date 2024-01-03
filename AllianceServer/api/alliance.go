package api

import (
	"AllianceServer/mgo"
	"AllianceServer/predefine"
	"AllianceServer/util"
	"context"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sort"
	"strconv"
)

// ======================= api =======================

// AllianceList 查询公会列表
func AllianceList() (ret string, err error) {
	cur, err := mgo.Find("testing", "alliance", bson.M{})
	if err != nil {
		return "", err
	}

	var results []bson.M
	err = cur.All(context.Background(), &results)
	ret = "all alliances :\n"
	for _, result := range results {
		ret = ret + result["allianceName"].(string) + "\n"
	}
	return
}

func WhichAlliance(userAccount string) (string, error) {
	cur, err := mgo.FindOne("testing", "user", bson.M{"account": userAccount})
	if err != nil {
		return "", err
	}

	var result bson.M
	err = cur.Decode(&result)

	if err != nil {
		return "", err
	}

	transResult := result["allianceName"]
	if result["allianceName"] == nil {
		return "no alliance", nil
	}

	cur1, err := mgo.Find("testing", "user", bson.M{"allianceName": transResult})
	ret := "alliance:" + transResult.(string) + "\n" + "members:"

	var results []bson.M
	err = cur1.All(context.Background(), &results)
	for _, result := range results {
		ret = ret + result["account"].(string) + " "
	}

	return ret, nil
}

// CreateAlliance 创建公会
func CreateAlliance(userAccount string, alliance string) error {
	allianceName, _ := allianceName(alliance)

	if allianceName != "" {
		return errors.New("already have alliance")
	}

	var items []bson.M
	for i := 1; i <= 5; i++ { // 初始化五种道具，数量都是1
		items = append(items, bson.M{"id": i, "name": "道具" + strconv.Itoa(i), "itemType": i, "number": 1})
	}

	for i := 6; i <= predefine.INITIAL_MAX_GRID_NUM; i++ {
		items = append(items, bson.M{"id": 0, "name": "道具" + strconv.Itoa(0), "itemType": 0, "number": 0})
	}
	_, err := mgo.UpdateOne("testing", "alliance", bson.M{"leader": userAccount},
		bson.M{"$set": bson.M{"leader": userAccount, "allianceName": alliance, "items": items}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	err = JoinAlliance(userAccount, alliance)

	return err
}

// JoinAlliance 加入公会
func JoinAlliance(userAccount string, alliance string) error {
	OldAllianceName, _ := allianceName(userAccount)

	if OldAllianceName != "" {
		return errors.New("already have alliance")
	}

	_, err := mgo.UpdateOne("testing", "user", bson.M{"account": userAccount}, bson.M{"$set": bson.M{"allianceName": alliance}}, options.Update().SetUpsert(true))

	return err
}

// DismissAlliance 解散公会（验证会长）
func DismissAlliance(userAccount string) error {
	allianceName, err := allianceName(userAccount)
	if err != nil {
		return err
	}

	isLeader, err := isLeader(userAccount, allianceName)
	if err != nil {
		return err
	}
	if !isLeader {
		return errors.New("not leader")
	}

	_, err = mgo.UpdateMany("testing", "user", bson.M{"allianceName": allianceName}, bson.M{"$set": bson.M{"allianceName": nil}}, options.Update().SetUpsert(true))
	_, err = mgo.DeleteOne("testing", "alliance", bson.M{"allianceName": allianceName})

	return err
}

// CommitItem 提交物品
func CommitItem(userAccount string, id int, num int) error {
	allianceName, err := allianceName(userAccount)
	if err != nil {
		return err
	}
	cur, err := mgo.FindOne("testing", "alliance", bson.M{"allianceName": allianceName})
	if err != nil {
		return err
	}
	var result bson.M
	err = cur.Decode(&result)
	if err != nil {
		return err
	}

	transResult := result["items"].(bson.A)
	var newResult []bson.M
	isOK := false
	for _, v := range transResult {
		v1 := v.(bson.M)
		OldNum := int(v1["number"].(int32))
		if !isOK && int(v1["id"].(int32)) == id && OldNum < predefine.ITEM_MAX_NUM_ONE_GRID {
			newNumber := util.Min(predefine.ITEM_MAX_NUM_ONE_GRID, OldNum+num)
			v1["number"] = newNumber
			newResult = append(newResult, v1)
			num -= predefine.ITEM_MAX_NUM_ONE_GRID - OldNum
			if num <= 0 {
				isOK = true
			}
		} else if !isOK && int(v1["id"].(int32)) == 0 {
			newResult = append(newResult, bson.M{"id": id, "name": "道具" + strconv.Itoa(id), "itemType": id, "number": util.Min(predefine.ITEM_MAX_NUM_ONE_GRID, num)})
			num -= predefine.ITEM_MAX_NUM_ONE_GRID - OldNum
			if num <= 0 {
				isOK = true
			}
		} else {
			newResult = append(newResult, v1)
		}
	}

	_, err = mgo.UpdateOne("testing", "alliance", bson.M{"allianceName": allianceName}, bson.M{"$set": bson.M{"items": newResult}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

// TidyItems 整理物品
func TidyItems(userAccount string) error {
	allianceName, err := allianceName(userAccount)
	if err != nil {
		return err
	}
	cur, err := mgo.FindOne("testing", "alliance", bson.M{"allianceName": allianceName})
	if err != nil {
		return err
	}
	var result bson.M
	err = cur.Decode(&result)
	if err != nil {
		return err
	}

	transResult := result["items"].(bson.A)
	statisticsMap := map[int32]int32{}
	length := len(transResult)

	for _, v := range transResult {
		itemType := v.(bson.M)["itemType"].(int32)
		itemNum := v.(bson.M)["number"].(int32)
		if itemType != 0 {
			if _, ok := statisticsMap[itemType]; ok {
				statisticsMap[itemType] += itemNum
			} else {
				statisticsMap[itemType] = itemNum
			}
		}
	}

	var itemTypes []int32

	for k := range statisticsMap {
		itemTypes = append(itemTypes, k)
	}

	sort.Slice(itemTypes, func(i, j int) bool {
		return itemTypes[i] < itemTypes[j]
	})
	var items []bson.M
	for _, itemType := range itemTypes {
		num := statisticsMap[itemType]
		for num > 0 {
			items = append(items, bson.M{"id": int(itemType), "name": "道具" + strconv.Itoa(int(itemType)), "itemType": int(itemType), "number": util.Min(5, int(num))})
			num = int32(util.Max(0, int(num)-5))
		}
	}

	length1 := len(items)

	for i := 0; i < length-length1; i++ {
		items = append(items, bson.M{"id": 0, "name": "道具" + strconv.Itoa(0), "itemType": 0, "number": util.Min(5, 0)})
	}
	_, err = mgo.UpdateOne("testing", "alliance", bson.M{"allianceName": allianceName}, bson.M{"$set": bson.M{"items": items}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

// IncreaseCapacity 扩容公会仓库（验证会长）
func IncreaseCapacity(userAccount string) error {
	allianceName, err := allianceName(userAccount)
	if err != nil {
		return err
	}

	isLeader, err := isLeader(userAccount, allianceName)
	if err != nil {
		return err
	}
	if !isLeader {
		return errors.New("not leader")
	}

	cur, err := mgo.FindOne("testing", "alliance", bson.M{"allianceName": allianceName})
	if err != nil {
		return err
	}
	var result bson.M
	err = cur.Decode(&result)
	if err != nil {
		return err
	}

	transResult := result["items"].(bson.A)
	if len(transResult) >= predefine.FINAL_MAX_GRID_NUM {
		return errors.New("max grid num")
	}
	var newResult []bson.M

	for _, v := range transResult {
		newResult = append(newResult, v.(bson.M))
	}

	for i := 0; i < predefine.INCREASE_NUM; i++ {
		newResult = append(newResult, bson.M{"id": 0, "name": "道具" + strconv.Itoa(0), "itemType": 0, "number": 0})
	}

	_, err = mgo.UpdateOne("testing", "alliance", bson.M{"allianceName": allianceName}, bson.M{"$set": bson.M{"items": newResult}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

// DestroyItem 销毁格子物品（验证会长）
func DestroyItem(userAccount string, index int) error {
	allianceName, err := allianceName(userAccount)
	if err != nil {
		return err
	}
	cur, err := mgo.FindOne("testing", "alliance", bson.M{"allianceName": allianceName})
	if err != nil {
		return err
	}
	var result bson.M
	err = cur.Decode(&result)
	if err != nil {
		return err
	}

	transResult := result["items"].(bson.A)
	leader := result["leader"].(string)
	if leader != userAccount {
		return errors.New("not leader")
	}
	var newResult []bson.M
	if index < 1 || index > len(transResult) {
		return errors.New("index error")
	}
	for k, v := range transResult {
		if k == index-1 {
			newResult = append(newResult, bson.M{"id": 0, "name": "道具" + strconv.Itoa(0), "itemType": 0, "number": 0})
		} else {
			newResult = append(newResult, v.(bson.M))
		}
	}

	_, err = mgo.UpdateOne("testing", "alliance", bson.M{"allianceName": allianceName}, bson.M{"$set": bson.M{"items": newResult}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

// AllianceItems 查看公会物品
func AllianceItems(userAccount string) (string, error) {
	allianceName, err := allianceName(userAccount)
	if err != nil {
		return "", err
	}
	cur, err := mgo.FindOne("testing", "alliance", bson.M{"allianceName": allianceName})
	if err != nil {
		return "", err
	}
	var result bson.M
	err = cur.Decode(&result)
	if err != nil {
		return "", err
	}

	transResult := result["items"].(bson.A)
	var retString string
	for k, v := range transResult {
		jsonStr, _ := json.Marshal(v)
		retString = retString + strconv.Itoa(k+1) + "." + string(jsonStr) + "\n"
	}

	return retString, nil
}

// ===================== common =====================

// 获取公会名
func allianceName(userAccount string) (string, error) {
	cur1, err := mgo.FindOne("testing", "user", bson.M{"account": userAccount})
	if err != nil {
		return "", err
	}

	var result1 bson.M
	err = cur1.Decode(&result1)
	if err != nil {
		return "", err
	}
	if result1["allianceName"] == nil {
		return "", errors.New("no alliance")
	}

	allianceName := result1["allianceName"].(string)
	return allianceName, nil
}

// 是否是指定公会的会长
func isLeader(userAccount string, allianceName string) (bool, error) {
	cur, err := mgo.FindOne("testing", "alliance", bson.M{"allianceName": allianceName})
	if err != nil {
		return false, err
	}
	var result bson.M
	err = cur.Decode(&result)
	if err != nil {
		return false, err
	}

	leader := result["leader"].(string)
	return userAccount == leader, nil
}
