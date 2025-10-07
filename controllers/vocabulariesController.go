package controllers

var is_refreshed = map[bool]string{true: "yes", false: "no"}

var entryStatuses = map[string]string{
	"es_to_be_processed": "To Be Processed",
	"es_processed":       "Proccessed",
	"es_deaccessioned":   "Deaccessioned",
	"":                   "unknown",
}

func GetEntryStatuses() map[string]string { return entryStatuses }

func GetEntryStatus(s string) string {
	for k, v := range entryStatuses {
		if k == s {
			return v
		}
	}
	return "No Match"
}

var storageLocations = map[string]string{
	"sl_rsw_acm_born_digital": "RW ACM Born Digital",
	"sl_rsw_spec_coll":        "RW Special Collections",
	"sl_rsw_amatica_staging":  "RW Archivematica Staging",
	"sl_rStar":                "R*",
	"sl_cooper":               "806 Cooper FTK Workstation",
	"sl_bobst":                "806 Bobst FTK Workstation",
	"sl_fred":                 "ACM FRED FTK Workstation",
	"sl_wilma":                "ACM WILMA FTK Workstation",
	"sl_mac":                  "ACM Mac Workstation",
	"sl_not_imaged":           "Not Imaged",
	"sl_unknown":              "unknown",
}

func GetStorageLocations() map[string]string { return storageLocations }

func GetStorageLocation(s string) string {
	for k, v := range storageLocations {
		if k == s {
			return v
		}
	}
	return "No Match"
}

func GetMediatypes() map[string]string { return Mediatypes }

func GetMediaType(s string) string {

	for k, v := range Mediatypes {
		if k == s {
			return v
		}
	}
	return "No Match"
}

var Mediatypes = map[string]string{
	"":                           "",
	"mediatype_floppy_3_5":       "3.5 in. Floppy Disk",
	"mediatype_floppy_5_25":      "5.25 in. Floppy Disk",
	"mediatype_floppy_8":         "8 in. Floppy Disk",
	"mediatype_cd":               "CD commercial",
	"mediatype_cdr":              "CD-R",
	"mediatype_cdrw":             "CD-RW",
	"mediatype_data_cartridge":   "Data Cartridge",
	"mediatype_dvd":              "DVD commercial",
	"mediatype_dvdr":             "DVD-R",
	"mediatype_dvdrw":            "DVD-RW",
	"mediatype_file_transfer":    "File Transfer",
	"mediatype_flash_drive":      "Flash Drive",
	"mediatype_hard_disk_drive":  "Hard Disk Drive",
	"mediatype_jaz":              "Jaz Drive",
	"mediatype_laserdisc":        "Laserdisc",
	"mediatype_minidisc":         "MiniDisc",
	"mediatype_network_transfer": "Network Transfer",
	"mediatype_orb":              "Orb Disk",
	"mediatype_sd":               "SD Card",
	"mediatype_zip":              "Zip Disk",
}

func getInterfaces() map[string]string { return interfaces }

var interfaces = map[string]string{
	"":                           "",
	"interface_tableau_ultrabay": "Tableau Ultrabay",
	"interface_kryoflux":         "KryoFlux",
	"interface_tableau_t8r2":     "Tableau T8-R2",
	"interface_optical_HP":       "HP CD/DVD Drive",
}

func getInterface(s string) string {
	for k, y := range interfaces {
		if k == s {
			return y
		}
	}
	return ""
}

func getImagingSoftware() map[string]string { return imaging_software }

var imaging_software = map[string]string{
	"":                                      "",
	"imaging_software_kryoflux_imager_v220": "KryoFlux Imager (DTC 2.20)",
	"imaging_software_kryoflux_imager_v30":  "KryoFlux Imager (DTC 3.0)",
	"imaging_software_kryoflux_imager_v35":  "KryoFlux Imager (DTC 3.5)",
	"imaging_software_ftk_imager_v3146":     "FTK Imager (v3.1.4.6)",
	"imaging_software_ftk_imager_v42013":    "FTK Imager (v4.2.0.13)",
	"imaging_software_isobusterpro_v43":     "IsoBuster Pro (v4.3)",
	"imaging_software_eac_v13":              "Exact Audio Copy (v1.3)",
}

func getHDDInterfaces() map[string]string { return hdd_interfaces }

var hdd_interfaces = map[string]string{
	"":                    "",
	"hdd_interface_usb":   "USB",
	"hdd_interface_sata":  "SATA",
	"hdd_interface_fw400": "FW400",
	"hdd_interface_fw800": "FW800",
	"hdd_interface_scsi":  "SCSI",
	"hdd_interface_ide":   "IDE",
}

func getImageFormats() map[string]string { return image_formats }

var image_formats = map[string]string{
	"":                     "",
	"image_format_raw":     "Raw (dd)",
	"image_format_e01":     "E01",
	"image_format_ad1":     "AD1",
	"image_format_iso":     "ISO - Userspace",
	"image_format_iso_raw": "ISO - Raw",
	"image_format_bincue":  "BIN/CUE",
	"image_format_wavcue":  "WAV/CUE",
}

var encoding_schemes = map[string]string{
	"":                              "",
	"encoding_scheme_mfm":           "MFM",
	"encoding_scheme_apple_400_800": "Apple 400/800",
}

var filesystems = map[string]string{
	"":                       "",
	"filesystem_fat12":       "FAT12",
	"filesystem_fat16":       "FAT16",
	"filesystem_fat32":       "FAT32",
	"filesystem_ntfs":        "NTFS",
	"filesystem_hfs":         "HFS",
	"filesystem_hfs_plus":    "HFS+",
	"filesystem_amiga_os":    "AmigaOS",
	"filesystem_9660":        "ISO 9660",
	"filesystem_9660_joliet": "ISO 9660 Joliet",
	"filesystem_udf":         "UDF",
	"filesystem_uknown":      "Uknown",
}

func getImageSuccess() map[string]string { return image_success }

var image_success = map[string]string{
	"":                  "",
	"image_success_yes": "Yes",
	"image_success_no":  "No",
}

func getInterpretSuccess() map[string]string { return interpret_success }

var interpret_success = map[string]string{
	"":                             "",
	"interpret_success_yes":        "Yes",
	"interpret_success_yes_errors": "Yes W/Errors",
	"interpret_success_no":         "No",
}

func getStockUnits() map[string]string { return stock_unit }

var stock_unit = map[string]string{
	"":   "",
	"KB": "Kilobytes",
	"MB": "Megabytes",
	"GB": "Gigabytes",
	"TB": "Terabytes",
}

var accession_state = map[string]string{
	"":                             "",
	"accession_not_started":        "Not Started",
	"accession_queued":             "Queued",
	"accession_in_progress":        "In Progress",
	"accession_needs_qa":           "QA",
	"accession_ready_for_transfer": "Transferred",
	"acceession_complete":          "Archivesspace Updated /  Complete",
}

func getOpticalContentTypes() map[string]string { return content_type }

var content_type = map[string]string{
	"":              "",
	"content_video": "Video",
	"content_audio": "Audio",
	"content_data":  "Data",
	"content_email": "Email",
}

var structure = map[string]string{
	"":                   "",
	"structure_data":     "Data Disc",
	"structure_dvdvideo": "DVD-Video",
	"structure_cdda":     "Compact Disc Digital Audio",
	"structure_complex":  "Complex Optical Image",
}
