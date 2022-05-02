package cloudflare

import (
	"context"
	"fmt"
	"log"
	"strings"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudflareZoneCacheVariants() *schema.Resource {
	return &schema.Resource{
		Schema: resourceCloudflareZoneCacheVariantsSchema(),
		CreateContext: resourceCloudflareZoneCacheVariantsUpdate,
		ReadContext: resourceCloudflareZoneCacheVariantsRead,
		UpdateContext: resourceCloudflareZoneCacheVariantsUpdate,
		DeleteContext: resourceCloudflareZoneCacheVariantsDelete,
	}
}

func resourceCloudflareZoneCacheVariantsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudflare.API)

	log.Printf("[INFO] Reading Zone Cache Variants in zone %q", d.Id())

	zoneCacheVariants, err := client.ZoneCacheVariants(context.Background(), d.Id())

	if err != nil {
		if strings.Contains(err.Error(), "HTTP status 404") {
			log.Printf("[INFO] Zone Cache Variants for zone %q not found", d.Id())
			d.SetId("")
			return nil
		} else {
			return fmt.Errorf("Error reading cache variants for zone %q: %w", d.Id(), err)
		}
	}

	value := zoneCacheVariants.Value

	if err := d.Set("avif", value.Avif); err != nil {
		return fmt.Errorf("failed to set avif: %w", err)
	}

	if err := d.Set("bmp", value.Bmp); err != nil {
		return fmt.Errorf("failed to set bmp: %w", err)
	}

	if err := d.Set("gif", value.Gif); err != nil {
		return fmt.Errorf("failed to set gif: %w", err)
	}

	if err := d.Set("jpeg", value.Jpeg); err != nil {
		return fmt.Errorf("failed to set jpeg: %w", err)
	}

	if err := d.Set("jpg", value.Jpg); err != nil {
		return fmt.Errorf("failed to set jpg: %w", err)
	}

	if err := d.Set("jp2", value.Jp2); err != nil {
		return fmt.Errorf("failed to set jp2: %w", err)
	}

	if err := d.Set("jpg2", value.Jpg2); err != nil {
		return fmt.Errorf("failed to set jpg2: %w", err)
	}

	if err := d.Set("png", value.Png); err != nil {
		return fmt.Errorf("failed to set png: %w", err)
	}

	if err := d.Set("tif", value.Tif); err != nil {
		return fmt.Errorf("failed to set tif: %w", err)
	}

	if err := d.Set("tiff", value.Tiff); err != nil {
		return fmt.Errorf("failed to set tiff: %w", err)
	}

	if err := d.Set("webp", value.Webp); err != nil {
		return fmt.Errorf("failed to set webp: %w", err)
	}

	return nil
}

func resourceCloudflareZoneCacheVariantsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudflare.API)

	zoneID := d.Get("zone_id").(string)
	d.SetId(zoneID)

	variantsValue := cacheVariantsValuesFromResource(d)
	log.Printf("[INFO] Setting Zone Cache Variants to struct: %+v for zone ID: %q", variantsValue, d.Id())

	_, err := client.UpdateZoneCacheVariants(context.Background(), d.Id(), variantsValue)

	if err != nil {
		return fmt.Errorf("error setting cache variants for zone %q: %w", d.Id(), err)
	}

	return resourceCloudflareZoneCacheVariantsRead(d, meta)
}

func resourceCloudflareZoneCacheVariantsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudflare.API)

	log.Printf("[INFO] Deleting Zone Cache Variants for zone ID: %q", d.Id())

	err := client.DeleteZoneCacheVariants(context.Background(), d.Id())

	if err != nil {
		return fmt.Errorf("error deleting cache variants for zone %v: %w", d.Id(), err)
	}

	return nil
}

func cacheVariantsValuesFromResource(d *schema.ResourceData) cloudflare.ZoneCacheVariantsValues {
	variantsValue := cloudflare.ZoneCacheVariantsValues{}

	if value, ok := d.GetOk("avif"); ok {
		variantsValue.Avif = expandInterfaceToStringList(value.(*schema.Set).List())
	}

	if value, ok := d.GetOk("bmp"); ok {
		variantsValue.Bmp = expandInterfaceToStringList(value.(*schema.Set).List())
	}

	if value, ok := d.GetOk("gif"); ok {
		variantsValue.Gif = expandInterfaceToStringList(value.(*schema.Set).List())
	}

	if value, ok := d.GetOk("jpeg"); ok {
		variantsValue.Jpeg = expandInterfaceToStringList(value.(*schema.Set).List())
	}

	if value, ok := d.GetOk("jpg"); ok {
		variantsValue.Jpg = expandInterfaceToStringList(value.(*schema.Set).List())
	}

	if value, ok := d.GetOk("jp2"); ok {
		variantsValue.Jp2 = expandInterfaceToStringList(value.(*schema.Set).List())
	}

	if value, ok := d.GetOk("jpg2"); ok {
		variantsValue.Jpg2 = expandInterfaceToStringList(value.(*schema.Set).List())
	}

	if value, ok := d.GetOk("png"); ok {
		variantsValue.Png = expandInterfaceToStringList(value.(*schema.Set).List())
	}

	if value, ok := d.GetOk("tif"); ok {
		variantsValue.Tif = expandInterfaceToStringList(value.(*schema.Set).List())
	}

	if value, ok := d.GetOk("tiff"); ok {
		variantsValue.Tiff = expandInterfaceToStringList(value.(*schema.Set).List())
	}

	if value, ok := d.GetOk("webp"); ok {
		variantsValue.Webp = expandInterfaceToStringList(value.(*schema.Set).List())
	}

	return variantsValue
}
