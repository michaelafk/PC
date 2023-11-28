package def_monitor is

   protected type puente is
      entry entrada_del_sur();
      entry entrada_del_norte();
      procedure salir();

   private
      puente_esta_vacio : Boolean := True;
   end puente;
end def_monitor;
