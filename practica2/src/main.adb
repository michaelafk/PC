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
      if(del_norte) then
         delay 0.1;--esepra a posar-se en marxa
         Put_Line("El cotxe" &Integer'Image(nom)&"est� en ruta en direcci� Nord");
         delay 0.1;--espera per arribar al pont
         monitor.entrada_del_norte(nom);
         delay 0.1;--espera per travesar el pont
         monitor.salir(nom);--surt del pont
      else
         delay 0.1;--espera a posar-se en marxa
         Put_Line("El cotxe" &Integer'Image(nom)&"est� en ruta en direcci� Sud");
         delay 0.1;--espera per arribar al pont
         monitor.entrada_del_sur(nom);
         delay 0.1;--espera per travesar el pont
         monitor.salir(nom);--surt del pont
      end if;

   end vehiculo;

   type vehiculos is array(1..THREADS) of vehiculo;
   coches        :vehiculos;

   local: Boolean;
begin
   --inicio de la simulacion
   Put_Line("INICIO DE LA SIMULACION");
   --comienzo de las tasks
   for I in coches'Range loop
      coches(I).Start(I,(I%2 == 0))
   null;
end Main;
