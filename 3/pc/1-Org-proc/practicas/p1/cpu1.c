/*	cpu1.c	*/

/* trata correctamente rd en STORE como fuente */

#include "cpu.h"

extern void chivato();
extern void inichivato();
extern char get_instr();
extern void etapa_PreDecode();

IREG inula = {0, 0, 0, NOP, NO_USADO, NO_USADO, NO_USADO, NO_USADO, 0, 0, 0, 0, 0, NO_USADO, NO_USADO, 0, 0};
IREG etapa_Bin, etapa_Pin, etapa_Din, etapa_Ain, etapa_Min, etapa_Ein;
IREG etapa_Bout, etapa_Pout, etapa_Dout, etapa_Aout, etapa_Mout, etapa_Eout;

char carga_B = 1, carga_D = 1, carga_A = 1, carga_M = 1, carga_E = 1, carga_P = 1;

unsigned long int tiempo = 0;
unsigned long int instrucciones = 0;

unsigned long int cortos1=0,cortos2=0,cortos3=0,banco=0;

static unsigned long ccpu = 0, craw = 0, csalto = 0, cpredec = 0;
static unsigned long cfloat = 0;
static unsigned long loads = 0, stores = 0, saltos = 0, floats = 0;

static char tomado;
static FILE * fpout;


void actualizarCortos(int es_corto_usado, unsigned long int* corto, unsigned long int* cortos_usados, unsigned int * regs_anticipados, uint32_t reg_destino){
    if(es_corto_usado){
        *corto += 1;
        regs_anticipados[*cortos_usados] = reg_destino;
        *cortos_usados += 1;
    }
}


/* Funcion que devuelve true */
int noEstaAnticipado(uint32_t reg_destino, uint32_t* regs_anticipados){
    int fin = 0;
    for(int i = 0; i < 3 && !fin; i++){
        fin = regs_anticipados[i] == reg_destino;
    }
	return !fin;
}

void sim(IREG instr)
{
    /* cargamos la instruccion en etapa_Bin  */
    etapa_Bin = instr;
    instrucciones++;
    
    if (etapa_Bin.co == LOAD) loads++;
    if (etapa_Bin.co == STORE) stores++;
    if (etapa_Bin.co == FLOAT) floats++;
    if (etapa_Bin.co == BRCON || etapa_Bin.co == BRINC) saltos++;

    do
    {
        /* etapa Escritura en BR */
        /* nada que simular en esta etapa */	
        etapa_Eout = etapa_Ein;
        carga_E = 1;

        /* etapa Memoria: si no hay problemas la instruccion pasa a E */
        /* nada que simular en esta etapa */	
        etapa_Mout = etapa_Min;
        carga_M = 1;

        /* etapa Alu: si no hay problemas la instruccion pasa a M */
	/* nada que simular en esta etapa */
        etapa_Aout = etapa_Ain;
        carga_A = 1;

	/* etapa Decode: si no hay problemas la instruccion pasa a A */
	/* a partir de pre-decode solo puede haber un rd valido, el rd0 */
        if (etapa_Ain.co == LOAD && 
                  ((etapa_Din.cf0 && etapa_Din.rf0 == etapa_Ain.rd0)
                || (etapa_Din.cf1 && etapa_Din.rf1 == etapa_Ain.rd0)
                || (etapa_Din.cf2 && etapa_Din.rf2 == etapa_Ain.rd0)))
        {
           etapa_Dout = inula;
           craw++;
           carga_D = 0;
        }
        else {
	        unsigned int regs_corto[3] = {-1, -1, -1};
            unsigned long int cortos_usados = 0;

            /* Si op en MEM se usa el RD y este coincide con alguno de los fuentes se actualizan los contadores corto1 y cortos_usados */
            if(etapa_Ain.cd0 && noEstaAnticipado(etapa_Ain.rd0, regs_corto))
	        {
                actualizarCortos((etapa_Din.cf0 && etapa_Din.rf0 == etapa_Ain.rd0), &cortos1, &cortos_usados, regs_corto, etapa_Ain.rd0);
                actualizarCortos((etapa_Din.cf1 && etapa_Din.rf1 == etapa_Ain.rd0), &cortos1, &cortos_usados, regs_corto, etapa_Ain.rd0);   
                actualizarCortos((etapa_Din.cf2 && etapa_Din.rf2 == etapa_Ain.rd0), &cortos1, &cortos_usados, regs_corto, etapa_Ain.rd0);
	        }

            /* Si op en MEM se usa el RD y este coincide con alguno de los fuentes y no se repite el RD se actualizan los contadores corto2 y cortos_usados */
            if(etapa_Min.cd0 && noEstaAnticipado(etapa_Min.rd0, regs_corto))
	        {
                actualizarCortos((etapa_Din.cf0 && etapa_Din.rf0 == etapa_Min.rd0), &cortos2, &cortos_usados, regs_corto, etapa_Min.rd0);
                actualizarCortos((etapa_Din.cf1 && etapa_Din.rf1 == etapa_Min.rd0), &cortos2, &cortos_usados, regs_corto, etapa_Min.rd0);   
                actualizarCortos((etapa_Din.cf2 && etapa_Din.rf2 == etapa_Min.rd0), &cortos2, &cortos_usados, regs_corto, etapa_Min.rd0);
	        }
            
            /* Si op en WB se usa el RD y este coincide con alguno de los fuentes y no se repite el RD se actualizan los contadores corto3 y cortos_usados */
            if(etapa_Ein.cd0 && noEstaAnticipado(etapa_Ein.rd0, regs_corto))
            {
                actualizarCortos((etapa_Din.cf0 && etapa_Din.rf0 == etapa_Ein.rd0), &cortos3, &cortos_usados, regs_corto, etapa_Ein.rd0);
                actualizarCortos((etapa_Din.cf1 && etapa_Din.rf1 == etapa_Ein.rd0), &cortos3, &cortos_usados, regs_corto, etapa_Ein.rd0);   
                actualizarCortos((etapa_Din.cf2 && etapa_Din.rf2 == etapa_Ein.rd0), &cortos3, &cortos_usados, regs_corto, etapa_Ein.rd0);
            }
            
            banco += etapa_Din.cf0 + etapa_Din.cf1 + etapa_Din.cf2 - cortos_usados;
            etapa_Dout = etapa_Din;
            carga_D = 1;
        }

        /* etapa Pre-decodificacion */
	/* desdobla las instrucciones ld/st dobles y con pre/pos incremento/decremento */
	/* ademas, si detecta salto tomado para a la etapa B mediante la variable tomado */

        carga_P = carga_D;
        tomado = 0;
        if (carga_P)
        {
            if ((etapa_Pin.co == BRCON && etapa_Pin.taken == 1) 
                || (etapa_Pin.co == BRINC))
            {
                tomado = 1;
            }

            etapa_PreDecode();
	    if (!carga_P) cpredec++;
        }
        
        /* etapa Busqueda: si no hay problemas la instruccion pasa a D */
        /* nada que simular en esta etapa */	

        if ( tomado ) 
        {
             etapa_Bout = inula;
             carga_B = 0; 
            /* Solo contamos si no hemos parado la carga por algo previamente
               para cuadrar los ciclos de parada */
             if (carga_P) csalto++;
        } else 
        {            
            etapa_Bout = etapa_Bin;
            carga_B = carga_P;
        }
        
        if (carga_B) ccpu++;

        reloj();
    } while(!carga_B);
}

