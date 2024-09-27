package body def_monitor is

   protected body PuenteMonitor is

      -- Tránsito del norte --
      -- Entra
      entry norteLock(id: Integer) when (not ambulancia) and
        (not ocupado) and (esperandoNorte >= esperandoSur) is
      begin
         ocupado := True;
         esperandoNorte := esperandoNorte - 1;
         Put_Line("El cotxe" & Integer'Image(id) & " entra al pont"
                  & " Esperen al nord:" & Integer'Image(esperandoNorte));
      end norteLock;

      -- Tránsito del sur --
      -- Entra
      entry surLock(id: Integer) when (not ambulancia) and
        (not ocupado) and (esperandoSur > esperandoNorte) is
      begin
         ocupado := True;
         esperandoSur := esperandoSur - 1;
         Put_Line("El cotxe" & Integer'Image(id) & " entra al pont."
                  & " Esperen al sud: " & Integer'Image(esperandoSur));
      end surLock;


      -- Ambulancia --
      -- Entra
      entry ambulanciaLock(id: Integer) when (not ocupado) is
      begin
         Put_Line("+++++Ambulància" & Integer'Image(id) & " és al pont");
      end ambulanciaLock;


      -- Salida --
      procedure unlock(id: Integer) is
      begin
         Put_Line("----> El vehicle" & Integer'Image(id) & " surt del pont");
         if (id = 112) then
            ambulancia := False;
         end if;
         ocupado := False;
      end unlock;


      -- Llegada al puente --
      -- Espera para incorporarse a la cola
      entry espera(id: Integer; direc: Boolean) when (True) is
      begin
         if (id /= 112) then -- Se mira si es o no ambulancia
            if (direc = True) then -- y la dirección
               esperandoNorte := esperandoNorte + 1;
               Put_Line("El cotxe" & Integer'Image(id) & " espera a l'entrada nord."
                        & " Esperen al nord:" & Integer'Image(esperandoNorte));
            else
               esperandoSur := esperandoSur + 1;
               Put_Line("El cotxe" & Integer'Image(id) & " espera a l'entrada sud."
                        & " Esperen al sud: " & Integer'Image(esperandoSur));
            end if;

         else
            -- Considero que la ambulancia no cuenta para las colas,
            -- no afecta en absoluto y es un poco (muy poco) menos costoso
            ambulancia := True; -- Se activa el flag de prioridad de ambulancia
            Put_Line("+++++Ambulància" & Integer'Image(id) & " espera per entrar");
         end if;
      end espera;


  end PuenteMonitor;

end def_monitor;





