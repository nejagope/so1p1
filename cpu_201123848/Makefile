obj-m += cpu_201123848.o
all:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules
clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) clean
test:
	dmesg -C
	insmod cpu_201123848.ko
	rmmod cpu_201123848.ko
 	dmesg