#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/sched.h>		// pr_info
#include <linux/sched/signal.h> // for_each_process
#include <linux/proc_fs.h>
#include <linux/seq_file.h>

MODULE_LICENSE("MIT");
MODULE_AUTHOR("Mario Alvarado");
MODULE_DESCRIPTION("Un monitor de CPU sencillo.");
MODULE_VERSION("0.01");

static int imprimirInfo(struct seq_file *m, void *v)
{
	struct task_struct *procesos;
	size_t process_counter = 0;
	for_each_process(procesos)
	{
		char *estado = procesos->state == -1 ? "Inejecutable" : procesos->state == 0 ? "Ejecutable" : "Detenido";
		seq_printf(m, "Nombre: %s\n\tID: %d\n\tEstado: %s(%ld)\n", procesos->comm, procesos->pid, estado, procesos->state);
		++process_counter;
	}
	seq_printf(m, "Número de procesos: %zu\n", process_counter);
	return 0;
}

// Método auxiliar para la lectura
static int abrir_proceso_cpu(struct inode *inode, struct file *file)
{
	return single_open(file, imprimirInfo, NULL);
}

// Estructura para llevar el control del asunto.
static struct file_operations opciones_de_archivo_procedimiento_cpu = {
	.owner = THIS_MODULE,
	.open = abrir_proceso_cpu,
	.read = seq_read,
	.llseek = seq_lseek,
	.release = single_release,
};

int init_module(void)
{
	printk(KERN_INFO "Mairo Alvarado - 201123848\n");
	printk(KERN_INFO "Nelson Gonzalez - 200412956\n");
	proc_create("cpu_201123848", 0, NULL, &opciones_de_archivo_procedimiento_cpu);
	return 0;
}

void cleanup_module(void)
{
	printk(KERN_INFO "Sistemas Operativos 1\n");
}
