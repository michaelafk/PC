package body def_monitor is
   protected body puente is


      --entrar a la entrada norte del puente
      entry entrada_del_sur(Id : Integer) when puente_esta_vacio is
      begin
         coches_en_sur := coches_en_sur + 1;
         Put_Line("\tEL cotxe"&Integer'Image(Id)&"espera a l'entrada SUD, esperen al SUD: "&Integer'Image(coches_en_sur));
         if coches_en_sur > coches_en_norte and puente_esta_vacio then
            coches_en_sur := coches_en_sur - 1;
            puente_esta_vacio := False;
            Put_Line("\tEL cotxe"&Integer'Image(Id)&"entra al pont. Esperen al SUD: "&Integer'Image(coches_en_sur));
         end if;
      end entrada_del_sur;

      --entrar a la entrada sur del puente
      entry entrada_del_norte(Id : Integer) when puente_esta_vacio is
      begin
         coches_en_norte := coches_en_norte + 1;
         Put_Line("\tEL cotxe"&Integer'Image(Id)&"espera a l'entrada Nord, esperen al Nord: "&Integer'Image(coches_en_norte));
         if coches_en_norte >= coches_en_sur and puente_esta_vacio then
            coches_en_norte := coches_en_norte - 1;
            puente_esta_vacio := False;
            Put_Line("\tEL cotxe"&Integer'Image(Id)&"entra al pont. Esperen al Nord: "&Integer'Image(coches_en_norte));
         end if;
      end entrada_del_norte;

      --salir del puente y darle paso a que entre otro
      procedure salir(Id: Integer) is
      begin
         puente_esta_vacio := True;
         Put_Line("---->El vehicle "&Integer'Image(Id)&"surt del pont");
      end salir;

   end puente;
end def_monitor;
