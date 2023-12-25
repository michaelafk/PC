with Ada.Text_IO; use Ada.Text_IO;

package def_monitor is
   protected type puente is
      entry entrada_del_sur(Id: Integer);--entry per contabilitzar quants esperen al sur
      entry entrada_del_norte(Id: Integer);--entry per contabilitzar quants esperen al nord
      entry salir_del_sur(Id: Integer);--entry per disminuir els vehicles que hi esperen al sur y per poder pasar al pont
      entry salir_del_norte(Id: Integer);--entry per disminuir els vehicles que hi esperen al nord y per poder pasar al pont
      entry entrada_ambulancia(Id: Integer);--entry per poder donarli prioritat a l'ambulancia
      procedure ambulancia_es_al_pont(Id: Integer);--procedure per poder avisar a les taks que la ambulacia esta esperant
      procedure salir(Id: Integer);--procedure per sortir del pont
   private
      coches_en_sur : Integer := 0;--variable per contabilitzar el nombre de cotxes que hi esperen en el sur
      coches_en_norte: Integer := 0;--variable per contabilitzar el nombre de cotxes que hi esperen en el nord
      puente_esta_vacio : Boolean := True;--variable per poder gestionar el acces al pont
      ambulancia_esperando: Boolean := False;--Variable per poder saber si hi esta esperant l'ambulacia per poder pasar
   end puente;
end def_monitor;
