package policy

import (
	"log"
	"fmt"
	"os"

	intersight "github.com/CiscoDevNet/intersight-go"
)

func createSpanGroup() *intersight.StorageSpanDrives {
	spanGroup := intersight.NewStorageSpanDrives("storage.SpanDrives","storage.SpanDrives")
	spanGroup.SetSlots("2")
	return spanGroup   
}

func createVirtualDrivePolicy() *intersight.StorageVirtualDrivePolicy {
	virtualDrivePolicy := intersight.NewStorageVirtualDrivePolicy("storage.VirtualDrivePolicy", "storage.VirtualDrivePolicy")
	virtualDrivePolicy.SetStripSize(int32(64))
	virtualDrivePolicy.SetWritePolicy("Default")
	virtualDrivePolicy.SetReadPolicy("Default")
	virtualDrivePolicy.SetAccessPolicy("Default")
	virtualDrivePolicy.SetDriveCache("Default")
	return virtualDrivePolicy
}

func createStoragePolicyRelationship(moid string) intersight.StorageStoragePolicyRelationship {
	storagePolicy := new(intersight.StorageStoragePolicy)
	storagePolicy.ClassId = "mo.MoRef"
	storagePolicy.ObjectType = "storage.StoragePolicy"
	storagePolicy.Moid = &moid
	storageRelationship := intersight.StorageStoragePolicyAsStorageStoragePolicyRelationship(storagePolicy)
	return storageRelationship
}

func ReturnPolicyAbstractPolicyRelationshipStorageDriveGroup(config *Config, storageMoid sting) intersight.PolicyAbstractPolicyRelationship {
	moid := CreateSorageDriveGroup(config,storageMoid)
	storageDriveGroup1 := new(intersight.PolicyAbstractPolicy)
	storageDriveGroup1.SetClassId("mo.MoRef")
	storageDriveGroup1.ObjectType("storage.DriveGroup")
	storageDriveGroup1.SetMoid(moid)
	storageDriveGroupRelationship := intersight.PolicyAbstractPolicyAsPolicyAbstractPolicyRelationship(storageDriveGroup1)
	return storageDriveGroupRelationship
}

func CreateSorageDriveGroup(config *Config, storageMoid string) string {
	var err error
	cfg := getApiClient(config)
	apiClient := cfg.ApiClient   
	ctx := cfg.ctx
	storageDriveGroup := intersight.NewStorageDriveGroupWithDefaults()
	storageDriveGroup.SetName("storage_drive_gp_test")
	storageDriveGroup.SetRaidLevel("Raid0")
	storageDriveGroup.SetType(int32(0))
	spanGroup1 := createSpanGroup()
	spanGroups := []intersight.StorageSpanDrives{*spanGroup1}
	storageManualDriveGroup := intersight.NewStorageManualDriveGroupWithDefaults()
	storageManualDriveGroup.SetSpanGroups(spanGroups)
	manualDriveGroup := intersight.NewNullableStorageManualDriveGroup(storageManualDriveGroup)
	manualDriveGroup1 := manualDriveGroup.Get()
	storageDriveGroup.SetManualDriveGroup(*manualDriveGroup1)
	
	virtualDrivePolicy1 := createVirtualDrivePolicy()
	virtualPolicy := intersight.NewNullableStorageVirtualDrivePolicy(virtualDrivePolicy1)
	virtualPolicy1 := virtualPolicy.Get()
	
	virtualDrive1 := intersight.NewStorageVirtualDriveConfigurationWithDefaults()
	virtualDrive1.SetName("drive_gp_vd")
	virtualDrive1.SetSize(int64(100))
	virtualDrive1.SetExpandToAvailable(false)
	virtualDrive1.SetBootDrive(false)
	virtualDrive1.SetVirtualDrivePolicy(*virtualPolicy1)
	
	virtualDrive2 := intersight.NewStorageVirtualDriveConfigurationWithDefaults()
	virtualDrive2.SetName("drive_gp_vd_01")
	virtualDrive2.SetSize(int64(100))
	virtualDrive2.SetExpandToAvailable(false)
	virtualDrive2.SetBootDrive(false)
	virtualDrive2.SetVirtualDrivePolicy(*virtualPolicy1)
	virtualDrives := []intersight.StorageVirtualDriveConfiguration{*virtualDrive1, *virtualDrive2}
	storageDriveGroup.SetVirtualDrives(virtualDrives)

	storagePolicyRelationship := createStoragePolicyRelationship(storageMoid)
	storageDriveGroup.SetStoragePolicy(storagePolicyRelationship)

	ifMatch := ""
	ifNoneMatch := ""
	resp, r, err := apiClient.StorageApi.CreateStorageDriveGroup(ctx).StorageDriveGroup(*storageDriveGroup).IfMatch(ifMatch).IfNoneMatch(ifNoneMatch).Execute()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
		log.Fatalf("HTTP response: %v\n", r)
	}
	fmt.Fprintf(os.Stdout, "Response: %v\n", resp)
	moid := resp.GetMoid()
	return moid
}