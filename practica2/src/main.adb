with Ada.Text_IO; use Ada.Text_IO;
with def_monitor; use def_monitor;
procedure Main is
   --5 coches y una ambulacia
   THREADS : constant Integer := 6;
   --dato compartido con las tasks
   monitor : puente
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
         --hacer algo aqui, todavia no se
         nom := id;
         del_norte := Ruta_del_norte;
      end Start;
      delay 0.1;
      if(del_norte) then
         --viene del norte y por tanto hace los prints que toca del norte
         --junto con las funciones del norte
      else
         --viene del sur y por tanto hace los prints que toca del sur
         --junto con las funciones del sur
      end if;

   end vehiculo;

   coches is array(1..THREADS) of vehiculo;


begin
   --inicio de la simulacion
   Put_Line("INICIO DE LA SIMULACION");
   --comienzo de las tasks
   for I in coches'Range loop
      coches(I).Start(I,(I%2 == 0))
   null;
end Main;
