package def_monitor is

   protected type Monitor is
      -- Si el pont no esta ocupat, l'ambulancia no te intencions de entar i el num de cotxes que espera al
      --nord es mayor que els que esperen al sud, entra. Si no, bloqueja al cotxe fins que es compleixin les condicions
      -- Quan entra bloqueja el pont i decrementa el num de cotxes que esperen al nord
      entry cotxeNordLock;
      -- fa el mateix que la entry anterior, pero per als cotxes del sud
      entry cotxeSudLock;
      -- Bloqueja la entrada de l'ambulancia si el pont ja esta ocupat, si no esta ocupat, entra, bloqueja el pont i
      -- desactiva variable ambuVolAccedir
      entry ambuLock; -- per bloquejar l'ambulancia si el pont esta ocupat per un cotxe
      -- allibera el pont una vegada el vehicle ha passat. Es podria haver fet un procedure unlock per a cada una de les entrys, pero
      -- com que la unica acció es modificar la variable pontOcupat basta un
      procedure vehicleUnlock; -- per alliberar el pont
      -- activa ambuVolAccedir indicant que l'ambulancia esta llesta per entrar al pont
      procedure ambuVolEntrar;
      -- incrementen les variables cotxesNord o cotxesSud, depenent de si es un cotxe que ve del nord o del sud,
      -- per evitar problemees amb la  concurrencia, s'emplea la variable arribadaOcupat que s'activa i desactiva abans
      -- i despres de incrementar les variables. Indica que el cotxe ha arribat al pont i esta esperant per entrar
      entry cotxeNordVolEntrar;
      entry cotxeSudVolEntrar;

      -- getters de les variables privades que s'empren al main
      function getCotxesNord return Integer;
      function getCotxesSud return Integer;

    private
        cotxesNord : integer:=0; -- comptador per als cotxes que esperen al nord
        cotxesSud : integer:=0; -- comptador per als cotxes que esperen al sud
        ambuVolAccedir: boolean:= false; -- s'activa quan l'ambulancia ha arribat i vol entrar al pont
        pontOcupat: boolean:= false; -- s'activa quan el pont està ocupat per un cotxe o l'ambulancia
        arribadaOcupat: boolean:= false; -- s'activa cada vegada que un cotxe arriba a l'entrada del pont, ja que es modifica la variable
                                         -- cotxesNord o cotxesSud i pot donar lloc a incongruencies degut a la concurrencia

    end Monitor;

end def_monitor;
