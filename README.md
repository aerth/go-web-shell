# go-web-shell

this implements a simple web shell (+web server).

listens on port 8080 by default

accepts simple shell commands, no terminal features

meant for single user

## building


````shell
$ make
````

## cross compilation for android or raspberry pi (ARM)

step 1: run on build host (linux)

```

make arm && adb push web-shell /sdcard/

```

step 2: install on an android (adb shell)

```
su
mount -o remount,rw /system
mv /sdcard/web-shell /system/xbin/

```

## license

[MIT](https://github.com/matiasinsaurralde/go-web-shell/blob/master/LICENSE)
