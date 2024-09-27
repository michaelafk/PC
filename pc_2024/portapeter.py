import threading,time,random

def entrada(id):
    #no critical section
    global want
    global entrades
    global Numeroentrades
    global minwaittime
    global maxwaittime
    global last
    print(f"Entrada {id}")
    for i in range(Numeroentrades):
        time.sleep(random.uniform(minwaittime,maxwaittime))
        want[id] = True
        last = id
        pos = (id+1) % 2
        while want[pos] != False | last == pos:
            #esperando
            print()
        #critical section
        entrades[id]+=1
        print(f"Porta {id}: {entrades[id]} entrades de : {i} Temps: {time.ctime()}")
        want[id] = False
        time.sleep(random.uniform(minwaittime,maxwaittime))

entrades = [0,0]
want = [False,False] #wantp = want[0], wantq = want[1]
Numeroentrades = 10
minwaittime = 500/1000
maxwaittime = 5000/1000
last = 0
p0 = threading.Thread(target = entrada,args=[0])
p1 = threading.Thread(target = entrada,args=[1])
p0.start()
p1.start()
p0.join()
p1.join()