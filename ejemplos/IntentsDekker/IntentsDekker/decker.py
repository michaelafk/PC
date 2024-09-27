import threading
THREADS = 2
wantp = False
wantq = False
turn = 1
n = 0
def p():
    global wantp,wantq,turn,n
    #seccion no critica
    wantp = True
    for idx in range(9):
        print(idx)
        while wantq:
            if turn == 2:
                wantp = False
                while turn ==2:
                    print("Proceso p esperando\n")
                wantp = True
        n+=1
        turn = 2
        wantp = False
def q():
    global wantp,wantq,turn,n
    #seccion no critica
    wantq = True
    for idx1 in range(9):
        print(idx1)
        while wantp:
            if turn == 1:
                wantq = False
                while turn ==1:
                    print("Proceso q esperando\n")
                wantq = True
        n+=1
        turn = 1
        wantq = False

def main():
    threads = []
    t1 = threading.Thread(target=p)
    threads.append(t1)
    t1.start()
    t2 = threading.Thread(target=q)
    threads.append(t2)
    t2.start()
    
    for t in threads:
        t.join()
    print(n)

if __name__ == "__main__":
    main()