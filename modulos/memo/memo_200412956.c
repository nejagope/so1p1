#include <linux/init.h>
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/sys.h>
#include <linux/mm.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>


static int memo_proc_show(struct seq_file *m, void *v) {
	//variables que contendran la información del estado de la memoria
	int factorBytes;
	long total_ram;
	long free_ram;
	int used_ram;
	struct sysinfo info;	

	//se obtiene la información del estado de la memoria
	si_meminfo(&info);

	//se rellenan las variables
	factorBytes = info.mem_unit;
	total_ram = info.totalram * factorBytes;
	free_ram = info.freeram * factorBytes;
	used_ram = ((total_ram - free_ram)* 100)/total_ram ;

	//se escribe el archivo
	seq_printf(m, "Carnets: 200412956 201123848\nNombres:Nelson Gonzalez, Mario Alvarado\nCantidad total de memoria: %ld bytes\nCantidad de memoria disponible: %ld bytes\nMemoria utilizada: %d por ciento\n", total_ram, free_ram, used_ram);
	return 0;
}

static int memo_proc_open(struct inode *inode, struct  file *file) {
     return single_open(file, memo_proc_show, NULL);
}

//operaciones a realizar
static struct file_operations memo_proc_fops = {
     .owner = THIS_MODULE,
     .open = memo_proc_open,
     .read = seq_read,
     .llseek = seq_lseek,
     .release = single_release,
};

static int memo_init(void){
	//Mostrar carnets
	printk(KERN_ALERT "Carnets: 200412956 201123848\n");
	
	//Create proc file	
	proc_create("memo_200412956", 0, NULL, &memo_proc_fops);

	return 0;
}

static void memo_exit(void){
	//Remover proc entry
	remove_proc_entry("memo_200412956", NULL);
	//Mostrar nombre del curso
	printk(KERN_ALERT "Sistemas Operativos 1\n");
}

module_init(memo_init);
module_exit(memo_exit);

//-------------- metadata -------------------------
MODULE_AUTHOR("Nelson González");
MODULE_DESCRIPTION("[SO1] Moduclo de memoria");
MODULE_LICENSE("GPL");
