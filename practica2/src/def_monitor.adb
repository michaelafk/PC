package body def_monitor is
   protected body puente is


      --entrar a la entrada norte del puente
      entry entrada_del_sur(Id : Integer) when coches_en_sur>=0 is
      begin
         coches_en_sur := coches_en_sur + 1;
         Put_Line("     EL cotxe "&Integer'Image(Id)&" espera a l'entrada SUD, esperen al SUD: "&Integer'Image(coches_en_sur));
      end entrada_del_sur;

      --entrar a la entrada sur del puente
      entry entrada_del_norte(Id : Integer) when coches_en_norte>=0 is
      begin
         coches_en_norte := coches_en_norte + 1;
         Put_Line("EL cotxe "&Integer'Image(Id)&" espera a l'entrada Nord, esperen al Nord: "&Integer'Image(coches_en_norte));
      end entrada_del_norte;

      --salir del puente y darle paso a que entre otro
      entry salir_del_sur(Id: Integer) when puente_esta_vacio and coches_en_sur >=coches_en_norte and not ambulancia_esperando is
      begin
         puente_esta_vacio := False;
         coches_en_sur := coches_en_sur - 1;
         Put_Line("     El vehicle "&Integer'Image(Id)&" entra al pont, esperen al SUD: "&Integer'Image(coches_en_sur));
      end salir_del_sur;

      entry salir_del_norte(Id: Integer) when puente_esta_vacio and coches_en_norte > coches_en_sur and not ambulancia_esperando is
      begin
         puente_esta_vacio := False;
         coches_en_norte := coches_en_norte - 1;
         Put_Line("El vehicle "&Integer'Image(Id)&" entra al pont, esperen al Nord: "&Integer'Image(coches_en_norte));
      end salir_del_norte;

      procedure ambulancia_es_al_pont(Id: Integer) is
      begin
         ambulancia_esperando := True;
         Put_Line("+++++Ambulancia "&Integer'Image(Id)&" espera per entrar");
      end ambulancia_es_al_pont;

      entry entrada_ambulancia(Id: Integer) when puente_esta_vacio is
      begin
         puente_esta_vacio := False;
         ambulancia_esperando := False;
         Put_Line("+++++Ambulancia "&Integer'Image(Id)&" es al pont");
      end entrada_ambulancia;
      procedure salir(Id: Integer) is
      begin
         Put_Line("---->El vehicle "&Integer'Image(Id)&" surt del pont");
         puente_esta_vacio:= True;
      end salir;

   end puente;
end def_monitor;
