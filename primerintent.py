import time
import threading

def proceso_0(etiqueta):
    global const_global
    global turno
    while True:
        while turno !=1:
            print(f" proceso {etiqueta}: esperando")
        const_global+=1
        print(f" proceso {etiqueta}: {const_global}")
        turno = 2
def proceso_1(etiqueta):
    global const_global
    global turno
    while True:
        while turno !=2:
            print(f" proceso {etiqueta}: esperando")
        const_global+=1
        print(f" proceso {etiqueta}: {const_global}")
        turno = 1

const_global = 0
turno = 1
t1 = threading.Thread(target= proceso_0,args=("1",))
t2 = threading.Thread(target= proceso_1,args=("2",))
    
t1.start()
t2.start()
    
t1.join()
t2.join()