void reloj()
{
    if (carga_P) etapa_Pin = etapa_Bout;
    if (carga_D) etapa_Din = etapa_Pout;
    if (carga_A) etapa_Ain = etapa_Dout;
    if (carga_M) etapa_Min = etapa_Aout;
    if (carga_E) etapa_Ein = etapa_Mout;
    tiempo++;
    //chivato();
}

void inicpu()
{
    int i;
	fpout = fopen("/dev/tty","w");

    etapa_Ein = inula;
    etapa_Eout = inula;
    etapa_Min = inula;
    etapa_Mout = inula;
    etapa_Ain = inula;
    etapa_Aout = inula;
    etapa_Din = inula;
    etapa_Dout = inula;
    etapa_Bout = inula;
    //inichivato();
}

void fincpu()
{
	unsigned long todocpu;

    fprintf(fpout, "CPI: %lu inst. %lu ciclos %2.2f ciclos/inst.\n",
	    instrucciones, tiempo, tiempo/(float)instrucciones);

    todocpu =  ccpu + craw + cfloat + csalto + cpredec;
    fprintf(fpout, "Ciclos CPU: %lu cpu %lu PREDEC %lu RAW %lu FLOAT %lu SALTOS\t(Total: %lu)\n",
	    ccpu, cpredec, craw, cfloat, csalto, todocpu);
    
    fprintf(fpout, "%lu c1; %lu c2; %lu c3; %lu banco registros\n",
		cortos1,
		cortos2,
		cortos3,
		banco);

    int total = cortos1 + cortos2 + cortos3 + banco;
    fprintf(fpout, "%.2f c1; %.2f c2; %.2f c3; %.2f banco registros\n",
	   (double) cortos1 / total * 100,
	   (double) cortos2 / total * 100,
	   (double) cortos3 / total * 100,
	   (double) banco / total * 100);

    fprintf(fpout, "%lu loads %lu stores %lu floats %lu saltos\n", loads, 
        stores, floats, saltos);
}
