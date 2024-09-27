import threading
salons = [3,3,3]
noms = ["Alexandra","Betlem","Cinta","Adalia","Siara",
        "Feliu","Adela","Margalida","Ester","Isidre","Joradana",
        "Anna"]
Asig = [-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1]
ProcesosTotals = 6 #tant fumadors como no fumadors
class Maitre:
    def __init__(self):
        self.mutex = threading.Lock()
        self.x = threading.Condition(self.mutex)
        self.y = threading.Condition(self.mutex)
        self.Tipus = [-1,-1,-1] #-1 no asignado,1 Fumador,0 No Fumador
    def demanarTaula(self, id):
        Trobat = False
        SaloAsig = -1
        for i in range(len(salons)):
            if salons[i]!=0:
                if self.Tipus[i]==-1 or (id%2 == self.Tipus[i]):
                    self.Tipus[i] = id%2
                    Trobat = True
                    SaloAsig = i
                    salons[i] = salons[i] - 1
                    break
                else:
                    Trobat = False
        #while not Trobat:
            #aux2=""
            #if id%2==0:
                #aux2 = "NOFUMADORS"
            #else:
                #aux2 = "FUMADORS"
            #print(f"No hi ha cap taula disponible per a {noms[id]} a {aux2}")
            #self.x.wait() #esperar a que hi hagi una taula dispnible per al seu tipus
        if Trobat:
            Asig[id] = SaloAsig
            aux = ""
            if self.Tipus[i]==0: 
                aux = "NOFUMADORs"
            else: 
                aux = "FUMADORs"
            print(f"***** El sr./sra. {noms[id]} te taula al salo {SaloAsig} de {aux}")
            if salons[SaloAsig] == 0:
                print(f"salo {SaloAsig} ple")
    
    def deixarTaula(self,id):
        aux = ""
        i = Asig[id]
        if self.Tipus[i] == 1:
            aux = "FUMADOR"
        else:
            aux = "NOFUMADOR"
        salons[i] = salons[i] - 1
        print(f"S'allibera un lloc del salo {i} {aux} Queden {3-salons[i]}")
        print(f"A reveure sr./sra. {noms[id]}")

        

maitre = Maitre()
def Fumadors(fumador_id):
    global maitre
    print(f"Hola el meu nom es {noms[fumador_id]}, voldria dinar i som fumador/a")
    maitre.demanarTaula(fumador_id)
    print(f"{noms[fumador_id]} diu: M'agrada molt el salo {Asig[fumador_id]}")
def NoFumadors(Nofumador_id):
    global maitre
    print(f"Hola el meu nom es {noms[Nofumador_id]}, voldria dinar i no som fumador/a")
    maitre.demanarTaula(Nofumador_id)
    print(f"{noms[Nofumador_id]} diu: M'agrada molt el salo {Asig[Nofumador_id]}")

def main():
    threads = []
    for i in range(ProcesosTotals*2):
        if(i%2 == 1):
            #processos fumadors
            t = threading.Thread(target=Fumadors,args = (i,))
            threads.append(t)
        else:
            #processos no fumadors
            t = threading.Thread(target=NoFumadors,args=(i,))
            threads.append(t)
        t.start()

    for t in threads:
        t.join()

if __name__ == '__main__':
    main()
