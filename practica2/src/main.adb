with Ada.Text_IO; use Ada.Text_IO;
with def_monitor; use def_monitor;
procedure Main is
   --5 coches y una ambulacia
   THREADS : constant Integer := 6;
   --dato compartido con las tasks
   monitor : puente;
   --especificacion de la tasks
   task type vehiculo is
      entry Start(id : Integer; Ruta_del_norte: boolean);
   end vehiculo;
   --cuerpo de la task
   task body vehiculo is
      nom: Integer;
      del_norte: Boolean;
   begin
      accept Start(id: Integer; Ruta_del_norte: Boolean) do
         --obtenir varibales per als Put_lines
         nom := id;
         del_norte := Ruta_del_norte;
      end Start;
      if nom /= 112 then
         if(del_norte) then
            delay 0.1;--espera per posar-se en marxa
            Put_Line("El cotxe " &Integer'Image(nom)&" esta en ruta en direccion Nord");
            monitor.entrada_del_norte(nom);--arriba a l'entrada del nord
            delay 0.1;--espera per arribar al pont
            monitor.salir_del_norte(nom);--salir de l'entrada i entrar en el pont
            delay 0.1;--espera per sortir del pont
            monitor.salir(nom);--sortir del pont
         else
            delay 0.1;--espera a posar-se en marxa
            Put_Line("El cotxe " &Integer'Image(nom)&" esta en ruta en direccion Sud");
            monitor.entrada_del_sur(nom);--arriba a l'entrada del nord
            delay 0.1;--espera per sortir del pont
            monitor.salir_del_sur(nom);--salir de l'entrada i entrar en el pont
            delay 0.1;--espera per sortir del pont
            monitor.salir(nom);--surt del pont
         end if;
      else
         delay 0.1;--espera per posar-se en marxa
         Put_Line("+++++Ambulancia "&Integer'Image(nom)&" esta en ruta");
         delay 0.1;--espera per arribar al pont
         monitor.ambulancia_es_al_pont(nom);--avisar que la ambulancia hi es al pont
         monitor.entrada_ambulancia(nom);--entrar al pont
         delay 0.1;--espera per travessar el pont
         monitor.salir(nom);
      end if;
   end vehiculo;

   --array de las tasks vehivulo
   type vehiculos is array(1..THREADS) of vehiculo;
   coches        :vehiculos;

   --variable local para obtener la direccion de los coches
   local: Boolean;
begin
   --inici de la simulacio
   Put_Line("***************INICI DE LA SIMULACIO DEL PONT***************");
   --comen�ament de les tasks
   for I in 1..THREADS loop
      local := (I mod 2) = 0;
      coches(I).Start((if I /= 6 then I else 112),local);
   end loop;
end Main;
