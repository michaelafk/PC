with Ada.Text_Io; use Ada.Text_Io;
with def_monitor; use def_monitor;
with Ada.Numerics.Discrete_Random;

procedure main is

  nSur : integer := 3;
  nNorte : integer := 2;
  NVEHICULOS : constant integer := 5; -- Vehículos en total (menos ambulancia)

  monitor : PuenteMonitor;

  task type vehiculo is
    entry Start (Idx : in integer; direccion : in boolean);
  end vehiculo;
  
  -- Rangos de delay
  type sleepRange is range 1 .. 10; -- Delays de 0,1 a 1 segundo
  package time is new Ada.Numerics.Discrete_Random (sleepRange);
  Generador : time.Generator; -- RNG

  -- Delay de task aleatorio
  procedure sleep is
  begin
    delay Duration (time.Random (Generador) / 10);
  end sleep;

  -- vehiculo task
  task body vehiculo is
      locIdx : integer;
      locDir : boolean;
   begin

      accept Start (Idx : in integer; direccion : boolean) do
         locIdx := Idx;
         locDir := direccion;
         if (locIdx /= 112) then -- Se mira que no sea el id de la ambulancia
            if (locDir) then
               Put_Line("El cotxe" & Integer'Image(locIdx)
                        & " està en ruta en direcció nord!");
            else
               Put_Line("El cotxe" & Integer'Image(locIdx)
                        & " està en ruta en direcció sud!");
            end if;

         else
            Put_Line("L'ambulància" & Integer'Image(locIdx) & " està en ruta");
         end if;
      end Start;

      -- Llega a la cola para entrar al puente
      sleep;
      monitor.espera(locIdx, locDir);

      -- Se comprueba si se trata de la ambulancia
      if (locIdx /= 112) then
      -- Se comprueba la dirección (norte = True, sur = False)
         if (locDir) then
            -- Viene del norte
            monitor.norteLock(locIdx); -- Espera para entrar
            sleep;
            monitor.unlock(locIdx); -- Sale
            sleep;
         else
            -- Viene del sur
            monitor.surLock(locIdx); -- Espera para entrar
            sleep;
            monitor.unlock(locIdx); -- Sale
         end if;
      else
         -- Entra la ambulancia
         monitor.ambulanciaLock(locIdx); -- Espera para entrar
         sleep;
         monitor.unlock(locIdx); -- Sale
      end if;

  end vehiculo;

  -- Array de vehículos (tasks)
  type arrayNVehiculos is array (1..NVEHICULOS + 1) of vehiculo;
  vehiculos : arrayNVehiculos;

begin
  -- Inicialización
   for Idx in 1..NVEHICULOS loop -- Start vehículos
      if (nNorte > 0) then
         vehiculos(Idx).Start(Idx, True);
         nNorte := nNorte - 1;
      else
         vehiculos(Idx).Start(Idx, False);
      end if;

   end loop;
   -- Start ambulancia, se identifica con id = 112
   vehiculos(NVEHICULOS + 1).Start(112, True);

end main;
