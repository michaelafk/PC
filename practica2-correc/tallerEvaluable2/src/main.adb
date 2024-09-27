with Text_Io; use  Text_Io;
with def_monitor; use def_monitor;
-- Marta Gonzalez Juan
-- Segon exercici avaluable: El problema del pont

procedure Main is

   THREADS : constant integer := 6; --- total cotxes
   mon : Monitor;
   ambu_num: integer := 112; --- identificador ambulancia

   task type cotxe is
        entry Start (Idx : in integer);
   end cotxe;

  task body cotxe is
    My_Idx : integer;
    entradaNord : boolean;

    begin
      accept Start (Idx : in integer) do
         My_Idx := Idx;
      end Start;

      if (My_Idx = ambu_num) then  -- si es l'ambulancia

        delay 2.0; ---temps que el vehicle espera a posar-se en marxa
        Put_Line ("L'ambulancia " & My_Idx'img & " est� en ruta");
        delay 6.0; ---temps que el vehicle demora per arribar al pont
        mon.ambuVolEntrar; -- ha arribat al pont, activa la variable ambuVolAccedir
        Put_Line ("+++++Ambul�ncia " & My_Idx'img & " espera per entrar");

        mon.ambuLock;
        Put_Line ("+++++Ambul�ncia " & My_Idx'img & " �s al pont");
        delay 2.0; --- temps que el vehicle demora per travessar el pont
        Put_Line ("---->El vehicle " & My_Idx'img & " surt del pont");
        mon.vehicleUnlock;

      else
        ---comprobar si es entrada nord o sud
        if (My_Idx mod 2 =0) then
            entradaNord:=true;
        else
            entradaNord:= false;
        end if;

        if (entradaNord) then
            delay 1.0; ---temps que el vehicle espera a posar-se en marxa
            Put_Line ("El cotxe " & My_Idx'img & " est� en ruta en direcci� Nord");
            delay 3.0; ---temps que el vehicle demora per arribar al pont
            mon.cotxeNordVolEntrar; -- ha arribat al pont, incrementa el nombre de cotxes que esperen al nord
            Put_Line ("El cotxe " & My_Idx'img & " espera a l'entrada NORD, esperen al NORD: "& mon.getCotxesNord'img);

            mon.cotxeNordLock;
            Put_Line ("El cotxe " & My_Idx'img & " entra al pont, esperen al Nord: "& mon.getCotxesNord'img);
            delay 3.0; --- temps que el vehicle demora per travessar el pont
            Put_Line ("---->El vehicle " & My_Idx'img & " surt del pont");
            mon.vehicleUnlock;


        else

            delay 1.0;  --- temps que el vehicle espera a posar-se en marxa
            Put_Line ("El cotxe " & My_Idx'img & " est� en ruta en direcci� Sud");
            delay 3.0;  --- temps que el vehicle demora per arribar al pont
            mon.cotxeSudVolEntrar; -- ha arribat al pont, incrementa el nombre de cotxes que esperen al sud
            Put_Line ("El cotxe " & My_Idx'img & " espera a l'entrada SUD, esperen al SUD: "& mon.getCotxesSud'img);

            mon.cotxeSudLock;
            Put_Line ("El cotxe " & My_Idx'img & " entra al pont, esperen al Sud: "& mon.getCotxesSud'img);
            delay 3.0;  --- temps que el vehicle demora per travessar el pont
            Put_Line ("---->El vehicle " & My_Idx'img & " surt del pont");
            mon.vehicleUnlock;


        end if;

      end if;

  end cotxe;

  -----
  -- Array de tasques
  -----
  type cotxesArray is array (1..THREADS) of cotxe;
  ca : cotxesArray;
  ambu : cotxe;

begin

  -----
  -- Start les tasques
  -----
  for Idx in 1..THREADS loop
    ca(Idx).Start(Idx);
  end loop;
    ambu.Start(ambu_num);

end Main;

