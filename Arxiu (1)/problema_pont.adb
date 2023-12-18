---------------------------------------------------------------------------------
-- Problema del pont amb cotxes i ambulància
-- Descripció de la simulació
-- Hi ha un pont que només pot ser creuat per un vehicle a la vegada, el pont té dos accessos: 
-- l’accés del Nord i l’accés del Sud. Al llarg de la simulació 5 cotxes i una ambulància han de
-- creuar el pont. 
-- •	Al pont només hi pot haver com a màxim un vehicle.
-- •	Quan un vehicle arribi a qualsevol accés entrarà al pont si no hi ha ningú i si està ocupat
--    esperarà.
-- •	Quan un vehicle surti del pont en podrà entrar un altre. Si hi ha l’ambulància esperant 
--    aquesta tindrà preferència damunt tots els altres vehicles i serà la que entrarà al pont 
--    immediatament. Si no hi ha ambulància, tindrà preferència el vehicle que estigui primer a
--    l’accés on hi hagi més vehicles esperant.
-- •	El temps de travessar el pont s’ajustarà per poder veure com s’organitzen les esperes als
--    dos accessos.
-- •	Tots els vehicles travessen el pont una vegada i acaben la simulació.
---------------------------------------------------------------------------------

with Ada.Text_IO;              use Ada.Text_IO;
with Ada.Strings.Unbounded;    use Ada.Strings.Unbounded;
with Ada.Text_IO.Unbounded_IO; use Ada.Text_IO.Unbounded_IO;
with Ada.Numerics.Discrete_Random;
with Def_Monitor_Pont;     use Def_Monitor_Pont;

procedure problema_pont is

   -------------------------------
   -- Constants de la simulació --
   -------------------------------

   Num_Cotxes : constant := 6;
   Num_Ambulancies : constant := 1;

   -----------------------
   -- Nombres aleatoris --
   -----------------------

   Random_Duration : Duration;
   type Custom is range 500..2000;
   package Rand_Cust is new Ada.Numerics.Discrete_Random(Custom);
   use Rand_Cust;
   Seed : Generator;
   Num  : Custom;

   procedure Retard_Aleatori is
   begin
      -- Create the seed for the random number generator
      Reset(Seed);
      -- Generate a random integer from 500 to 2000
      Num := Random(Seed);
      -- Convert Num to a Duration value from 0.5 to 2.0
      Random_Duration := Duration(Num) / 1000.0;
      delay Random_Duration;
   end Retard_Aleatori;

   -------------
   -- Monitor --
   -------------

   monitor : Monitor_Pont;

   ----------------------
   -- Tasca Cotxe --
   ----------------------

   task type Cotxe is
      entry Start (Idx : in Integer);
   end Cotxe;

   task body Cotxe is

      My_Idx : Integer;
      Direccio : Unbounded_String;

   begin
      accept Start (Idx : in Integer) do
         My_Idx := Idx;
         if (My_Idx mod 2 = 0) then
            Direccio := To_Unbounded_String("Nord");
         else 
            Direccio := To_Unbounded_String("Sud");
         end if;
      end Start;
      Retard_Aleatori;
      Put_Line ("El cotxe " & My_Idx'Img & " està en ruta en direcció " & Direccio);
      Retard_Aleatori;
      if (Direccio = "Nord") then
         monitor.esperaNord(My_Idx);
         monitor.entraNord(My_Idx);         
      else
         monitor.esperaSud(My_Idx);
         monitor.entraSud(My_Idx);
      end if;
      Retard_Aleatori;
      monitor.sortirPont(My_Idx);
   end Cotxe;

   ---------------------
   -- Tasca Ambulancia --
   ---------------------

   task type Ambulancia is
      entry Start (Idx : in integer);
   end Ambulancia;

   task body Ambulancia is

      My_Idx : Integer;

   begin
      accept Start (Idx : in Integer) do
         My_Idx := Idx;
      end Start;
      Retard_Aleatori;
      Put_Line ("     L'ambulancia " & My_Idx'Img & " està en ruta");
      Retard_Aleatori;
      monitor.esperaAmbulancia (My_Idx);
      monitor.entraAmbulancia (My_Idx);
      Retard_Aleatori;
      monitor.sortirPont (My_Idx);
   end Ambulancia;

   ---------------
   -- Principal --
   ---------------

   type Cotxes is array (1 .. Num_Cotxes) of Cotxe;

   cot : Cotxes;
   amb : Ambulancia;

begin

   for I in 1 .. Num_Cotxes loop
      cot(I).Start(I);
   end loop;
   amb.Start(112);
end problema_pont;
