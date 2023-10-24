import threading,time,random,collections

capacidad  = 5
maxabelles = 5

class monitor:
    def __init__(self,size):
        self.buffer = collections.deque([],size)
        self.mutex = threading.Lock()
        self.potmenjar = threading.Condition(self.mutex)
        self.potproduir = threading.Condition(self.mutex)
       
    
    def produir(self,item):
        with self.mutex:
            while len(self.buffer) == self.buffer.maxlen:
                self.potproduir.wait()
            if len(self.buffer) == self.buffer.maxlen:
                self.potmenjar.notify()
            else:
                self.buffer.append(item)
                
        
        
    def consumir(self):
        with self.mutex:
            while len(self.buffer) == 0:
                self.potmenjar.wait()
            item = self.buffer.popleft()
            if len(self.buffer) == 0:
                self.potproduir.notify_all()
            return item
            
def productor(buffer, id):
    for i in range(capacidad):
        item = random.randint(1,100)
        buffer.produir(item)
        #time.sleep(random.uniform(0.2,0.5))
        print(f"PRODUCTOR {id} ha produit: {item}")
        

def consumidor(buffer):
    for i in range(capacidad):
        item = buffer.consumir()
        #time.sleep(0.2)
        print(f"            CONSUMIDOR ha menjat: {item}")

def main():
    buffer = monitor(capacidad)
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