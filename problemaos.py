import threading,time,random

capacidad  = 5
maxabelles = 5

class monitor:
    def __init__(self):
        self.buffer = []
        self.mutex = threading.Lock()
        self.potmenjar = threading.Condition(self.mutex)
        self.potproduir = threading.Condition(self.mutex)
       
    
    def produir(self,item):
        with self.mutex:
            while len(self.buffer) == capacidad:
                self.potproduir.wait()
            self.buffer.append(item)
            self.potmenjar.notify()
        
        
    def consumir(self):
        with self.mutex:
            while not self.buffer:
                self.potmenjar.wait()
            item = self.buffer.pop()
            self.potproduir.notify()
            return item
            
def productor(buffer, id):
    for i in range(capacidad):
        item = id + i + random.randint(10,100)
        buffer.produir(item)
        print(f"Productor {id} ha produit: {item}")
        time.sleep(random.uniform(0.2,0.5))
        

def consumidor(buffer):
    for i in range(capacidad):
        item = buffer.consumir()
        print(f"            Consumidor ha menjat: {item}")
        time.sleep(0.2)
    
def main():
    buffer = monitor()
    abelles = []
    for i in range(maxabelles):
        abellaproductora = threading.Thread(target = productor, args = (buffer, i))
        abelles.append(abellaproductora)
    
    os = threading.Thread(target = consumidor, args = (buffer,))
    
    
    for a in abelles:
        a.start()
    os.start()
    
    for b in abelles:
        b.join()
    os.join()
    
if __name__ == "__main__":
    main()