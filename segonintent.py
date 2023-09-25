import time
import threading

def threads(etiqueta):
    const_global = 0
    global wantp
    global wantq
    while True:
        if wantp == False:
            while(wantq != False):
                print("proceso 1 esperando")
        elif wantq == False:
            while(wantp != False):
                print("proceso 2 esperando")
        if wantq == False:
            wantp = True
        elif wantp == False:
            wantq = True
        const_global+=1
        if wantq == False:
            wantp = False
        elif wantp == False:
            wantq = False
        print(f"proceso {etiqueta}: {const_global}")

wantp = False
wantq = False
t1 = threading.Thread(target= threads,args=("1",))
t2 = threading.Thread(target= threads,args=("2",))
    
t1.start()
t2.start()
    
t1.join()
t2.join()