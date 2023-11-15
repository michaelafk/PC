
/**
 * @author Sergio Vega García
 */

import java.util.Random;
import java.util.concurrent.Semaphore;

public class App {

    private static final int CANTIDAD_JUECES = 1;
    private static final int CANTIDAD_SOSPECHOSOS = 20;

    private static final Thread[] HILOS = new Thread[CANTIDAD_JUECES + CANTIDAD_SOSPECHOSOS];

    // Semáforo que va a controlar el acceso a la sala del juzgado
    private static final Semaphore SALA = new Semaphore(1);
    // Semáforo que va a controlar el acceso al contador de sospechosos fichados
    private static final Semaphore FICHAR = new Semaphore(1);
    // Semáforo que va a controlar el acceso al contador de sospechosos declarados
    // Inicialmente sin permisos, una vez todos los sospechosos han fichado el juez
    // libera este semáforo
    private static final Semaphore DECLARAR = new Semaphore(0);
    // Semáforo que le indica al juez si los sospechosos han finalizado lo que esten haciendo
    private static final Semaphore SOSPECHOSOS_LISTOS = new Semaphore(0);

    // Contador de sospechosos en la sala del juzgado, se garantiza la exclusión
    // mutua con el semáforo SALA
    private static int sospechososEnLaSala = 0;
    // Contador de sospechosos fichados, se garantiza la exclusión mutua con el
    // semáforo FICHAR
    private static int sospechososFichados = 0;
    // Contador de sospechosos declarados, se garantiza la exclusión mutua con el
    // semáforo DECLARAR
    private static int sospechososDeclarados = 0;

    // Booleano que controla el estado del juzgado
    // Se garantiza la exclusión mutua con el semáforo SALA
    // Cuando vale false los sospechosos pueden fihcar y declarar cuando se les
    // permita
    // Cuando vale true, al entrar a la sala se les actualizara su sentencia y no
    // podran ni fichar ni declarar y saldrán de la sala
    private static boolean finalizado = false;
    // Boolean que le indica al juez si los sospechosos han finalizado lo que esten haciendo
    private static boolean sospechososListos = false;

    public static void main(String[] args) throws Exception {
        crearHilos();
        System.out.println("La policia de Gotham ha detenido a " + CANTIDAD_SOSPECHOSOS + " sospechosos");
        System.out.println("El juez tomara declaracion a los que pueda");
        iniciarHilos();
        finalizarHilos();
    }

    private static void crearHilos() {
        for (int i = 0; i < HILOS.length; i++) {
            // Cambiar los booleanos para activar (true) o desactivar (false) los tiempos de
            // espera de los hilos
            HILOS[i] = new Thread(i == 0 ? new Juez(false) : new Sospechoso(true));
        }
    }

    private static void iniciarHilos() {
        for (int i = 0; i < HILOS.length; i++) {
            HILOS[i].start();
        }
    }

    private static void finalizarHilos() {
        for (int i = 0; i < HILOS.length; i++) {
            try {
                HILOS[i].join();
            } catch (InterruptedException e) {
                System.err.println("Error intentando finalizar hilo " + HILOS[i]);
                e.printStackTrace();
            }
        }
    }

    private static abstract class Persona implements Runnable {
        protected static enum MENSAJE {
            INICIO, ENTRADA_SALA, SALA_SIN_SOSPECHOSOS, SALA_CON_SOSPECHOSOS, SOSPECHOSOS_FICHADOS,
            SOSPECHOSOS_DECLARADOS, ACABAR, SOSPECHOSO_FICHA, SOSPECHOSO_DECLARA, SOSPECHOSO_VEREDICTO, SOSPECHOSO_FUERA
        }

        protected String Nombre;
        // Variable que controla si las esperas estan activas para esta Persona
        protected final boolean ESPERAS_ACTIVAS;

        public Persona(Boolean esperasActivas) {
            this.Nombre = obtenerNombreAleatorio();
            this.ESPERAS_ACTIVAS = esperasActivas;
        }

        @Override
        public abstract void run();

