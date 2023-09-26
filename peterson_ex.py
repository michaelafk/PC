import threading

def entrada_0():
    #no critical section
    global want
    global count
    want[0] = True
    last = 1
    while want[1] != False | last == 2:
        #esperando
        print()
    #critical section
    count+=1
    want[0] = False 

def entrada_1():
    #no critical section
    global want
    global count
    want[1] = True 
    last = 2
    while want[0] != False | last == 1:
        #esperando
        print()
    #critical section
    count+=1
    want[1] = False 

count = 0
want = [False,False] #wantp = want[0], wantq = want[1]
