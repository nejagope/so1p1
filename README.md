# so1p1
Sistemas Operativos 1 - Proyecto 1 


---------------Sever-------------
Carpeta goserver

Métodos disponibles en la api
/memo
/cpu
/proc
/procs
/procsinfo
/killproc

Descripción de los métodos
Todos los métodos son de tipo GET

/memo
Estado y uso de memoria
No requiere parámetros

/cpu
Estado y uso de cpu
No requiere parámetros

/proc
Detalles y estado de un proceso
Requiere el parámetro pid de tipo entero

/procs
Detalles y estado de todos los procesos
No requiere parámetros

/procsinfo
Conteo de procesos en estado running, sleeping, stoped, zombie y otros estados
No requiere parámetros

/killproc
Mata a un proceso segun su pid. retorna lo mismo que el método /proc
Requiere el parámetro pid de tipo entero