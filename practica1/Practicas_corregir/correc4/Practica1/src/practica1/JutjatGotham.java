package practica1;

import java.util.concurrent.Semaphore;

/**
 *
 * @author josep
 */
public class JutjatGotham {

    //Variable final que defineix el nombre de sospitosos totals
    final static int NUM_SOSPITOSOS = 20;
    //Array que aglutina els noms dels sospitosos
    final static String[] NOMS = {"Max", "Bella", "Charlie", "Lucy", "Cooper", "Lola",
        "Rocky", "Daisy", "Tucker", "Sadie", "Oliver", "Zoe", "Milo", "Ruby", "Leo",
        "Rosie", "Buddy", "Luna", "Duke", "Molly"};
    //Variable entera que emmagatzema el nombre de sospitosos que fitxen
    int nFitxats = 0;
    //Variable entera que emmagatzema el nombre de sospitosos que declaren
    int nDeclaracions = 0;
    //Variable entera que emmagatzema el nombre de sospitosos que accedeixen a la sala
    int nSospitosos = 0;

    //Aquest semàfor coordina l'arribada dels sospitosos a la sala, quan un
    //sospitós arriba es fitxa i s'allibera el semàfor.
    Semaphore fitxar_sem = new Semaphore(0);
    //Aquest semàfor s'utilitza per que el jutge esperi a tots els sospitosos
    //a que estiguin fitxats, abans de començar amb les declaracions.
    Semaphore iniciDeclaracions_sem = new Semaphore(0);
    //Aquest semàfor espera a que el jutge prengui declaracions a tots els presents
    //abans de rebre la sentència
    Semaphore fiDeclaracions_sem = new Semaphore(0);
    //Aquest semàfor serveix per a que el jutge esperi a que tots els sospitosos
    //ja estiguin sentenciats, per aixi poder finalitzar la simulació 
    Semaphore sentenciar_sem = new Semaphore(0);

    //Mètode main
    public static void main(String[] args) throws Exception {

        System.out.println("\nLa policia de Gotham ha detingut a 20 sospitosos");
        System.out.println("El jutge prendra declaracio als que pugui\n");

        //Objecte jutjat, perque el jutge i els sospitosos pertanyin al mateix objecte
        JutjatGotham jutjat = new JutjatGotham();

        //Iniciam el Jutge
        Thread p_jutge = new Thread(new Jutge(jutjat));
        p_jutge.start();

        //Iniciam els sospitosos
        Thread[] p_sospitosos = new Thread[NUM_SOSPITOSOS];

        for (int i = 0; i < NUM_SOSPITOSOS; i++) {
            p_sospitosos[i] = new Thread(new Sospitos(jutjat, NOMS[i], i));
            p_sospitosos[i].start();
        }

        p_jutge.join();

        for (int i = 0; i < NUM_SOSPITOSOS; i++) {
            p_sospitosos[i].join();
        }
    }

    //Mètode que genera un sleep de n milisegons
    public void sleep(int n) {
        try {
            Thread.sleep(n);
        } catch (InterruptedException e) {
            System.err.println("ERROR: " + e.getMessage());
        }
    }

    //Mètode que realitza la simulació de la entrada i del fitxatge d'un sospitós
    public void fitxarSospitos(String nom) {

        System.out.println("\t" + nom + " entra al jutjat. Sospitosos: " + nSospitosos);
        nSospitosos++;

        sleep(500); //Temps de fitxar d'un sospitós

        System.out.println("\t" + nom + " fitxa. Fitxats: " + nFitxats);
        nFitxats++;
        fitxar_sem.release();
    }

    //Mètode que espera que els sospitosos hagin fitxat
    public void esperaIniciDeclaracions() {
        try {

            fitxar_sem.acquire();

        } catch (InterruptedException e) {
            System.err.println("ERROR: " + e.getMessage());
        }
    }

    //Mètode que allibera un permís per a que puguin declarar els sospitosos
    public void iniciDeclaracions() {
        System.out.println("----> Jutge Dredd: Som a la sala, tanqueu la porta!");

        if (nFitxats != 0) {
            System.out.println("----> Jutge Dredd: Fitxeu als sospitosos presents.");
        }

        iniciDeclaracions_sem.release();
    }

