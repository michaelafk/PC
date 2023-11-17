/*
* Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Main.java to edit this template
 */
package practica1;

import java.util.Random;
import java.util.concurrent.Semaphore;
import java.util.logging.Level;
import java.util.logging.Logger;

/**
 *
 * @author Pedro Font
 */
public class Practica1 {

    //Nombre de presos 
    static final int NUMSOSPITOSOS = 20;
    //Nombre de sospitosos en sala, fitxats y declarats
    static int sospitosos = 0;
    static int fitxats = 0;
    static int declarats = 0;
    
    //Semàfors emprats
    static Semaphore declarar = new Semaphore(0);
    static Semaphore acaba = new Semaphore(0);
    static Semaphore porta = new Semaphore(1);
    static Semaphore fitxar = new Semaphore(1);
    static Semaphore bloqueigSospitosos = new Semaphore(1);
    
    public static void main(String[] args) throws InterruptedException{
        new Practica1().inici();

    }
    
    public void inici() throws InterruptedException{
        Thread[] threads = new Thread[NUMSOSPITOSOS];
        String[] nomsSospitosos = {"Two-Face", "Catwoman", "Killer Croc", "Mad Hatter", "Riddler", "Harley Quinn", 
            "Hush", "Clayface", "Mr. Freeze", "Deathstroke", "Poison Ivy", "Bane", "Deadshot", "Joker", "Ra's al Ghul",
            "Penguin", "Talia al Ghul", "Jason Todd", "Scarecrow", "Hugo Strange"};
        
        Thread jutge = new Thread(new Juez());
        jutge.start();
         System.out.println("La policia de Gotham ha detingut a "+ NUMSOSPITOSOS+" sospitosos");
         System.out.println("El jutge prendrà declaració als que pugui");
        for (int i = 0; i<NUMSOSPITOSOS; i++){
            threads[i] = new Thread(new Sospitos(nomsSospitosos[i]));
            threads[i].start();
        }

        for (Thread thread : threads) {
            thread.join();
        }
        jutge.join();
    }
    
    private class Juez implements Runnable{
        Random ran = new Random();
        private final String nombre;
        
        public Juez(){
            nombre = "Jutge Dredd";
        }
        
        @Override
        public void run() {
            try {
                Thread.sleep(ran.nextInt(100));
            } catch (InterruptedException ex) {
                Logger.getLogger(Practica1.class.getName()).log(Level.SEVERE, null, ex);
            }
            System.out.println("----> " + nombre + ": Jo sóc la llei!");
            try{
                Thread.sleep(ran.nextInt(300));
                porta.acquire();
                bloqueigSospitosos.acquire();
                System.out.println("----> " + nombre + ": estic a la sala, tancau la porta!");
                
                if (sospitosos>0){
                    System.out.println("----> " + nombre + ": fitxau als sospitosos presents");
                    while (fitxats<sospitosos){}
                    fitxar.acquire();
                    System.out.println("----> " + nombre + ": preneu declaració a tots els presents");
                    fitxar.release();
                    declarar.release();
                    while (declarats<sospitosos){}
                    declarar.acquire();
                    System.out.println("----> " + nombre + " podeu abandonar la sala, tots a l'Asil!");
                    declarar.release();
                }else{
                    System.out.println("----> " + nombre + ": Si no hi ha ningú, me'n vaig!");
                }
                
                System.out.println("----> " + nombre + ": la justicia descansa, demà prendré declaració als sospitosos que queden");
                porta.release();
                acaba.release();
            } catch (InterruptedException ex) {
                Logger.getLogger(Practica1.class.getName()).log(Level.SEVERE, null, ex);
            }
        }
        
    }
    
      private class Sospitos implements Runnable{
        Random ran = new Random();
        private final String nom;
        
        public Sospitos(String nombre){
            this.nom = nombre;
        }
        
        @Override
        public void run() {
            try {
                Thread.sleep(ran.nextInt(100));
            } catch (InterruptedException ex) {
                Logger.getLogger(Practica1.class.getName()).log(Level.SEVERE, null, ex);
            }
            System.out.println(nom + ": Soc innocent!");
            try{
                Thread.sleep(ran.nextInt(1000));
                porta.acquire();
                if (bloqueigSospitosos.tryAcquire()){
                    bloqueigSospitosos.release();
                    sospitosos++;
                    System.out.println(nom + " entra al jutjat. Sospitosos " + sospitosos);
                    porta.release();

                    fitxar.acquire();
                    fitxats++;
                    System.out.println(nom + " fitxa. Fitxats: " + fitxats);
                    fitxar.release();

                    declarar.acquire();
                    declarats++;
                    System.out.println(nom + " declara. Declaracions: " + declarats);
                    declarar.release();

                    acaba.acquire();
                    Thread.sleep(ran.nextInt(50));
                    System.out.println(nom + " entra a l'Asil d'Arkham");
                    acaba.release();
                }else{
                    porta.release();
                    Thread.sleep(ran.nextInt(50));
                    acaba.acquire();
                    System.out.println(nom + ": No és just, vull ser jutjat");
                    acaba.release();
                }
                
            } catch (InterruptedException ex) {
                Logger.getLogger(Practica1.class.getName()).log(Level.SEVERE, null, ex);
            }
        }
        
    }
    
}