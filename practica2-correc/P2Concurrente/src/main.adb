with Ada.Text_IO;
use Ada.Text_IO;
--with Coche_Tarea;
--use Coche_Tarea;
with Puente_Monitor;
use Puente_Monitor;

procedure Main is
   N_COCHES : constant Integer := 6;
   ID_AMBULANCIA : constant Integer := 112;
   Puente : Acceso_Puente;
   ---- Especificación de la task ------------
   task type Coche_Tarea is
      entry Start(Id : Integer);
   end Coche_Tarea;

   ---- Cuerpo de la task --------------------
   task body Coche_Tarea is
      --------- N_Coches : Integer := 0;
      Direccion : Character;
      Id_Local : Integer;
   begin
      accept Start(Id : Integer) do
         Id_Local := Id;
      end Start;
      if Id_Local = ID_AMBULANCIA then
         delay 1.0;  -- tiempo para ponerse en marcha
         Put_Line("L'ambulancia" & Integer'Image(Id_Local) & " està en ruta");
         delay 2.5;  -- tiempo para llegar al puente
         Puente.Ambulancia_Espera;  -- set de Ambulancia_Esperando := true y print
         Puente.Entrar_Ambulancia(Id_Local);  -- espera que el puente esté vacío para entrar con prioridad
         delay 3.0;  -- tiempo en el puente
         -- fin de proceso ambulancia
      elsif (Id_Local mod 2) = 0 then  -- id par -> Norte, si no, Sur
         Direccion := 'N';
         delay 1.0;  -- tiempo para ponerse en marcha
         Put_Line("El cotxe " & Integer'Image(Id_Local) & " està en ruta en direcció Nord");
         delay 2.0;  -- tiempo para llegar al puente
         Puente.Incrementar_Coches(Id_Local, Direccion);
         Puente.Entrar_Norte(Id_Local);  -- espera en la cola para entrar al puente
         delay 2.0;  -- tiempo de cruzar puente
      else
         Direccion := 'S';
         delay 1.0;  -- tiempo para ponerse en marcha
         Put_Line("El cotxe " & Integer'Image(Id_Local) & " està en ruta en direcció Sud");
         delay 2.0;  -- tiempo para llegar al puente
         Puente.Incrementar_Coches(Id_Local, Direccion);
         Puente.Entrar_Sur(Id_Local);  -- espera en la cola para entrar al puente
         delay 2.0;  -- tiempo de cruzar puente
      end if;
      Puente.Salir_Puente(Id_Local);
   end Coche_Tarea;

   type Array_Coches is array (1 .. N_COCHES) of Coche_Tarea;
   Coches : Array_Coches;
   Ambulancia : Coche_Tarea;

begin
   Ambulancia.Start(ID_AMBULANCIA);
   for Idx in 1..N_COCHES loop
      Coches(Idx).Start(Idx);
   end loop;
end Main;
