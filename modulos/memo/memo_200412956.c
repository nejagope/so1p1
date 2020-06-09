#include <linux/init.h>
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/sys.h>
#include <linux/mm.h>

static int memo_init(void){
    int factorBytes;
    long total_ram;
    long free_ram;
    int used_ram;
	struct sysinfo info;

	si_meminfo(&info);

	factorBytes = info.mem_unit;
	total_ram = info.totalram * factorBytes;
	free_ram = info.freeram * factorBytes;
	used_ram = ((total_ram - free_ram)* 100)/total_ram ;

	printk(KERN_ALERT "Carnets: 200412956 200412956\n");
	printk(KERN_ALERT "Cantidad total de memoria: %ld bytes\n", total_ram);
	printk(KERN_ALERT "Cantidad de memoria disponible: %ld bytes\n", free_ram);
	printk(KERN_ALERT "Memoria utilizada: %d por ciento\n", used_ram);
	return 0;
}

static void memo_exit(void){
    printk(KERN_ALERT "Sistemas Operativos 1\n");
}

module_init(memo_init);
module_exit(memo_exit);

//--------------
MODULE_AUTHOR("Nelson González");
MODULE_DESCRIPTION("[SO1] Memory module");
MODULE_LICENSE("GPL");
