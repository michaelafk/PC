import time
import threading

def threads(etiqueta):
    global turno
    const_global = 0
    while True:
        time.sleep(1)
        if turno == 1:
            while(turno !=1):
                print("proceso 1 esperando")
        else:
            while(turno !=2):
                print("proceso 2 esperando")
        const_global+=1
        print(f"proceso {etiqueta}: {const_global}")
        if turno == 1:
            turno = 2
        else:
            turno = 1

turno = 1
t1 = threading.Thread(target= threads,args=("1",))
t2 = threading.Thread(target= threads,args=("2",))
    
t1.start()
t2.start()
    
t1.join()
t2.join()