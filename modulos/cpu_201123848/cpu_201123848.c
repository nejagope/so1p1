#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/sched.h>	// pr_info
#include <linux/sched/signal.h>	// for_each_process

MODULE_LICENSE("MIT");
MODULE_AUTHOR("Mario Alvarado");
MODULE_DESCRIPTION("Un monitor de CPU sencillo.");
MODULE_VERSION("0.01");

void imprimirInfo(void)
{
	struct task_struct* procesos;
	size_t process_counter = 0;
	for_each_process(procesos) {
		char* estado = procesos->state == -1 ? "Inejecutable" : procesos->state == 0 ? "Ejecutable" : "Detenido";
		printk(KERN_INFO "Nombre: %s\n\tID: %d\n\tEstado: %s(%ld)\n", procesos->comm, procesos->pid, estado, procesos->state);
		++process_counter;
	}
	printk(KERN_INFO "Número de procesos: %zu\n", process_counter);
}

int init_module(void)
{
	printk(KERN_INFO "Mairo Alvarado - 201123848\n");
	printk(KERN_INFO "Nelson Gonzalez - 200412956\n");
	
	imprimirInfo();
	return 0;
}

void cleanup_module(void)
{
	printk(KERN_INFO "Sistemas Operativos 1\n");
}
