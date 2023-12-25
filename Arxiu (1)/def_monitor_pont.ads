---------------------------------------------------------------------------------
-- Especificacions del monitor Pont que regula el trànsit
---------------------------------------------------------------------------------

package Def_Monitor_Pont is

   protected type Monitor_Pont is
      
      -- Increment del comptador de cotxes esperant al Nord
      procedure esperaNord (Idx : in Integer);
      -- Bloqueig del cotxe al nord
      entry entraNord (Idx : in Integer);
      -- Increment del comptador de cotxes esperant al Sud
      procedure esperaSud (Idx : in Integer);
      -- Bloqueig del cotxe al Sud
      entry entraSud (Idx : in Integer);
      -- Actualització de qualsevol vehicle quan surt del pont
      procedure sortirPont (Idx : in Integer);
      -- Actualització de l'ambulància quan arriba al pont
      procedure esperaAmbulancia (Idx : in Integer);
      -- Bloqueig de l'ambulància
      entry entraAmbulancia (Idx : in Integer);
      
   private

      Cotxe_pont      : Boolean := False; -- Si hi ha vehicle al pont
      Hi_ha_ambulacia : Boolean := False; -- Si hi ha ambulància esperant
      nCotxesNord     : Integer := 0; -- Nombre de cotxes al Nord
      nCotxesSud      : Integer := 0; -- Nombre de cotxes al Sud
      
   end Monitor_Pont;

end Def_Monitor_Pont;
