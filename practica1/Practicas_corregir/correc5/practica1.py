
import threading
import time
import random

#SIMULACIONES (solo una puede valer 1)
SIMULACION1 = 0
SIMULACION2 = 0
SIMULACION3 = 0

SOSPECHOSOS = 20

nombres = ["Deadshot","Harley Quinn","Penguin","Riddler","Bane","Talia al Ghul","Ra's al Ghul","Hugo Strange","Killer Croc", "Catwoman",
           "Poison Ivy","Mr. Freeze","Jason Todd","Hush","Joker","Deathstroke","Mad Hatter","Two-Face","Scarecrow","Clayface"]


sala = threading.Semaphore(1)
fitxa = threading.Semaphore(1)
declara = threading.Semaphore(0)

cant_sospechosos = 0
fitxats = 0
declarats = 0
justicia_descansa = False

class Sospechoso(threading.Thread):
    
    
    def __init__(self, name):
        super(Sospechoso, self).__init__()
        self.__name = name
        
        
        
    def run(self):
        
        global cant_sospechosos, fitxats, declarats
        if SIMULACION2:time.sleep(3)
        elif SIMULACION1 or SIMULACION3: time.sleep(random.randrange(5))
        
        print(self.__name + ": Som innocent!")
        
        
        sala.acquire() #sospechoso quiere entrar en la sala
        
        if not justicia_descansa: #si la justicia no descansa entra en la sala, ficha y declara
            time.sleep(random.randrange(5))  
            cant_sospechosos += 1 #modifica variable global de cantidad de sospechosos en la sala
            print(f"{self.__name} entra al jutjat. Sospitosos: {cant_sospechosos}")
            
            sala.release() #deja entrar a otro sospechoso en la sala
            
            fitxa.acquire() #sospechoso quiere fichar
            time.sleep(random.randrange(5))
            fitxats += 1 #modifica variable global de cantidad de fichados
            print(f"{self.__name} fitxa. Fitxats: {fitxats}")
            fitxa.release() #deja fichar a otro sospechoso
            
            declara.acquire() #sospechoso quiere declarar
            declarats += 1 #modifica variable global de cantidad de declaracions
            print(f"{self.__name}  declara. Declaracions: {declarats}")
            declara.release() #deja declarar a otro sospechoso
            
            while not justicia_descansa:
                pass
            
            time.sleep(random.randrange(5)) 
            print(f"{self.__name}  entra a l'Asil d'Arkham")
                    
        else: #si la justicia descansa se queja
            print(f"{self.__name}: No és just vull declarar! Som innocent!")        
            sala.release()
    
        
        
            

        


class Juez(threading.Thread):
    
    def __init__(self):
        super(Juez, self).__init__()
        
    
    def run(self):
        global justicia_descansa
        if SIMULACION1: time.sleep(2)
        elif SIMULACION3: time.sleep(5)
        print("----> Jutge Dredd: Jo som la llei!")
        
        sala.acquire() #el juez quiere entrar en la sala
        if SIMULACION1 or SIMULACION3: time.sleep(random.randrange(5))
        print("----> Jutge Dredd: Som a la sala, tanqueu la porta!")
        
        if cant_sospechosos == 0: #en el caso que el juez entre antes que cualquier sospechoso
            print("----> Jutge Dredd: Si no hi ha ningú me'n vaig!")
           
        else: #si hay sospechosos en la sala
            print("----> Jutge Dredd: Fitxeu als sospitosos presents")
            while not fitxats == cant_sospechosos: #espera a que todos los sospechosos fichen
                pass
            print("----> Jutge Dredd: Preniu declaració als presents")
            declara.release() #permite que el primer sospechoso declare
            while not declarats == cant_sospechosos: #espera a que todos los sospechosos declaren
                pass
            print("----> Jutge Dredd: Podeu abandonar la sala tots a l'asil!")
        
        print("----> Jutge Dredd: La justícia descansa, demà prendré declaració als sospitosos que queden")
        time.sleep(random.randrange(5))
        justicia_descansa = True             
        sala.release() #el juez sale de la sala
        
        



def main():
    global SIMULACION1, SIMULACION2, SIMULACION3
    
    simulacion = int(input("Elige simulación (1, 2 o 3):"))
    while simulacion not in {1,2,3}:
        print("\033[91mError. La simulación solo puede ser 1, 2 o 3\033[0m")
        simulacion = int(input("Elige simulación (1, 2 o 3):"))
    if simulacion == 1: SIMULACION1 = 1
    elif simulacion == 2: SIMULACION2 = 1
    elif simulacion == 2: SIMULACION3 = 1
 
    threads = []
    

    for i in range(SOSPECHOSOS):
        threads.append(Sospechoso(nombres[i]))
        
    threads.append(Juez())    
    for i in range(len(threads)):
        threads[i].start()
        
        
    for i in range(len(threads)):
        threads[i].join()
        
    
    
        


if __name__ == "__main__":
    main()