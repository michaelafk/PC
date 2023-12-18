
---------------------------------------------------------------------------------
-- Cos del monitor Pont que regula el trànsit
---------------------------------------------------------------------------------
with Ada.Text_IO; use Ada.Text_IO;

package body Def_Monitor_Pont is

   protected body Monitor_Pont is

      -- Increment del comptador de cotxes esperant al Nord
      procedure esperaNord (Idx : in Integer) is
      begin
         nCotxesNord := nCotxesNord + 1;
         Put_Line ("El cotxe " & Idx'Img & " espera a l'entrada NORD, esperen al NORD: " & nCotxesNord'Img);
      end esperaNord;

      -- Bloqueig del cotxe al nord. El cotxe no podrà entrar si hi ha un altre cotxe, o si hi ha
      -- ambulància esperant, o si hi ha més cotxes esperant a l'altre accés
      -- Decrementa el nombre de cotxes que esperen al Nord i actualitza el cotxe al pont
      entry entraNord (Idx : in Integer) when not Cotxe_pont and not Hi_ha_ambulacia and nCotxesNord >= nCotxesSud is
      begin
         Cotxe_pont := True;
         nCotxesNord := nCotxesNord - 1;
         Put_Line ("El cotxe " & Idx'Img & " entra al pont. Esperen al NORD:" & nCotxesNord'Img);
      end entraNord;

      -- Increment del comptador de cotxes esperant al Sud
      procedure esperaSud (Idx : in Integer) is
      begin
         nCotxesSud := nCotxesSud + 1;
         Put_Line ("     El cotxe " & Idx'Img & " espera a l'entrada SUD, esperen al SUD: " & nCotxesSud'Img);
      end esperaSud;

      -- Bloqueig del cotxe al sud. El cotxe no podrà entrar si hi ha un altre cotxe, o si hi ha
      -- ambulància esperant, o si hi ha més cotxes esperant a l'altre accés
      -- Decrementa el nombre de cotxes que esperen al Sud i actualitza el cotxe al pont
      entry entraSud (Idx : in Integer) when not Cotxe_pont and not Hi_ha_ambulacia and  nCotxesSud >= nCotxesNord is
      begin
         Cotxe_pont := True;
         nCotxesSud := nCotxesSud - 1;
         Put_Line ("     El cotxe " & Idx'Img & " entra al pont. Esperen al SUD:" & nCotxesSud'Img);
      end entraSud;

      -- Actualitza el cotxe al pont
      procedure sortirPont (Idx : in Integer) is
      begin
         Cotxe_pont := False;
         Put_Line("---->El vehicle " & Idx'Img & " surt del pont");
      end sortirPont;
      
      -- Actualitza l'ambulància esperant
      procedure esperaAmbulancia (Idx : in Integer) is
      begin
         Hi_ha_ambulacia := True;
         Put_Line ("+++++Ambulància " & Idx'Img & " espera per entrar");
      end esperaAmbulancia;

      -- Bloqueig de l'ambulància. L'ambulància no podrà entrar si hi ha un cotxe al pont
      -- Actualitza el cotxe al pont i l'ambulàcia esperant
      entry entraAmbulancia (Idx : in Integer) when not Cotxe_pont is
      begin
         Cotxe_pont := True;
         Hi_ha_ambulacia := False;
         Put_Line ("+++++Ambulància " & Idx'Img & " és al pont");
      end entraAmbulancia;
      
   end Monitor_Pont;

end Def_Monitor_Pont;
