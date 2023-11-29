with Ada.Text_IO; use Ada.Text_IO;

package def_monitor is
   protected type puente is
      entry entrada_del_sur(Id: Integer);
      entry entrada_del_norte(Id: Integer);
      procedure salir(Id: Integer);

   private
      coches_en_sur : Integer := 0;
      coches_en_norte: Integer := 0;
      puente_esta_vacio : Boolean := True;
   end puente;
end def_monitor;