        protected void esperar(int ms) {
            if (!ESPERAS_ACTIVAS) {
                return;
            }
            try {
                Thread.sleep(ms);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        }

        protected abstract void comunicar(MENSAJE tipo);

        // La persona tarda un tiempo en llegar o no y se presenta
        protected void llegarJuzgado() {
            esperar(new Random().nextInt(1000) + 1000); // Tiempo de llegada
            comunicar(MENSAJE.INICIO);
        }

        private String obtenerNombreAleatorio() {
            String[] nombresDisponibles = new String[] {
                    "Alejandro", "Diego", "Jose", "Luis", "Manuel", "Alberto", "Javier", "Francisco", "Miguel",
                    "Carlos", "Antonio", "David", "Pablo", "Daniel", "Juan", "Sergio", "Andres", "Rafael", "Martin",
                    "Eduardo", "Maria", "Ana", "Laura", "Carmen", "Isabel", "Paula", "Elena", "Raquel", "Marta",
                    "Lucia", "Beatriz", "Adriana", "Sonia", "Rosa", "Sara", "Nuria", "Victoria", "Carolina", "Blanca",
                    "Clara"
            };
            return nombresDisponibles[new Random().nextInt(nombresDisponibles.length)];
        }
    }

    private static class Juez extends Persona {
        public Juez(Boolean esperasActivas) {
            super(esperasActivas);
        }

        @Override
        public void run() {
            // Llegada al juzgado
            llegarJuzgado();

            // region El juez entra a la sala

            esperar(new Random().nextInt(800) + 200); // Tiempo de entrada a la sala
            try {
                // Adquiere el permiso de la sala, no deja entrar a nadie más hasta finalizar el
                // juzgado
                SALA.acquire();
            } catch (InterruptedException e) {
                System.err.println("Error al obtener permiso para acceder a la sala por parte del juez");
                e.printStackTrace();
            }
            comunicar(MENSAJE.ENTRADA_SALA);
            if (sospechososEnLaSala == 0) {
                // El juez anuncia que no hay sospechosos en la sala
                comunicar(MENSAJE.SALA_SIN_SOSPECHOSOS);
                try {
                    // Si no hay ningun sospechoso en la sala se deshabilita la opción de fichar
                    FICHAR.acquire();
                } catch (InterruptedException e) {
                    System.err.println("Error al obtener permiso para deshabilitar los fichados por parte del juez");
                    e.printStackTrace();
                }
            } else {
                // El juez pide que se fichen al resto de sospechosos
                comunicar(MENSAJE.SALA_CON_SOSPECHOSOS);
            }

            // endregion

            // region Si hay sospechosos esperamos a que fichen y declaren

            if (sospechososEnLaSala > 0) {

                // region El juez espera a que todos los sospechosos en la sala fichen

                while (!sospechososListos) {
                    try {
                        SOSPECHOSOS_LISTOS.acquire();
                    } catch (InterruptedException e) {
                        e.printStackTrace();
                    }
                }
                // endregion

                // region El juez permite que los sospechosos declaren

                // Comunica que los sospechosos y pueden declarar
                comunicar(MENSAJE.SOSPECHOSOS_FICHADOS);
                // Libera las declaraciones
                sospechososListos = false;
                DECLARAR.release();

                // endregion

                // region El juez espera a que todos los sospechosos en la sala declaren

                while (!sospechososListos) {
                    try {
                        SOSPECHOSOS_LISTOS.acquire();
                    } catch (InterruptedException e) {
                        e.printStackTrace();
                    }
                }

                // endregion

            }

            // endregion

            // region Fin del juzgado

            // Comunica que todos los sospechosos han de ir al asilo si es que hay
            if (sospechososEnLaSala > 0) {
                comunicar(MENSAJE.SOSPECHOSOS_DECLARADOS);
            }
            // Comunica que acaba por hoy y mañana continua
            comunicar(MENSAJE.ACABAR);
            // Tiempo de salida
            esperar(new Random().nextInt(800) + 200); // Tiempo de salida de la sala
            // Indica que el juicio ha finalizado
            finalizado = true;
            // Libera la sala
            SALA.release();

            // endregion
        }

        @Override
        protected void comunicar(App.Persona.MENSAJE tipo) {
            String mensaje = "----> \tJuez " + this.Nombre + ": ";
            switch (tipo) {
                case INICIO:
                    mensaje += "Yo soy la ley!";
                    break;
                case ENTRADA_SALA:
                    mensaje += "Estoy en la sala, cerrad la puerta!";
                    break;
                case SALA_SIN_SOSPECHOSOS:
                    mensaje += "Si no hay nadie me voy!";
                    break;
                case SALA_CON_SOSPECHOSOS:
                    mensaje += "Fixad a los sospechosos presentes.";
                    break;
                case SOSPECHOSOS_FICHADOS:
                    mensaje += "Tomad declaracion a los presentes.";
                    break;
                case SOSPECHOSOS_DECLARADOS:
                    mensaje += "Podeis abandonar la sala todos al asilo!";
                    break;
                case ACABAR:
                    mensaje += "La justicia descansa, mañana tomare declaracion a los sospechosos restantes.";
                    break;
                default:
                    new Exception("Mensaje invalido para el juez");
                    break;
            }
            System.out.println(mensaje);
        }
    }

