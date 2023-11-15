import threading,time,random

#SEMAFOROS
SemSosJutjat = threading.Semaphore(1) #semaforo de los sospechosos y para exclusion mutua
SemSosEnd = threading.Semaphore(0) #semaforo de los sospechosos para que vayan acabando
SemJut = threading.Semaphore(1) #semaforo del juez para despertarlo
SemSosFitxat= threading.Semaphore(0) #semaforo para empezar el fitxat dels sospitosos
SemSosDeclarar = threading.Semaphore(0) #semaforo para empezar declaraciones
#VARIABLES
SosCountInSala = 0
SosCountFitxats = 0
SosCountDeclaracions = 0
open = True
LlistaNoFitxats = []
#CONSTANTES
MaxSospitosos = 20
SosInsala = 6
TimeInSos = 0.01
TimeInJut = 0.015
TimeSalaJut = 101
def sospitos(id):
    global SosCountInSala,SosCountFitxats,SosCountDeclaracions,open,LlistaNoFitxats
    time.sleep(TimeInSos)
    print(f"sospitos {id}: \tSom innocent!")
    
    #comprobar si es pot pasar
    if(open):
        #entrar al jutjat
        SemSosJutjat.acquire()
        SosCountInSala+=1
        print(f"sospitos {id}: \tentra al jutjat. Sospitos: {SosCountInSala}")
        SemSosJutjat.release()
    
        #començar el fitxat
        SemSosFitxat.acquire()
        open = False            #tanquen porta
        SosCountFitxats+=1
        print(f"sospitos {id}: \tfitxa. Fitxats: {SosCountFitxats}")
        SemSosFitxat.release()
        
        #esperar
        SemJut.release()
        
        #començar declaracions
        SemSosDeclarar.acquire()
        SosCountDeclaracions+=1
        print(f"sospitos {id}: \tDeclara. Declaracions: {SosCountDeclaracions}")
        SemSosDeclarar.release()
        
        #esperar
        SemJut.release()
    else:
        LlistaNoFitxats.append(id)
    if id in LlistaNoFitxats:
        print(f"{id}: No és just vull declarar! Som innocent!")
    else:
        print(f"sospitos {id}: \tentra a l'Asil d'Arkham")

def jutge():
    global open,SosCountInSala
    time.sleep(TimeInJut)
    #comença simulacio
    print("Jutge Dredd: Jo som la llei!")
    
    #avisa per a començar a fitxar
    SemSosFitxat.release()
    #espera fins que arribe eljutge
    time.sleep(TimeSalaJut)
    
    #jutge en sala, tancar porta
    print("Jutge Dredd: \tSom a la sala, tanqueu la porta!")
    print("Jutge Dredd: \tFitxeu als sospitosos presents")
    #esperar fins que tots fitxen
    for i in range(SosCountInSala):
        SemJut.acquire()
    print("Judge Dredd: Preniu declaració als presents")
    #esperar fins que tos declarin
    for i in range(SosCountInSala):
        SemJut.acquire()
    time.sleep(0.1)
    print("Jutge Dredd: Podeu abandonar la sala tots a l'asil!")
    print("Jutge Dredd: La justícia descansa, demà prendré declaració als sospitosos que queden")
    

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