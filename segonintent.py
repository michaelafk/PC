import time
import threading

const_global = 0
def proceso_p(etiqueta):
    global wantp
    global wantq
    k = 1
    while k==1:
        print("proceso %s",etiqueta)
        while(wantq):
            print("proceso p esperando su turno")
        wantp = True
        const_global+=1
        print(const_global) 
        wantp = False
def proceso_q(etiqueta):
    global wantp
    global wantq
    k = 1
    while k==1:
        print("proceso %s",etiqueta)
        while(wantp):
            print("proceso q esperando su turno")
        wantq = True
        const_global+=1
        print(const_global)        
        wantq = False

wantp = False
wantq = False
t1 = threading.Thread(target= proceso_p,args=("1",))
t2 = threading.Thread(target= proceso_q,args=("2",))
    
t1.start()
t2.start()
    
t1.join()
t2.join()