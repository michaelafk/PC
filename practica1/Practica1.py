import threading,time,random

#SEMAFOROS
SemSosInsala = threading.Semaphore(1) #semaforo de los sospechosos
SemSosFitxar = threading.Semaphore(1) #semaforo de los sospechosos para que vayan fichando
SemSosDeclarar = threading.Semaphore(0)
SemJut = threading.Semaphore(0) #semaforo para liberar al juez
Bloqueig = threading.Semaphore(0) #semaforo para la exclusion mutua sobre la variable turn
#SemJutS = threading.Semaphore(0)
#VARIABLES 
SosCountInSala = 0
SosCountFitxats = 0
SosCountDeclaracions = 0
portaOberta = True
JutgeInSala = False
#CONSTANTES
MaxSospitosos = 20
SosInsala = 6
TimeInSos = 0.01
TimeInJut = 0.015
TimeSalaJut = 101
def sospitos(id):
    global SosCountInSala,SosCountFitxats,SosCountDeclaracions,portaOberta
    time.sleep(0.02)
    print(f"sospitos {id}: \tSom innocent!")
    SemSosInsala.acquire()
    if portaOberta == True:
        #SemSosInsala.release
        SosCountInSala+=1
        time.sleep(0.001)
        print(f"sospitos {id}: \tha entrat a la sala, sospitos {SosCountInSala}")
        if JutgeInSala == False:
            SemSosInsala.release()
        
        SemSosFitxar.acquire()
        SosCountFitxats+=1
        time.sleep(0.01)
        print(f"sospitos {id}: \tfitxa, fitxat {SosCountFitxats}")
        if SosCountFitxats == SosCountInSala:
            SemJut.release()
        SemSosFitxar.release()

        SemSosDeclarar.acquire()
        SosCountDeclaracions+=1
        time.sleep(0.01)
        print(f"sospitos {id}: \tdeclara, declaracio {SosCountDeclaracions}")
        if SosCountDeclaracions == SosCountInSala:
            SemJut.release()
        SemSosDeclarar.release()

        Bloqueig.acquire()
        print(f"sospitos {id}: \tentra a l'Asil d'Arkham")
        Bloqueig.release()
    else:
        Bloqueig.acquire()
        print(f"sospitos {id}: \t no es just vull declarar")
        Bloqueig.release()
        SemSosInsala.release()
def jutge():
    global SemSosInsala,SosCountFitxats,SosCountDeclaracions,portaOberta,JutgeInSala
    time.sleep(0.0226)
    print("--->Jutge Dredd: Jo som la llei!")
    time.sleep(0.0105)
    SemSosInsala.acquire
    portaOberta = False
    JutgeInSala = True
    print("--->Jutge Dred: Ja som en la sala, tanqueu sa porta")
    print("--->Jutge Dred: Fitxeu als sospitosos presents")
    SemSosInsala.release()
    SemJut.acquire()#esperar fins que els sospitosos fitxen
    print("--->Jutge Dred: Preniu declaracio als presents")
    SemSosDeclarar.release()
    SemJut.acquire#esperar fins que els sospitosos declarin
    time.sleep(0.1)
    print("Jutge Dredd: Podeu abandonar la sala tots a l'asil!")
    print("Jutge Dredd: La justícia descansa, demà prendré declaració als sospitosos que queden")
    #SemSosInsala.release()
    Bloqueig.release()#donar permis per a que els sospitosos acabin
    

def main():
    threads = []
    
    jut = threading.Thread(target = jutge)
    threads.append(jut)
    for i in range(MaxSospitosos):
        sos = threading.Thread(target=sospitos,args=(i+1,))
        sos.name = f"{i+1}"
        threads.append(sos)
        
    for thread in threads:
        thread.start()

    for thread in threads:
        thread.join()



if __name__ == "__main__":
    main()