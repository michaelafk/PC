import threading,time,random

#SEMAFOROS
SemSos = threading.Semaphore(0) #semaforo de los sospechosos
SemSosE = threading.Semaphore(0) #semaforo de los sospechosos para que vayan acabando
SemJut = threading.Semaphore(0) #semaforo del juez
#SemJutS = threading.Semaphore(0)
#VARIABLES
SosCount = 0
#CONSTANTES
MaxSospitosos = 20
SosInsala = 6
def sospitos(id):
    global SosCount
    print(f"sospitos {id}: \tSom innocent!")
    #time.sleep(random.randint(1,2))
    
    SosCount+=1
    print(f"sospitos {id}: \tha entrat a la sala, suspect {SosCount}")
    if SosCount == SosInsala:
        SemJut.release() #hay 6 sospechosos en sala y por tanto llamamos a el juez para que entre
    SemSosE.acquire()
    if SosCount == SosInsala:
        SemJut.release() #acabar jutge
    print(f"sospitos {id}: \tentra a l'Asil d'Arkham")

def jutge():
    global SosCount
    SemJut.acquire() #Jutge espera
    print("Jutge Dredd: Jo som la llei!")
    for i in range(SosInsala):
        SemSosE.release()
    #SemSosE.release() #los procesos sospitosos quedan liberados
    for i in range(SosInsala):
        SemJut.acquire()#esperar a  que todos acaben
    time.sleep(0.1)
    print("Jutge Dredd: La justícia descansa, demà prendré declaració als sospitosos que queden")
    

def main():
    threads = []
    
    jut = threading.Thread(target = jutge)
    threads.append(jut)
    for i in range(SosInsala):
        sos = threading.Thread(target=sospitos,args=(i+1,))
        sos.name = f"{i+1}"
        threads.append(sos)
        
    for thread in threads:
        thread.start()

    for thread in threads:
        thread.join()



if __name__ == "__main__":
    main()