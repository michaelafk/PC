import threading
from random import randint, random
from time import sleep

MAYA_ASCII = """
    _  _
   | )/ )
\\\\ |//,' __
(")(_)-"()))=-
   (\\\\
"""

POOH_ASCII = """
    _
 __( )_
(      (o____
 |          |
 |      (__/
   \     /   ___
   /     \  \___/
 /    ^    /     \\
|   |  |__|_HUNNY |
|    \______)____/
 \         /
   \     /_
    |  ( __)
    (____)
"""

MAYA_AND_FRIENDS = 3
MAYA_SLEEP = random() * 5
POOH_SLEEP = random() * 10


class HoneyJar:
    
    def __init__(self, jar_size = 10):
        self.jar = 0
        self.jar_size = jar_size

        self.mutex = threading.Lock()
        self.maya = threading.Condition(self.mutex)
        self.pooh = threading.Condition(self.mutex)

    def is_full(self):
        return self.jar == self.jar_size

    def put_honey(self, bee_id):
        """Producer method

        Blocks Maya (bees) while the jar is full.

        Release Pooh (bear) when a bee puts honey and results into a full jar.
        """
        with self.mutex:
            while self.is_full():
                self.maya.wait()

            self.jar += 1
            print(f'Maya friend {bee_id} has put some honey: {self.jar} / {self.jar_size}')
            print(MAYA_ASCII)

            if self.is_full():
                self.pooh.notify()

    def eat_honey(self):
        """Consumer method

        Blocks Pooh (bear) while the jar is not full.

        After Pooh enters the critical section and eat the honey jar,
        release Maya (bees) currently locked.
        """
        with self.mutex:
            while not self.is_full():
                self.pooh.wait()

            print('Pooh eats hunny UwU')
            print(POOH_ASCII)
            sleep(POOH_SLEEP)
            self.jar = 0

            self.maya.notify()


honey_jar = HoneyJar()


def maya(bee_id):
    """Maya and friends produce honey"""
    global honey_jar

    while True:
        sleep(MAYA_SLEEP)
        honey_jar.put_honey(bee_id)


def pooh():
    """Pooh eats honey"""
    global honey_jar

    while True:
        honey_jar.eat_honey()


def main():
    maya_and_friends_threads = []

    for i in range(MAYA_AND_FRIENDS):
        t = threading.Thread(target=maya, args=(i, ))
        maya_and_friends_threads.append(t)
        t.start()

    pooh_thread = threading.Thread(target=pooh)
    pooh_thread.start()

    for t in maya_and_friends_threads:
        t.join()

    pooh_thread.join()



if __name__ == '__main__':
    main()