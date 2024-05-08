package controllers

var is_refreshed = map[bool]string{true: "yes", false: "no"}

var PartnerCodes = map[string]string{
	"":            "",
	"fales":       "Fales Library & Special Collections",
	"nyuarchives": "NYU Archives",
	"tamwag":      "Tamiment Library and Robert F. Wagner Labor Archives",
	"abudhabi":    "Abu Dhabi",
}

var filename_partner_codes = map[string]string{
	"fa": "fales",
	"tw": "tamwag",
	"ua": "nyuarchives",
}

func getMediatypes() map[string]string { return mediatypes }

var mediatypes = map[string]string{
	"":                          "",
	"mediatype_transfer":        "Network Transfer",
	"mediatype_floppy_3_5":      "3.5 in. Floppy Disk",
	"mediatype_floppy_5_25":     "5.25 in. Floppy Disk",
	"mediatype_floppy_8":        "8 in. Floppy Disk",
	"mediatype_hard_disk_drive": "Hard Disk Drive",
	"mediatype_flash_drive":     "Flash Drive",
	"mediatype_dvd":             "DVD commercial",
	"mediatype_dvdr":            "DVD-R",
	"mediatype_dvdrw":           "DVD-RW",
	"mediatype_cd":              "CD commercial",
	"mediatype_cdr":             "CD-R",
	"mediatype_cdrw":            "CD-RW",
	"mediatype_cdda":            "Audio CD",
	"mediatype_jaz":             "Jaz Drive",
	"mediatype_zip":             "Zip Disk",
	"mediatype_sd":              "SD Card",
	"mediatype_minidisc":        "MiniDisc",
	"mediatype_data_cartridge":  "Data Cartridge",
	"mediatype_laserdisc":       "Laserdisc",
	"mediatype_orb":             "Orb Disk",
}

func GetMediaType(s string) string {

	for k, y := range mediatypes {
		if k == s {
			return y
		}
	}
	return "No Match"
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
	"":                 "",
	"image_sucess_yes": "Yes",
	"image_success_no": "No",
}

func getInterpretSuccess() map[string]string { return interpret_success }

var interpret_success = map[string]string{
	"":                             "",
	"interpret_success_yes":        "Yes",
	"interpret_success_yes_errors": "Yes W/Errors",
	"interpret_succes_no":          "No",
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
