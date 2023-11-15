/*
PRÁCTICA 1: El problema del juzgado de Gotham
CREADOR: Claudia Roca Castro
ASIGNATURA: Programación Concurrente
 */
package practica1_progconcurrente;

import java.util.Random;
import java.util.concurrent.Semaphore;
import java.util.logging.Level;
import java.util.logging.Logger;

/**
 *
 * @author croca
 */
public class Practica1_ProgConcurrente {

    // DECLARACIONES
    public static final int NUM_SOSPECHOSOS = 20;
    public static final String[] NOMBRES_SOSPECHOSOS = {"Bane", "Ra's al Ghul", 
        "Riddler", "Catwoman", "Two-Face", "Poison Ivy", "Joker", "Scarecrow", 
        "Mr. Freeze", "Penguin", "Harley Quinn", "Killer Croc", "Jason Todd",
        "Hush", "Hugo Strange", "Talia al Ghul", "Deathstroke", "Clayface",
        "Deadshot", "Mad Hatter"};
    
    // Este semáforo bloquea a los sopechosos que quieren entrar a la sala cuando 
    // el juez ya ha llegado. Además, permite la entrada de uno en uno para garantizar 
    //la exclusión mútua con la variable sospechososEnSala
    public static Semaphore semaforoSala = new Semaphore(1);
    // Este semáforo permite la entrada de uno en uno de los sospechosos a fichar 
    // para garantizarexclusión mútua con la variable sospechososFichados
    public static Semaphore mutexFichar = new Semaphore(1);
    // Este semáforo bloquea al juez hasta que todos los sospechosos hayan fichado
    public static Semaphore todosFichados = new Semaphore(0);
    // Este semáforo bloquea a los sopechosos hasta que el juez permita declarar 
    // Además, permite la entrada de uno en uno para garantizar la exclusión mútua con 
    // la variable sospechososDeclarados
    public static Semaphore puedenDeclarar = new Semaphore(0);
    //  Este semáforo bloquea al juez hasta que todos los sospechosos hayan declarado
    public static Semaphore todosDeclarados = new Semaphore(0);
    // Este semáforo bloquea a los sopechosos hasta que el juez haya dado su veredicto
    public static Semaphore semaforoVeredicto = new Semaphore(0);
    
    public  static int sospechososEnSala = 0;
    public  static int sospechososFichados = 0;
    public  static int sospechososDeclarados = 0;
    
    public static boolean salaCerrada = false;
            
    public static void main(String[] args) throws InterruptedException {
        // DECLARACIONES
        Thread[] sospechosos = new Thread[NUM_SOSPECHOSOS];
        Thread jutge;
        
        // ACCIONES
        System.out.println("La policia de Gotham ha detingut a "+NUM_SOSPECHOSOS+ " sospitosos");
        System.out.println("El jutge prendrà declaració als que pugui");
        
        // Bucle inicialización y generación de sospechosos
        for (int i = 0; i < NUM_SOSPECHOSOS; i++) {
            sospechosos[i] = new Thread(new Sospechoso(NOMBRES_SOSPECHOSOS[i]));
            sospechosos[i].start();
        }
        // Inicialización y generación juez
        jutge = new Thread(new Juez());
        jutge.start();
        
        
        for (int i = 0; i < NUM_SOSPECHOSOS; i++) {
            sospechosos[i].join();
        }
        jutge.join();
    }
    
    public static class Juez implements Runnable{
        public Juez() {}
        
        @Override
        public void run() {
            try {
                // Espera tiempo a que llegue
                retardoAleatorio(500);
                System.out.println("----> Jutge Dredd: Jo som la llei!");
                
                
                // Juez cierra la puerta
                semaforoSala.acquire(); 
                //Espera tiempo a que entre
                Thread.sleep(200);
                System.out.println("----> Jutge Dredd: Som a la sala, tanqueu porta!");
                salaCerrada = true;
                
                
                if (sospechososEnSala == 0) {
                    System.out.println("----> Jutge Dredd: Si no hi ha ningú me'n vaig!");
                } else {
                    System.out.println("----> Jutge Dredd: Fitxeu als sospitosos presents");
                    
                    //Espera a que todos estén fichados
                    todosFichados.acquire(); 
                    System.out.println("----> Jutge Dredd: Preniu declaració als presents");
                    // Permite a todos declarar
                    puedenDeclarar.release(); 
                    
                    // El juez espera hasta que todos hayan declarado
                    todosDeclarados.acquire();
                    System.out.println("----> Judge Dredd: Podeu abandonar la sala tots a l'asil!");
                }
                Thread.sleep(200);
                System.out.println("----> Jutge Dredd: La justícia descansa, demà prendré declaració als sospitosos que queden");
                //Se abre la puerta de nuevo
                semaforoSala.release();
                semaforoVeredicto.release();
                
                
            } catch (InterruptedException ex) {
                Logger.getLogger(Practica1_ProgConcurrente.class.getName()).log(Level.SEVERE, null, ex);
            } 
        }
        
        private void retardoAleatorio(int num) throws InterruptedException {
            Random rand = new Random();
            Thread.sleep(rand.nextInt(num));
        }
    }
    
    public static class Sospechoso implements Runnable{
        private final String nombre;

        public Sospechoso(String nombre) {
            this.nombre = nombre;
        }
        
        @Override
        public void run() {
            
            try {
                // Llegada
                retardoAleatorio(500);
                System.out.println("      " + nombre + ": Som innocent!");
                
                // Si el juez no ha entrado, entra
                semaforoSala.acquire();
                Thread.sleep(200);
                if (salaCerrada) {
                    retardoAleatorio(600);
                    System.out.println("      "  + nombre + ": No és just vull declarar! Som inocent!");
                    semaforoSala.release();
                }else{
                    // Entra en la sala
                    sospechososEnSala++;
                    System.out.println("      " + nombre + " entra al jutjat. Sospitosos: " + sospechososEnSala);
                    // Da permiso de entrar al siguiente
                    semaforoSala.release();
                    
                    // Espera a poder fichar
                    mutexFichar.acquire();
                    Thread.sleep(400);
                    sospechososFichados++;
                    // Avisa al juez si ya todos han fichado
                    System.out.println("      " + nombre + " fitxa. Fitxats: " + sospechososFichados);
                    if (sospechososFichados == sospechososEnSala){
                        todosFichados.release();
                    }
                    // Da permiso al siguiente
                    mutexFichar.release();

                    // Espera a tener permiso de declarar
                    puedenDeclarar.acquire();
                    Thread.sleep(400);
                    sospechososDeclarados++;
                    System.out.println("      " + nombre + " declara. Declaracions: " + sospechososDeclarados);
                    // Avisa al juez si ya todos han declarado
                    if(sospechososDeclarados == sospechososEnSala){
                        todosDeclarados.release();
                    }
                    // Da permiso al siguiente 
                    puedenDeclarar.release();
                    
                    // Esperar al veredicto del juez
                    semaforoVeredicto.acquire();
                    retardoAleatorio(600);
                    System.out.println("      " + nombre +  " entra a l'Asil d'Arkham");
                    // Da permiso al siguiente para declarar
                    semaforoVeredicto.release();
                }
                
            } catch (InterruptedException ex) {
                Logger.getLogger(Practica1_ProgConcurrente.class.getName()).log(Level.SEVERE, null, ex);
            }
        }
        
        private void retardoAleatorio(int num) throws InterruptedException {
            Random rand = new Random();
            Thread.sleep(rand.nextInt(num));
        }
    }
}
