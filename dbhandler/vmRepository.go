package dbhandler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/niluwats/bethel_dashboard/domain"
	"github.com/niluwats/bethel_dashboard/errs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func (d AuthRepositoryDb) SaveNode(vm domain.VmAll) (*domain.VirtualMachine, *errs.AppError) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)
	database := client.Database("bethel_dashboard")
	orgCollection := database.Collection("organizations")
	resCollection := database.Collection("resourcegroups")
	vmCollection := database.Collection("virtualmachines")
	//metCollection := database.Collection("metrics")

	or := domain.Organization{
		OrgId:   primitive.NewObjectID(),
		OrgName: "fcx",
	}

	rg := domain.ResourceGroup{
		Name:   "vmres",
		Region: "westus",
	}
	vmDetails := domain.VirtualMachine{
		VmName:     vm.VmName,
		VmUserName: vm.VmUserName,
		VmPassword: vm.VmPassword,
		IpAdd:      vm.IpAdd,
	}

	var filtered []bson.M
	filterCursor, err1 := orgCollection.Find(ctx, bson.M{"org_name": "fcx"})
	if err1 != nil {
		log.Fatal(err1)
	}
	if err = filterCursor.All(ctx, &filtered); err != nil {
		log.Fatal(err)
	}
	if filtered != nil {
		c := filtered[0]
		rg.OrgId = c["_id"]
	}
	if filtered == nil {
		insertResult, err := orgCollection.InsertOne(ctx, or)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertResult.InsertedID)
		rg.OrgId = insertResult.InsertedID
	}

	var filtered1 []bson.M
	filterCursor1, err1 := resCollection.Find(ctx, bson.M{"resourcegroup_name": "vmres"})
	if err1 != nil {
		log.Fatal(err1)
	}
	if err = filterCursor1.All(ctx, &filtered1); err != nil {
		log.Fatal(err)
	}
	if filtered1 != nil {
		c := filtered1[0]
		vmDetails.ResGrpId = c["_id"]
	}
	if filtered1 == nil {
		insertResult, err := resCollection.InsertOne(ctx, rg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertResult.InsertedID)
		vmDetails.ResGrpId = insertResult.InsertedID
	}

	var Filtered2 []bson.M
	filterCursor2, err1 := vmCollection.Find(ctx, bson.M{"vm_name": "vm1"})
	if err1 != nil {
		log.Fatal(err1)
	}
	if err = filterCursor2.All(ctx, &Filtered2); err != nil {
		log.Fatal(err)
	}
	if Filtered2 == nil {
		insertResult, err := vmCollection.InsertOne(ctx, vmDetails)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertResult.InsertedID)

	}
	return &vmDetails, nil
}

// vmLoginCred := domain.VmLogin{
// 	VmName:     vm.VmName,
// 	VmUserName: vm.VmUserName,
// 	VmPassword: vm.VmPassword,
// 	IpAdd:      vm.IpAdd,
// }
// resGrp := domain.ResourceGroup{
// 	Name:   vm.ResGrpName,
// 	Region: vm.Region,
// 	LoginDet: []domain.VmLogin{
// 		vmLoginCred,
// 	},
// }
// org := domain.Organization{
// 	OrgName:       vm.OrgName,
// 	ResourceGroup: []domain.ResourceGroup{resGrp},
// }
// loc := domain.Location{
// 	Region: vm.Region,
// }
// var resLoc domain.Location

// col := d.client.C("organizations")
// col2 := d.client.C("resourcegroup_locations")
// var res domain.Organization

// err0 := col2.Find(bson.M{"region": loc.Region}).One(&resLoc)
// if err0 == mgo.ErrNotFound {
// 	err0 = col2.Insert(&loc)
// 	if err0 != nil {
// 		logger.Error("error while inserting new resourcegroup location" + err0.Error())
// 		return nil, errs.NewUnexpectedError("unexpected DB error")
// 	}
// }

// err := col.Find(bson.M{"org_name": vm.OrgName}).One(&res)
// if err == mgo.ErrNotFound {
// 	err = col.Insert(&org)
// 	if err != nil {
// 		logger.Error("error while inserting new resource group" + err.Error())
// 		return nil, errs.NewUnexpectedError("unexpected DB error")
// 	}
// 	err = col2.Insert(&loc)
// 	if err != nil {
// 		logger.Error("error while inserting new region" + err.Error())
// 		return nil, errs.NewUnexpectedError("unexpected DB error")
// 	}
// } else {
// 	err = col.Find(bson.M{"org_name": vm.OrgName, "resourcegroup.resourcegroup_name": resGrp.Name}).One(&res)
// 	if err == mgo.ErrNotFound {
// 		pushQuery := bson.M{"resourcegroup": resGrp}
// 		err1 := col.Update(bson.M{"org_name": vm.OrgName}, bson.M{"$addToSet": pushQuery})
// 		if err1 != nil {
// 			logger.Error("error while updating resource group array" + err1.Error())
// 			return nil, errs.NewUnexpectedError("unexpected DB error")
// 		}

// 	} else {
// 		pushQuery := bson.M{"resourcegroup.$.virtual_machine": vmLoginCred}
// 		err1 := col.Update(bson.M{"org_name": vm.OrgName, "resourcegroup.resourcegroup_name": resGrp.Name}, bson.M{"$addToSet": pushQuery})
// 		if err1 != nil {
// 			logger.Error("error while updating resource group" + err1.Error())
// 			return nil, errs.NewUnexpectedError("unexpected DB error")
// 		}
// 	}
// }
// err = col.Find(bson.M{"org_name": vm.OrgName}).One(&res)
// if err != nil {
// 	logger.Error("error while fetching organization" + err.Error())
// 	return nil, errs.NewUnexpectedError("unexpected DB error")
// }
// return &res, nil
