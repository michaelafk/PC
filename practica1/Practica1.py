import threading

mutex = threading.Semaphore
MaxSospitosos = 20
SemSos = threading.Semaphore()
SemJut = threading.Semaphore()
counter = 0
def sospitos(id):
    SemSos.acquire()
    SemSos.release()

def jutge():
    SemJut.acquire()
    SemJut.release()
    SemSos.release()

def main():
    sospitosos = []
    for i in range(MaxSospitosos):
        sos = threading.Thread(target=sospitos,args=(i+1,))
        sos.name = f"{i+1}"
        sospitosos.append(sos)
    for s in sospitosos:
        s.start()
        print(f"sospitos {s.name}: \tSom innocent!")
    jut = threading.Thread(target = jutge)
    jut.start()
    print("Jutge Dredd: Jo som la llei!")

    for s in sospitosos:
        s.join()
        print(f"sospitos {s.name} entra a l'Asil d'Arkham")
    jut.join()
    print("Jutge Dredd: La justícia descansa, demà prendré declaració als sospitosos que queden")



if __name__ == "__main__":
    main()