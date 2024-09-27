Actividad 2 Programación concurrente:

Alumno: Noé Fernández Simmank

El código se encuentra en la carpeta EX3 y esta compuesto por tres archivos:

-tresorer.go
-client.go
Son los dos archivos principales
-go.sum
-go.mod
Son archivos necesarios para el entorno go

Para ejecutar los archivos, desde la consola, ubicado en el directorio EX3 
se ejecuta cada proceso por separado con:
go run .\tresorer.go
go run .\client.go "nombre del cliente"

Ahora se muestran unas simulaciones:
-----------------------------------------------------------------------------
SIMULACIÓN 1
-----------------------------------------------------------------------------
En este caso solo esta el tesorero y un cliente y se puede ver como en la 
primera operación se pide una retirada pero al no haber saldo no se realiza, 
después los siguientes se realizan.
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
EJECUCIÓN TRESORER.GO
PS D:\UNI\concurrent\repe\p3> go run .\tresorer.go
El tresorer és al despatx. El botí mínim és: 15
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
Operació rebuda: -10 del client: Juan
OPERACIÓ NO PERMESA NO HI HA FONS
Balanç: 0
<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
Operació rebuda: 8 del client: Juan
Balanç: 8
<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
Operació rebuda: -2 del client: Juan
Balanç: 6
<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
Operació rebuda: 3 del client: Juan
Balanç: 9
<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
EJECUCIÓN CLIENT.GO "JUAN"
PS D:\UNI\concurrent\repe\p3> go run .\client.go Juan
Hola el meu nom és: Juan
Juan vol fer 4 operacions
Juan operacion 1: -10
Operació sol·licitada...
NO HI HA SALDO
Balanç actual: 0
1----------------------------------------
Juan operacion 2: 8
Operació sol·licitada...
OPERACIÓ CORRECTE
Balanç actual: 8
2----------------------------------------
Juan operacion 3: -2
Operació sol·licitada...
ES FARÀ EL REINTEGRE SI HI HA SALDO
Balanç actual: 6
3----------------------------------------
Juan operacion 4: 3
Operació sol·licitada...
OPERACIÓ CORRECTE
Balanç actual: 9
4----------------------------------------
----------------------------------------------------------------------------
SIMUALCIÓN 2
----------------------------------------------------------------------------
En este caso se van intercalando operaciones de dos clientes diferentes, se
van produciendo las operaciones que son viables. Se puede observar que el
valor del balance en la consola del tesorero es correcto en todo momento, en
cambio en las consolas de los clientes no tiene porque mostrar valores
coherentes por la concurrencia de estos. 
Al llegar al botín mínimo, en este caso 15, el tesorero lo roba y los 
clientes no siguen enviando operaciones. 
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
EJECUCIÓN TRESORER.GO
PS D:\UNI\concurrent\repe\p3> go run .\tresorer.go
El tresorer és al despatx. El botí mínim és: 15
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
Operació rebuda: 9 del client: Juan
Balanç: 9
<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
Operació rebuda: 3 del client: Andreu
Balanç: 12
<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
Operació rebuda: -5 del client: Juan
Balanç: 7
<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
Operació rebuda: -4 del client: Andreu
Balanç: 3
<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
Operació rebuda: -3 del client: Andreu
Balanç: 0
<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
Operació rebuda: 6 del client: Juan
Balanç: 6
<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
Operació rebuda: 6 del client: Juan
Balanç: 12
<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
Operació rebuda: 7 del client: Andreu
Balanç: 19
<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
El tresorer decideix robar el deposit i tancar el despatx
El Tresorer se'n va


Cua 'colaDiposits' esborrada amb èxit
Cua 'colaBalanç' esborrada amb èxit
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
EJECUCIÓN CLIENT.GO "JUAN"
PS D:\UNI\concurrent\repe\p3> go run .\client.go Juan
Hola el meu nom és: Juan
Juan vol fer 2 operacions
Juan operacion 1: 9
Operació sol·licitada...
OPERACIÓ CORRECTE
Balanç actual: 9
1----------------------------------------
Juan operacion 2: -5
Operació sol·licitada...
ES FARÀ EL REINTEGRE SI HI HA SALDO
Balanç actual: 7
2----------------------------------------
PS D:\UNI\concurrent\repe\p3> go run .\client.go Juan
Hola el meu nom és: Juan
Juan vol fer 5 operacions
Juan operacion 1: 6
Operació sol·licitada...
OPERACIÓ CORRECTE
Balanç actual: 6
1----------------------------------------
Juan operacion 2: 6
Operació sol·licitada...
El tresorer a dit: L'oficina acaba de tancar
Jo també me'n vaig idò!
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
EJECUCIÓN CLIENT.GO "ANDREU"
PS D:\UNI\concurrent\repe\p3> go run .\client.go Andreu
Hola el meu nom és: Andreu
Andreu vol fer 3 operacions
Andreu operacion 1: 3
Operació sol·licitada...
OPERACIÓ CORRECTE
Balanç actual: 12
1----------------------------------------
Andreu operacion 2: -4
Operació sol·licitada...
ES FARÀ EL REINTEGRE SI HI HA SALDO
Balanç actual: 3
2----------------------------------------
Andreu operacion 3: -3
Operació sol·licitada...
NO HI HA SALDO
Balanç actual: 0
3----------------------------------------
PS D:\UNI\concurrent\repe\p3> go run .\client.go Andreu
Hola el meu nom és: Andreu
Andreu vol fer 1 operacions
Andreu operacion 1: 7
Operació sol·licitada...
El tresorer a dit: L'oficina acaba de tancar
Jo també me'n vaig idò!
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