    //Mètode que adquireix permisos del nombre de fitxats que hiha
    public void esperaFiDeclaracions() {
        try {
            System.out.println("----> Jutge Dredd: Preniu declaracio als presents.");
            iniciDeclaracions_sem.acquire(nFitxats);

        } catch (InterruptedException e) {
            System.err.println("ERROR: " + e.getMessage());
        }
    }

    //Mètode que realitza la declaracio del sospitós
    public void prendreDeclaracio(String nom) {
        System.out.println("\t" + nom + " declara. Declaracions: " + nDeclaracions);
        nDeclaracions++;
        fiDeclaracions_sem.release();
    }

    //Mètode que espera la sentècia del jutge per enviarlos a l'asil
    public void esperaSentencies() {
        try {
            fiDeclaracions_sem.acquire(nDeclaracions);

            if (nFitxats != 0) { //Si ha fitxat colcú
                System.out.println("----> Jutge Dredd: Poden abandonar la sala, tots a l'asil!");

            } else {
                System.out.println("----> Jutge Dredd: Si no hiha ningu me'n vaig!");
            }

            System.out.println("----> Jutge Dredd: La justicia descansa, dema prendre declaracio als sospitosos que queden.");
            sentenciar_sem.release();

        } catch (InterruptedException e) {
            System.err.println("ERROR: " + e.getMessage());
        }
    }

    //Mètode que imprimeix el sospitos que s'ha enviat a l'asil
    public void sentenciarSospitos(String nom) {
        System.out.println("\t" + nom + " entra a l'Asil d'Arkham.");
    }

    //Mètode que imprimeix el sospitós que vol declarar, que no ha declarat
    public void citaPerDema(String nom) {
        System.out.println("\t" + nom + ": ¡No es just! Vull declarar. ¡Som inocent!");
    }
}

//Clase Jutge que recull totes les operacions que realitza el procés del jutge
class Jutge implements Runnable {

    private JutjatGotham jutjat; //Objecte jutjat al cual pertanydrà el jutge

    //MÈTODE CONSTRUCTOR
    public Jutge(JutjatGotham j) {
        this.jutjat = j;
    }

    @Override
    public void run() {
        try {
            Thread.sleep(2000); //Temps que tarda en arribar el Jutge
            System.out.println("----> Jutge Dredd: Jo som la llei!");

            Thread.sleep(1000); //Temps que tarda en entrar el jutge

            if (jutjat.nSospitosos == 0) { //Si no ha entrat cap sospitós
                jutjat.esperaSentencies();

            } else if (jutjat.nFitxats == jutjat.NUM_SOSPITOSOS || (jutjat.nFitxats > 0 && jutjat.nSospitosos > 0)) {

                jutjat.iniciDeclaracions();
                jutjat.esperaFiDeclaracions();
                jutjat.esperaSentencies();

            } else {

                jutjat.esperaSentencies();
            }

        } catch (InterruptedException e) {
            System.err.println("ERROR: " + e.getMessage());
        }
    }
}

//Clase Sospitos que recull totes les operacions que realitza un procés sospitós
class Sospitos implements Runnable {

    JutjatGotham jutjat; //Objecte jutjat al cual pertanydrà el sospitos
    String nom; //Nom del sospitós
    int id; //Id del proces del sospitos, per asignarli el nom

    //MÈTODE CONSTRUCTOR
    public Sospitos(JutjatGotham j, String s, int id) {
        this.jutjat = j;
        this.nom = s;
        this.id = id;
    }

    @Override
    public void run() {

        try {

            System.out.println("\t" + nom + ": Som innocent!");
            Thread.sleep(1000); //Temps que tarda el sospitós a entrar a la sala
            jutjat.fitxarSospitos(nom);
            Thread.sleep(1000); //Temps de simulacio de la declaracio
            jutjat.prendreDeclaracio(nom);
            jutjat.sentenciarSospitos(nom);

        } catch (InterruptedException e) {
            System.err.println("ERROR: " + e.getMessage());
        }
    }
}
