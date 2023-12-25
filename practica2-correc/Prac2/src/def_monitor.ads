with Ada.Text_IO; use Ada.Text_IO;
package def_monitor is

   protected type PuenteMonitor is

      entry norteLock(id: Integer);
      entry surLock(id: Integer);
      entry ambulanciaLock(id: Integer);
      procedure unlock(id: Integer);
      entry espera(id: Integer; direc: Boolean);


    private
      esperandoNorte : Integer := 0;
      esperandoSur : Integer := 0;
      ambulancia : Boolean := False; -- True si la ambulancia está en cola
      ocupado : Boolean := False; -- Controla la concurrencia en el puente
      llegando : Boolean := False; -- Controla la concurrencia en las colas

    end PuenteMonitor;

end def_monitor;
