
with Ada.Text_IO; use Ada.Text_IO;
with Ada.Integer_Text_IO; use Ada.Integer_Text_IO;
with Ada.Calendar; use Ada.Calendar;

--No funciona

procedure Simulacio_Pont is
   protected Pont is
      --entry  per a l'ambulància
      entry Entra_Ambulancia;
      --entry per als cotxes
      entry Entra_Cotxe(ID : Integer; Direccio : String);
      --procedure que controla els contadors dels vehicles als accessos
      procedure Surte(ID : Integer; Direccio : String);
   private
      --Boolean que avisa si l'ambulància està esperant
      Ambulancia_Esperant : Boolean := False;
      --contador de cotxes a l'accés nord
      Cotxe_Nord_Esperant : Integer := 0;
      --contador de cotxes a l'accés sud
      Cotxe_Sud_Esperant : Integer := 0;
      procedure Missatge(Num : Integer; Tipus : String; Direccio : String);
   end Pont;
   --vehicles
   task type Vehicle(ID : Integer; Direccio : String) is
      entry Inicia;
   end Vehicle;

   task body Vehicle is
   begin
      accept Inicia;
      Pont.Missatge("el cotxe ", ID, "està en ruta en direcció", Direccio);
      delay 1.0; -- Simula el temps de travessia del pont
      Pont.Surte(ID, Direccio);
   end;

   task body Pont is
   begin
      --deixa a la ambulancia si veim que està esperant
      accept Entra_Ambulancia when Ambulancia_Esperant do
         Ambulancia_Esperant := True;
         Missatge( "L'ambulancia ", 112, " està en ruta");
      end Entra_Ambulancia;
      --pas dels cotxes
      accept Entra_Cotxe(ID : Integer; Direccio : String) when
         (ID mod 2 = 0 and Cotxe_Nord_Esperant = 0) or
         (ID mod 2 /= 0 and Cotxe_Sud_Esperant = 0) do
         if ID mod 2 = 0 then
            Cotxe_Nord_Esperant := Cotxe_Nord_Esperant - 1;
         else
            Cotxe_Sud_Esperant := Cotxe_Sud_Esperant - 1;
         end if;
         Missatge("el cotxe ", ID, "espera a l'entrada", Direccio, ", esperen al", Direccio );
      end Entra_Cotxe;
      --sortida del pont
      procedure Surte(ID : Integer; Direccio : String) is
      begin
         delay 0.5; -- Simula el temps de sortida del pont
         Missatge("---->El vehicle" , ID, "surt del pont");
         if ID = 112 then
            Ambulancia_Esperant := False;
         else
            if Direccio = "Nord" then
               Cotxe_Nord_Esperant := Cotxe_Nord_Esperant - 1;
            else
               Cotxe_Sud_Esperant := Cotxe_Sud_Esperant - 1;
            end if;
         end if;
      end Surte;

      procedure Missatge(Num : Integer; Tipus : String; Direccio : String) is
      begin
         Put("El vehicle "); Put(Item => Num, Fore => 2, Aft => 0, Exp => 0);
         Put(" "); Put(Tipus); Put(" "); Put(Direccio); Put_Line("!");
      end Missatge;
   end Pont;

   Ambulance : Vehicle(112, "");

   Cotxe1 : Vehicle(1, "Nord");
   Cotxe2 : Vehicle(2, "Sud");
   Cotxe3 : Vehicle(3, "Nord");
   Cotxe4 : Vehicle(4, "Sud");
   Cotxe5 : Vehicle(5, "Nord");
begin
   Ambulance.Inicia;
   Cotxe1.Inicia; Cotxe2.Inicia; Cotxe3.Inicia; Cotxe4.Inicia; Cotxe5.Inicia;
   delay 10.0; -- Espera per permetre que tots els vehicles travessin
end Simulacio_Pont;
