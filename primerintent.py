import time
import threading

const_global = 0
def proceso_p():
    global turno
    k = 1
    while k==1:
        while(turno!=1):
            print("proceso p esperando su turno")
        print("proceso p entrando en seccion critica")
        const_global+=1
        print(const_global)
        print("proceso p finalizando seccion critica")  
        turno = 2
def proceso_q():
    global turno
    k = 1
    while k==1:
        while(turno!=2):
            print("proceso q esperando su turno")
        print("proceso q entrando en seccion critica")
        const_global+=1
        print(const_global)
        print("proceso q finalizando seccion critica")  
        turno = 1

turno = 1
t1 = threading.Thread(target= proceso_p, daemon=True)
t2 = threading.Thread(target= proceso_q, daemon=False)
    
t1.start()
t2.start()
    
t1.join()
t2.join()