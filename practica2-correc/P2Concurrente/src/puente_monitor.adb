with Ada.Text_IO;
use Ada.Text_IO;
with Puente_Monitor;

--- Todos los procedimientos/funciones/entry del código tienen exclusión mutua
-- implementada por el lenguaje por ser un objeto protegido
package body Puente_Monitor is
   protected body Acceso_Puente is

      -- Si hay empate (coches norte = coches sur) he elegido que entran por norte
      entry Entrar_Norte(Id: Integer) when Coches_Norte >= Coches_Sur and not Puente_Ocupado and not Ambulancia_Esperando is
         begin
         Puente_Ocupado := True;
         Coches_Norte := Coches_Norte - 1;
         -- espera y print del puente
         Put_Line("El cotxe " & Integer'Image(Id) & " entra al pont. Esperen al NORD: " & Integer'Image(Coches_Norte));
      end Entrar_Norte;

      entry Entrar_Sur(Id: Integer) when Coches_Sur >= Coches_Norte and not Puente_Ocupado and not Ambulancia_Esperando is
         begin
         Puente_Ocupado := True;
         Coches_Sur := Coches_Sur - 1;
         -- espera y print del puente
         Put_Line("El cotxe " & Integer'Image(Id) & " entra al pont. Esperen al SUD: " & Integer'Image(Coches_Sur));
      end Entrar_Sur;


      entry Entrar_Ambulancia(Id: Integer) when not Puente_Ocupado is
      begin
         Puente_Ocupado := True;
         Ambulancia_Esperando := False;
         Put_Line("+++++Ambulància " & Integer'Image(Id) & " és al pont");
      end Entrar_Ambulancia;

      procedure Salir_Puente(Id : Integer) is
      begin
         if Id = 112 then
            Ambulancia_Esperando := False;
         end if;
         Put_Line("---->El vehicle " & Integer'Image(Id) & " surt del pont");
         Puente_Ocupado := False;
      end Salir_Puente;

      procedure Ambulancia_Espera is
      begin
         Ambulancia_Esperando := True;
         Put_Line("+++++Ambulància 112 espera per entrar");
      end Ambulancia_Espera;

      -- Procedimientos de acceso a los contadores de coches con exclusión mutua
      procedure Incrementar_Coches(Id: Integer; Direccion: Character) is
      begin
         if Direccion = 'N' then
            Coches_Norte := Coches_Norte + 1;
            Put_Line("El cotxe " & Integer'Image(Id) & " espera en la entrada NORD, esperen al NORD: " & Integer'Image(Coches_Norte));
         else
            Coches_Sur := Coches_Sur + 1;
            Put_Line("El cotxe " & Integer'Image(Id) & " espera en la entrada SUD, esperen al SUD: " & Integer'Image(Coches_Sur));
         end if;
      end Incrementar_Coches;

      function Get_Coches(Direccion : Character) return Integer is
      begin
         if Direccion = 'N' then
            return Coches_Norte;
         end if;
         return Coches_Sur;
      end Get_Coches;

   end Acceso_Puente;
end Puente_Monitor;
