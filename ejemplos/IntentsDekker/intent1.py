
import threading

count = 0
turn = 1
def p():
    global count,turn
    print("soy el proceso p")
    while turn!=1:
        print("esperando")
    count+=1
    print("proceso p = %d",count)
    turn = 2
def q():
    global count,turn
    print("soy el proceso q")
    while turn!=2:
        print("esperando")
    count+=1
    print("proceso q = %d",count)
    turn = 1


def main():
    Threads = []
    t1 = threading.Thread(target=p)
    Threads.append(t1)
    t1.start()
    t2 = threading.Thread(target=q)
    Threads.append(t2)
    t2.start()
    for i in Threads:
        i.join()

if __name__ == "__main__":
    main()