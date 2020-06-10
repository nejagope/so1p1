# cpu_201123848

Un pequeño y sencillo módulo de kernel de linux para monitoreo de procesos.

# Uso

Para cargar el módulo, ejecutar los suguientes comandos como `root`

```sh
# make
# insmod cpu_201123848.ko
# cat /proc/cpu_201123848
```

Para elminiar el módulo, ejecutar el suguiente comando como `root`

```sh
# rmmod cpu_201123848.ko
```

Para hacer una prueba rápida, ejecutar los suguientes comandos como `root`

```sh
# make
# make test
```
