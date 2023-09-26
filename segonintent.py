import time
import threading

want = [False,False]
def proceso_0(etiqueta):
    global const_global
    while const_global < 10:
        while want[1]!=False:
            print(f" proceso {etiqueta}: esperando")
        want[0] = True
        const_global+=1
        print(f" proceso {etiqueta}: {const_global}")
        want[0] = False
def proceso_1(etiqueta):
    global const_global
    while const_global < 9:
        while want[0] != False:
            print(f" proceso {etiqueta}: esperando")
        want[1] = True
        const_global+=1
        print(f" proceso {etiqueta}: {const_global}")
        want[1] = False

const_global = 0
t1 = threading.Thread(target= proceso_0,args=("1",))
t2 = threading.Thread(target= proceso_1,args=("2",))
    
t1.start()
t2.start()
    
t1.join()
t2.join()
    
t1.join()
t2.join()
