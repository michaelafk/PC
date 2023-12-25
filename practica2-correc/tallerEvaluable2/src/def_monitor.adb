package body def_monitor is

   protected body Monitor is
   
    entry cotxeNordLock when (cotxesSud < cotxesNord) and (not pontOcupat) and (not ambuVolAccedir) is
        begin
         pontOcupat := true;  -- ocupa el pont
         cotxesNord := cotxesNord-1; -- decrementa els cotxes que esperen a l'entrada nord
        end cotxeNordLock;

    entry cotxeSudLock when (cotxesNord <= cotxesSud) and (not pontOcupat) and (not ambuVolAccedir) is
        begin
         pontOcupat := true; -- ocupa el pont
         cotxesSud := cotxesSud-1; -- decrementa els cotxes que esperen a l'entrada sud
        end cotxeSudLock;

    entry ambuLock when not pontOcupat is
        begin
            pontOcupat:= true;
            ambuVolAccedir:= false;
            
        end ambuLock;

    procedure vehicleUnlock is
        begin
            pontOcupat:= false;
        end vehicleUnlock;

    procedure ambuVolEntrar is
        begin
            ambuVolAccedir:= true; 
      end ambuVolEntrar;
   
    entry cotxeNordVolEntrar when not arribadaOcupat is
        begin
          arribadaOcupat:= true;
          cotxesNord:= cotxesNord +1; -- incrementa els cotxes que esperen a l'entrada nord
          arribadaOcupat:= false;
        end cotxeNordVolEntrar;
      
    entry cotxeSudVolEntrar when not arribadaOcupat is
       begin
          arribadaOcupat:= true;
          cotxesSud:= cotxesSud+1; -- incrementa els cotxes que esperen a l'entrada sud
          arribadaOcupat:= false;
        end cotxeSudVolEntrar;

    function getCotxesNord return Integer is
        begin
            return cotxesNord;
        end getCotxesNord;

    function getCotxesSud return Integer is
        begin
            return cotxesSud;
        end getCotxesSud;

  end Monitor;  

end def_monitor;
