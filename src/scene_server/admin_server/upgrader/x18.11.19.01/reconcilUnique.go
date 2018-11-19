/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */
package x18_11_19_01

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func reconcilUnique(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	type Attribute struct {
		ID                int64       `json:"id" bson:"id"`
		OwnerID           string      `json:"bk_supplier_account" bson:"bk_supplier_account"`
		ObjectID          string      `json:"bk_obj_id" bson:"bk_obj_id"`
		PropertyID        string      `json:"bk_property_id" bson:"bk_property_id"`
		PropertyName      string      `json:"bk_property_name" bson:"bk_property_name"`
		PropertyGroup     string      `json:"bk_property_group" bson:"bk_property_group"`
		PropertyGroupName string      `json:"bk_property_group_name" bson:"-"`
		PropertyIndex     int64       `json:"bk_property_index" bson:"bk_property_index"`
		Unit              string      `json:"unit" bson:"unit"`
		Placeholder       string      `json:"placeholder" bson:"placeholder"`
		IsEditable        bool        `json:"editable" bson:"editable"`
		IsPre             bool        `json:"ispre" bson:"ispre"`
		IsRequired        bool        `json:"isrequired" bson:"isrequired"`
		IsReadOnly        bool        `json:"isreadonly" bson:"isreadonly"`
		IsOnly            bool        `json:"isonly" bson:"isonly"`
		IsSystem          bool        `json:"bk_issystem" bson:"bk_issystem"`
		IsAPI             bool        `json:"bk_isapi" bson:"bk_isapi"`
		PropertyType      string      `json:"bk_property_type" bson:"bk_property_type"`
		Option            interface{} `json:"option" bson:"option"`
		Description       string      `json:"description" bson:"description"`
		Creator           string      `json:"creator" bson:"creator"`
		CreateTime        *time.Time  `json:"create_time" bson:"creaet_time"`
		LastTime          *time.Time  `json:"last_time" bson:"last_time"`
	}

	oldAttributes := []Attribute{}
	err := db.Table(common.BKTableNameObjAttDes).Find(nil).All(ctx, &oldAttributes)
	if err != nil {
		return err
	}
	var obj2IsOnlyProperty = map[string][]Attribute{}
	var propertyIDToProperty = map[string]Attribute{}

	var keyfunc = func(a, b string) string { return a + ":" + b }
	for _, oldAttr := range oldAttributes {
		if oldAttr.IsOnly {
			obj2IsOnlyProperty[oldAttr.ObjectID] = append(obj2IsOnlyProperty[oldAttr.ObjectID], oldAttr)
		}
		propertyIDToProperty[keyfunc(oldAttr.ObjectID, oldAttr.PropertyID)] = oldAttr
	}

	uniques := []metadata.ObjectUnique{
		// host
		{
			ObjID:     common.BKInnerObjIDHost,
			MustCheck: true,
			Keys: []metadata.UinqueKeys{
				{
					Kind: metadata.UinqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDHost, common.BKAssetIDField)].ID),
				},
			},
			Ispre:    true,
			OwnerID:  conf.OwnerID,
			LastTime: metadata.Now(),
		},
		{
			ObjID:     common.BKInnerObjIDHost,
			MustCheck: false,
			Keys: []metadata.UinqueKeys{
				{
					Kind: metadata.UinqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDHost, common.BKCloudIDField)].ID),
				},
				{
					Kind: metadata.UinqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDHost, common.BKHostInnerIPField)].ID),
				},
			},
			Ispre:    true,
			OwnerID:  conf.OwnerID,
			LastTime: metadata.Now(),
		},
		// process
		{
			ObjID:     common.BKInnerObjIDProc,
			MustCheck: true,
			Keys: []metadata.UinqueKeys{
				{
					Kind: metadata.UinqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDProc, common.BKAppIDField)].ID),
				},
				{
					Kind: metadata.UinqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDProc, common.BKProcNameField)].ID),
				},
			},
			Ispre:    true,
			OwnerID:  conf.OwnerID,
			LastTime: metadata.Now(),
		},
		{
			ObjID:     common.BKInnerObjIDProc,
			MustCheck: false,
			Keys: []metadata.UinqueKeys{
				{
					Kind: metadata.UinqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDProc, common.BKAppIDField)].ID),
				},
				{
					Kind: metadata.UinqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDProc, common.BKFuncIDField)].ID),
				},
			},
			Ispre:    true,
			OwnerID:  conf.OwnerID,
			LastTime: metadata.Now(),
		},
		// biz
		{
			ObjID:     common.BKInnerObjIDApp,
			MustCheck: true,
			Keys: []metadata.UinqueKeys{
				{
					Kind: metadata.UinqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDApp, common.BKAppNameField)].ID),
				},
			},
			Ispre:    true,
			OwnerID:  conf.OwnerID,
			LastTime: metadata.Now(),
		},
		// set
		{
			ObjID:     common.BKInnerObjIDSet,
			MustCheck: true,
			Keys: []metadata.UinqueKeys{
				{
					Kind: metadata.UinqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDSet, common.BKAppIDField)].ID),
				},
				{
					Kind: metadata.UinqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDSet, common.BKSetNameField)].ID),
				},
			},
			Ispre:    true,
			OwnerID:  conf.OwnerID,
			LastTime: metadata.Now(),
		},
		// module
		{
			ObjID:     common.BKInnerObjIDModule,
			MustCheck: true,
			Keys: []metadata.UinqueKeys{
				{
					Kind: metadata.UinqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDModule, common.BKAppIDField)].ID),
				},
				{
					Kind: metadata.UinqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDModule, common.BKModuleNameField)].ID),
				},
			},
			Ispre:    true,
			OwnerID:  conf.OwnerID,
			LastTime: metadata.Now(),
		},
		// cloud area
		{
			ObjID:     common.BKInnerObjIDPlat,
			MustCheck: true,
			Keys: []metadata.UinqueKeys{
				{
					Kind: metadata.UinqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDPlat, common.BKCloudNameField)].ID),
				},
			},
			Ispre:    true,
			OwnerID:  conf.OwnerID,
			LastTime: metadata.Now(),
		},
	}

	for objID, oldAttrs := range obj2IsOnlyProperty {
		keys := []metadata.UinqueKeys{}
		ownerID := conf.OwnerID
		allPreset := true
		for _, oldAttr := range oldAttrs {
			keys = append(keys, metadata.UinqueKeys{
				Kind: metadata.UinqueKeyKindProperty,
				ID:   uint64(oldAttr.ID),
			})
			ownerID = oldAttr.OwnerID
			if !oldAttr.IsPre || (oldAttr.IsPre && oldAttr.PropertyID == common.BKInstNameField) {
				allPreset = false
			}
		}
		if allPreset {
			continue
		}

		unique := metadata.ObjectUnique{
			ObjID:     objID,
			MustCheck: true,
			Keys:      keys,
			Ispre:     false,
			OwnerID:   ownerID,
			LastTime:  metadata.Now(),
		}
		uniques = append(uniques, unique)
	}

	for _, unique := range uniques {
		uid, err := db.NextSequence(ctx, common.BKTableNameObjUnique)
		if err != nil {
			return err
		}
		unique.ID = uid
		if err := db.Table(common.BKTableNameObjUnique).Insert(ctx, unique); err != nil {
			return err
		}
	}

	if err := db.Table(common.BKTableNameObjAttDes).DropColumn(ctx, common.BKIsOnlyField); err != nil {
		return err
	}

	return nil
}
