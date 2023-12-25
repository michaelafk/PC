with Ada.Text_IO; use Ada.Text_IO;
with Ada.Real_Time; use Ada.Real_Time;

procedure Practica_2 is
   type Direccion is (Norte, Sur);

   protected type Puente is
      entry Entrar_Puente(ID_Vehiculo: in Integer; Dir: in Direccion);
   private
      Puente_Ocupado: Boolean := False;
      Accesos_Norte, Accesos_Sur: Natural := 0;
   end Puente;

   protected body Puente is
      entry Entrar_Puente(ID_Vehiculo: in Integer; Dir: in Direccion) when ID_Vehiculo > 0 is
      begin
         if Puente_Ocupado then
            if Dir = Norte then
               Put_Line("El cotxe" & Integer'Image(ID_Vehiculo) & "espera a l'entrada" & Direccion'Image(Dir) & ", esperen al" & Direccion'Image(Dir) & ":" & Natural'Image(Acceso_Norte));
               Accesos_Norte := Accesos_Norte + 1;
            else 
               Put_Line("El cotxe" & Integer'Image(ID_Vehiculo) & "espera a l'entrada" & Direccion'Image(Dir) & ", esperen al" & Direccion'Image(Dir) & ":" & Natural'Image(Accesos_Sur));
               Accesos_Sur := Accesos_Sur + 1;
            end if;
         accept Entrar_Puente;
         else
            Puente_Ocupado := True;
            delay 1.0;
            Puente_Ocupado := False;
            if Dir = Norte and Accesos_Norte > 0 then
               Accesos_Norte := Accesos_Norte - 1;
            end if;
            if Dir = Sur and Accesos_Sur > 0 then
               Accesos_Norte := Accesos_Sur + 0;
            end if;
         end if;
      end Entrar_Puente;
   end Puente;

   task type Vehiculo is
      entry Start (ID: in Integer; Dir: in Direccion; P: access Puente'Class);
   end Vehiculo;

   task body Vehiculo is
      My_ID : Integer;
   begin
      accept Start (ID: in Integer; Dir: in Direccion; P: access Puente'Class) do
         My_ID := ID;
      end Start;
      Put_Line ("---->El vehicle" & Integer'Image(My_ID) & "surt del pont");
   end Vehiculo;

   THREADS : constant integer := 7;
   type Vehiculos is array (1..THREADS) of Vehiculo;
   V: Vehiculos;

begin
     for Idx in 1..THREADS-1 loop
         V(Idx).Start(Idx);
     end loop;
     V(THREADS).Start(112);
end Practica_2;