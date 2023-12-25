------ Puente protegido de concurrencia
-- Los procedure/function acceden con exclusión mutua

package Puente_Monitor is
   protected type Acceso_Puente is

      entry Entrar_Norte(Id: Integer);
      entry Entrar_Sur(Id: Integer);
      entry Entrar_Ambulancia(Id: Integer);
      procedure Ambulancia_Espera;
      procedure Incrementar_Coches(Id : Integer; Direccion: Character);
      procedure Salir_Puente(Id : Integer);
      function Get_Coches(Direccion: Character) return Integer;
   private
      Coches_Norte : Integer := 0;
      Coches_Sur : Integer := 0;
      Puente_Ocupado : Boolean := False;
      Ambulancia_Esperando : Boolean := False;
   end Acceso_Puente;
end Puente_Monitor;
