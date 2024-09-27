import threading,time

def entrada_0():
    #no critical section
    global want
    global entries0
    global Numeroentrades
    global waittime
    global waittime2
    print("Entrada 0")
    for i in range(Numeroentrades):
        time.sleep(waittime)
        want[0] = True
        last = 1
        while want[1] != False | last == 2:
            #esperando
            print()
        #critical section
        print("Porta 0:")
        entries0+=1
        print(f"{entries0} entrades de : {i} Temps: {time.c}")
        want[0] = False
        time.sleep(waittime2)

def entrada_1():
    #no critical section
    global want
    global entries1
    global Numeroentrades
    global waittime
    global waittime2
    print("Entrada 1")
    for i in range(Numeroentrades):
        time.sleep(waittime)
        print("")
        want[1] = True 
        last = 2
        while want[0] != False | last == 1:
            #esperando
            print()
        #critical section
        print("Porta 1: ")
        entries1+=1
        print(f"{entries1} entrades de: {i} Temps: x")
        want[1] = False
        time.sleep(waittime2)

entries0 = 0
entries1 = 0
want = [False,False] #wantp = want[0], wantq = want[1]
Numeroentrades = 20
waittime = 500/1000
waittime2 = 200/1000
p0 = threading.Thread(target = entrada_0)
p1 = threading.Thread(target = entrada_1)
p0.start()
p1.start()
p0.join()
p1.join()