#Autor: Antonio Pujol

import threading
import datetime
import random
import time

VEGADES_OBRIR = 10
THREADS = 2
WAIT_SECS = 5.000
want = [0,0]
last = 0
n = 0
timePorta = [0,0]

def thread():
    global n, want, last, timePorta
    sum = 0
    currentTime = ""
    id = int(threading.current_thread().name)
    altre = (id + 1) % THREADS
    print("Entrada {}".format(id))
    for i in range(VEGADES_OBRIR):
        start = time.time()
        time.sleep(random.random()*WAIT_SECS)
        want[id] = 1
        last = id
        while want[altre] == 1 and last == id:
            pass
        # Start SC
        n = n + 1   
        currentTime = datetime.datetime.now().strftime("%H:%M:%S")
        #End SC
        want[id] = 0
        print("Porta {}: {} entrades de: {} Temps: {}".format(id, (i+1), n, currentTime))
        sum = sum + (time.time() - start)
    timePorta[id] = sum/VEGADES_OBRIR

def main():
    threads = []

    for i in range(THREADS):        
        t = threading.Thread(target=thread)
        t.name = i
        threads.append(t)
        t.start()

    for t in threads:
        t.join()

    print("Entrades totals {}".format(n))
    print("Temps mig porta 0: {} segons".format(round(timePorta[0], 3)))
    print("Temps mig porta 1: {} segons".format(round(timePorta[1], 3)))
    

if __name__ == "__main__":
    main()