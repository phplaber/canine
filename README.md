**canine** 是一款发现 Android 设备文件系统攻击面的工具，用于辅助漏洞挖掘。

## 功能
1. 在指定目录下，搜索指定用户或用户组具有可写权限的文件
2. 在指定目录下，搜索指定用户以 `root` 或 `system` 特权用户身份执行的 `SUID` 可执行文件 
3. 在指定目录下，搜索指定用户组以 `root` 或 `system` 特权用户组身份执行的 `SGID` 可执行文件 

## 使用
1. 下载源代码
``` bash
git clone https://github.com/phplaber/canine.git
```
2. 编译程序，生成可执行文件
``` bash
# 1、64位架构则为 arm64 
# 2、需使用 ndk 工具链 
# 3、需显示开启 CGO_ENABLED
env GOOS=android GOARCH=arm CC=/path/to/sdk/ndk/linux-androideabi-clang CGO_ENABLED=1 go build  -o ./bin/ -v ./cmd/canine
```
3. 使用 adb 将文件拷贝进 Android
``` bash
adb push /path/to/bin/canine /data/local/tmp/
```
4. 赋予文件执行权限，然后执行
``` bash
chmod +x canine
./canine --help                                    
canine v0.1
	A tool for find andriod attack surface of file system
	Usage: canine -u [user] -g [groups] dirpath1 dirpath2 ...
  -g string
    	groupname(s), e.g. shell,log,sdcard_rw
  -u string
    	username, e.g. shell
```
5. 使用示例
``` bash
./canine -u shell -g shell,log,adb,sdcard_rw /dev /data
[*] Scanning...
[*] Found 0 entries that are SUID executable
[*] Found 0 entries that are SGID executable
[*] Found 194 entries that are Writable
   chardev 0666 root root /dev/ashmem
   chardev 0666 root root /dev/binder
symbollink 0777 root root /dev/block/platform/bootdevice/by-name/boot
symbollink 0777 root root /dev/block/platform/bootdevice/by-name/boot_para
symbollink 0777 root root /dev/block/platform/bootdevice/by-name/cache
symbollink 0777 root root /dev/block/platform/bootdevice/by-name/cust
symbollink 0777 root root /dev/block/platform/bootdevice/by-name/devinfo
...
```