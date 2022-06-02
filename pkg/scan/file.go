package scan

import (
	"io/fs"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

// https://android.googlesource.com/platform/system/core/+/android-8.0.0_r4/libcutils/include/private/android_filesystem_config.h
var AID2Name = map[uint32]string{
	0:    "root",
	1000: "system",
	1001: "radio",
	1002: "bluetooth",
	1003: "graphics",
	1004: "input",
	1005: "audio",
	1006: "camera",
	1007: "log",
	1008: "compass",
	1009: "mount",
	1010: "wifi",
	1011: "adb",
	1012: "install",
	1013: "media",
	1014: "dhcp",
	1015: "sdcard_rw",
	1016: "vpn",
	1017: "keystore",
	1018: "usb",
	1019: "drm",
	1020: "mdnsr",
	1021: "gps",
	1022: "unused1",
	1023: "media_rw",
	1024: "mtp",
	1025: "unused2",
	1026: "drmrpc",
	1027: "nfc",
	1028: "sdcard_r",
	1029: "clat",
	1030: "loop_radio",
	1031: "media_drm",
	1032: "package_info",
	1033: "sdcard_pics",
	1034: "sdcard_av",
	1035: "sdcard_all",
	1036: "logd",
	1037: "shared_relro",
	1038: "dbus",
	1039: "tlsdate",
	1040: "media_ex",
	1041: "audioserver",
	1042: "metrics_coll",
	1043: "metricsd",
	1044: "webserv",
	1045: "debuggerd",
	1046: "media_codec",
	1047: "cameraserver",
	1048: "firewall",
	1049: "trunks",
	1050: "nvram",
	1051: "dns",
	1052: "dns_tether",
	1053: "webview_zygote",
	1054: "vehicle_network",
	1055: "media_audio",
	1056: "media_video",
	1057: "media_image",
	1058: "tombstoned",
	1059: "media_obb",
	1060: "ese",
	1061: "ota_update",
	2000: "shell",
	2001: "cache",
	2002: "diag",
	3001: "net_bt_admin",
	3002: "net_bt",
	3003: "inet",
	3004: "net_raw",
	3005: "net_admin",
	3006: "net_bw_stats",
	3007: "net_bw_acct",
	3009: "readproc",
	3010: "wakelock",
	9997: "everybody",
	9998: "misc",
	9999: "nobody",
}

func GetFileType(info fs.FileInfo) string {
	fms := info.Mode().String()

	short := "??" // 不可读的文件
	endPos := strings.Index(fms, "r")
	if endPos != -1 {
		short = fms[:endPos]
	}
	var filetype string
	// short -> filetype
	//  -: file
	//  d: directory
	//  L: symboliclink
	//  D: blockdev
	// Dc: chardev
	//  p: pipe
	//  S: socket
	switch short {
	case "-":
		filetype = "file"
	case "d":
		filetype = "directory"
	case "L":
		filetype = "symbollink"
	case "D":
		filetype = "blockdev"
	case "Dc":
		filetype = "chardev"
	case "p":
		filetype = "pipe"
	case "S":
		filetype = "socket"
	default:
		filetype = "other"
	}

	return filetype
}

func GetFilePerm(info fs.FileInfo) string {
	perm := info.Mode().Perm()

	return "0" + strconv.FormatUint(uint64(perm), 8)
}

func GetFileOwnership(info fs.FileInfo) (uint32, uint32, string, string) {
	var uid, gid uint32
	var owner, group string
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		uid = stat.Uid
		gid = stat.Gid

		// user: LookupId not implemented on android
		// user: LookupGroupId not implemented on android
		if runtime.GOOS == "android" {
			owner = AID2Name[uid]
			group = AID2Name[gid]
		} else {
			u := strconv.FormatUint(uint64(uid), 10)
			g := strconv.FormatUint(uint64(gid), 10)

			if uu, err := user.LookupId(u); err == nil {
				owner = uu.Username
			}
			if ug, err := user.LookupGroupId(g); err == nil {
				group = ug.Name
			}
		}
	}

	return uid, gid, owner, group
}
