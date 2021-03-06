package cloudsigma

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cloudsigma/cloudsigma-sdk-go/cloudsigma"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCloudSigmaDrive() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudSigmaDriveCreate,
		ReadContext:   resourceCloudSigmaDriveRead,
		UpdateContext: resourceCloudSigmaDriveUpdate,
		DeleteContext: resourceCloudSigmaDriveDelete,

		SchemaVersion: 0,

		Schema: map[string]*schema.Schema{
			"clone_drive_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"media": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"cdrom", "disk"}, false)),
			},

			"mounted_on": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_uri": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.NoZeroValues),
			},

			"resource_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"size": {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(536870912)), // 536870912 = 512MB
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"storage_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the drive resource",
			},
		},
	}
}

func resourceCloudSigmaDriveCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudsigma.Client)

	cloneDriveUUID := d.Get("clone_drive_id").(string)
	if cloneDriveUUID != "" {
		// Clone the Drive if 'clone_drive_id' is set
		cloneRequest := &cloudsigma.DriveCloneRequest{
			Drive: &cloudsigma.Drive{
				Media: d.Get("media").(string),
				Name:  d.Get("name").(string),
				Size:  d.Get("size").(int),
			},
		}
		log.Printf("[DEBUG] Drive clone configuration: %#+v", *cloneRequest)
		drive, _, err := client.Drives.Clone(ctx, cloneDriveUUID, cloneRequest)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(drive.UUID)
		log.Printf("[INFO] Drive ID: %s", d.Id())
	} else {
		// Create the Drive because 'clone_drive_id' is not set
		createRequest := &cloudsigma.DriveCreateRequest{
			Drives: []cloudsigma.Drive{
				{
					Media: d.Get("media").(string),
					Name:  d.Get("name").(string),
					Size:  d.Get("size").(int),
				},
			},
		}
		log.Printf("[DEBUG] Drive create configuration: %#v", *createRequest)
		drives, _, err := client.Drives.Create(ctx, createRequest)
		if err != nil {
			return diag.FromErr(err)
		}
		drive := &drives[0]

		d.SetId(drive.UUID)
		log.Printf("[INFO] Drive ID: %s", d.Id())
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"cloning_dst", "creating"},
		Target:     []string{"unmounted"},
		Refresh:    driveStateRefreshFunc(ctx, client, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	return resourceCloudSigmaDriveRead(ctx, d, meta)
}

func resourceCloudSigmaDriveRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudsigma.Client)

	// Refresh the Drive state
	drive, resp, err := client.Drives.Get(ctx, d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("media", drive.Media)
	_ = d.Set("name", drive.Name)
	_ = d.Set("resource_uri", drive.ResourceURI)
	_ = d.Set("size", drive.Size)
	_ = d.Set("status", drive.Status)
	_ = d.Set("storage_type", drive.StorageType)
	_ = d.Set("uuid", drive.UUID)

	if err := d.Set("mounted_on", flattenMountedOn(&drive.MountedOn)); err != nil {
		return diag.Errorf("[DEBUG] Error setting Drive mounted_on - error: %#v", err)
	}

	return nil
}

func resourceCloudSigmaDriveUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudsigma.Client)

	drive := &cloudsigma.Drive{
		UUID: d.Id(),
	}

	if media, ok := d.GetOk("media"); ok {
		drive.Media = media.(string)
	}
	if name, ok := d.GetOk("name"); ok {
		drive.Name = name.(string)
	}
	if size, ok := d.GetOk("size"); ok {
		drive.Size = size.(int)
	}

	updateRequest := &cloudsigma.DriveUpdateRequest{
		Drive: drive,
	}
	log.Printf("[DEBUG] Drive update configuration: %#v", *updateRequest)

	_, _, err := client.Drives.Update(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCloudSigmaDriveRead(ctx, d, meta)
}

func resourceCloudSigmaDriveDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudsigma.Client)

	if v, ok := d.GetOk("mounted_on"); ok {
		mountedOns, err := expandMountedOn(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}

		for _, mountedOn := range mountedOns {
			err := stopServer(ctx, client, mountedOn.UUID)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	_, err := client.Drives.Delete(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func driveStateRefreshFunc(ctx context.Context, client *cloudsigma.Client, uuid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		drive, _, err := client.Drives.Get(ctx, uuid)
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving drive with uuid %s: %s", uuid, err)
		}

		return drive, drive.Status, nil
	}
}

func expandMountedOn(config []interface{}) ([]cloudsigma.ResourceLink, error) {
	mountedOns := make([]cloudsigma.ResourceLink, 0, len(config))

	for _, res := range config {
		mountedOn := res.(map[string]interface{})

		m := cloudsigma.ResourceLink{
			ResourceURI: mountedOn["resource_uri"].(string),
			UUID:        mountedOn["uuid"].(string),
		}

		mountedOns = append(mountedOns, m)
	}

	return mountedOns, nil
}

func flattenMountedOn(mountedOns *[]cloudsigma.ResourceLink) []interface{} {
	if mountedOns != nil {
		mos := make([]interface{}, len(*mountedOns))

		for i, mountedOn := range *mountedOns {
			mo := make(map[string]interface{})

			mo["resource_uri"] = mountedOn.ResourceURI
			mo["uuid"] = mountedOn.UUID

			mos[i] = mo
		}

		return mos
	}

	return make([]interface{}, 0)
}