    private static class Sospechoso extends Persona {
        private MENSAJE veredicto = MENSAJE.SOSPECHOSO_VEREDICTO;

        public Sospechoso(Boolean esperasActivas) {
            super(esperasActivas);
        }

        @Override
        public void run() {
            llegarJuzgado();

            // region Sospechoso entra a la sala

            esperar(new Random().nextInt(800) + 200); // Tiempo de entrada a la sala
            try {
                // Adquiere el permiso para entrar a la sala
                SALA.acquire();
            } catch (InterruptedException e) {
                System.err.println("Error al obtener permiso para acceder a la sala por parte del sospechoso");
                e.printStackTrace();
            }
            // Si al entrar a la sala el juzgado ha finalizado se actualiza su veredicto al
            // de volver mañana
            // ni se anuncia su entrada a la sala ni se aumenta el numero de sospechosos en
            // la sala
            if (finalizado) {
                veredicto = MENSAJE.SOSPECHOSO_FUERA;
            } else {
                sospechososEnLaSala++;
                comunicar(MENSAJE.ENTRADA_SALA);
            }
            // Libera la sala para que entre otro sospechoso o el juez
            SALA.release();

            // endregion

            // region Fichado y declaración

            // Si el juicio no ha finalizado los sospechosos que hayan entrado a la sala
            // fichan y declaran
            if (!finalizado) {

                // region Fichado

                esperar(new Random().nextInt(400) + 100); // Tiempo para fichar
                try {
                    // Adquiere el permiso para fichar
                    FICHAR.acquire();
                } catch (InterruptedException e) {
                    System.err.println("Error al obtener permiso para fichar");
                    e.printStackTrace();
                }
                // Actualiza el numero de sospechosos que han fichado
                sospechososFichados++;
                // Anuncia que ha fichado
                comunicar(MENSAJE.SOSPECHOSO_FICHA);
                if (sospechososFichados == sospechososEnLaSala) {
                    sospechososListos = true;
                    SOSPECHOSOS_LISTOS.release();
                }
                // Libera el fichado para que el resto de sospechosos puedan fichar
                FICHAR.release();

                // endregion

                // region Declaración

                esperar(new Random().nextInt(400) + 100); // Tiempo para declarar
                try {
                    // Adquiere el permiso para declarar
                    DECLARAR.acquire();
                } catch (InterruptedException e) {
                    System.err.println("Error al obtener permiso para declarar");
                    e.printStackTrace();
                }
                // Actualiza el número de sospechosos que han declarado
                sospechososDeclarados++;
                // Anuncia que ha declarado
                comunicar(MENSAJE.SOSPECHOSO_DECLARA);
                if (sospechososDeclarados == sospechososEnLaSala) {
                    sospechososListos = true;
                    SOSPECHOSOS_LISTOS.release();
                }
                // Libera las declaraciones para que el resto de sospechosos puedan declarar
                DECLARAR.release();

                // endregion
            }

            // endregion

            // region Fin del juicio

            try {
                // El sospechoso adquiere el permiso para salir de la sala
                SALA.acquire();
            } catch (InterruptedException e) {
                System.err.println("Error al obtener permiso para salir de la sala por parte del sospechoso");
                e.printStackTrace();
            }
            // Anuncia su veredicto que dependera de si fue capaz de entrar a la sala antes
            // que el juez o no
            comunicar(veredicto);
            // Libera la sala para que el resto de sospechosos puedan salir
            SALA.release();

            // endregion
        }

        @Override
        protected void comunicar(MENSAJE tipo) {
            String mensaje = "\t" + this.Nombre;
            switch (tipo) {
                case INICIO:
                    mensaje += ": Soy inocente!";
                    break;
                case ENTRADA_SALA:
                    mensaje += " entra al juzgado. Sospechosos " + sospechososEnLaSala;
                    break;
                case SOSPECHOSO_FICHA:
                    mensaje += " fitxa. Fichados: " + sospechososFichados;
                    break;
                case SOSPECHOSO_DECLARA:
                    mensaje += " declara. Declaracionnes: " + sospechososDeclarados;
                    break;
                case SOSPECHOSO_VEREDICTO:
                    mensaje += " entra al asilo de Arkham";
                    break;
                case SOSPECHOSO_FUERA:
                    mensaje += ": No es justo quiero declarar! Soy inocente!";
                    break;
                default:
                    new Exception("Mensaje invalido para el sospechoso");
                    break;
            }
            System.out.println(mensaje);
        }
    }

}
