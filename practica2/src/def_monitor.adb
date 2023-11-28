package body def_monitor is
   protected body puente is


      --entrar a la entrada norte del puente
      entry entrada_del_sur() when puente_esta_vacio is
      begin
         puente_esta_vacio := False;
      end entrada_del_sur;

      --entrar a la entrada sur del puente
      entry entrada_del_norte() when puente_esta_vacio is
      begin
         puente_esta_vacio := False;
      end entrada_del_norte;

      --salir del puente y darle paso a que entre otro
      procedure salir() is
      begin
         puente_esta_vacio := True;
      end salir;

   end puente;
end def_monitor;
