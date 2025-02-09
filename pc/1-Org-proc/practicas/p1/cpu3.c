/*	cpu3.c	*/

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

static unsigned long ccpu = 0, craw = 0, cwaw = 0, cbr = 0, csalto = 0, cpredec = 0;
static unsigned long cfloat = 0;
static unsigned long loads = 0, stores = 0, saltos = 0, floats = 0;

static char tomado;
static FILE * fpout;

/* Ciclos de latencia de operación FLOAT */
static const unsigned int FLOAT_LAT = 5;
static const unsigned int LD_LAT = 2;
static const unsigned int ARITM_LAT = 1;

/* Contador de ciclos en estado ocupado.  */
static unsigned int is_UF_FLOAT_in_use = 0;

/* Vector de bits puerto de escritura */
const unsigned int v_WP_BR_lenght = 7;
static unsigned int v_WP_BR_in_use[7] = {0,0,0,0,0,0,0};


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

        /* Primero se comprueba que no haya un Load productor de uno de los registros fuente */
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
            etapa_Dout = etapa_Din;
            carga_D = 1;
        }

        if(etapa_Din.co == FLOAT && carga_D){

            if(is_UF_FLOAT_in_use){
                // Contabilizamos los ciclos que se para por culpa de UF de FLOAT ocupada.
                cfloat++;
                carga_D = 0;
            }
        }

        if(etapa_Din.cd0 && carga_D){
            unsigned int lat = 0;
            switch(etapa_Din.co){
                case FLOAT:
                    lat = FLOAT_LAT+1;
                    break;
                case LOAD:
                    lat = LD_LAT+1;
                    break;
                default:
                    lat = ARITM_LAT+1;
                    break;
            }
            if(v_WP_BR_in_use[lat]){
                cbr++;
                carga_D = 0;
                etapa_Dout = inula;
            } else{
                if(etapa_Din.co == FLOAT){
                    is_UF_FLOAT_in_use = FLOAT_LAT;
                }
                v_WP_BR_in_use[lat] = 1;
            }
                
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
    if (is_UF_FLOAT_in_use) is_UF_FLOAT_in_use -= 1;    
    for(unsigned int i = 0; i < v_WP_BR_lenght-1; i++){
        v_WP_BR_in_use[i] = v_WP_BR_in_use[i+1];
    }
    v_WP_BR_in_use[v_WP_BR_lenght-1] = 0;
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

    todocpu =  ccpu + cpredec + craw + cwaw + cbr + cfloat + csalto;
    fprintf(fpout, "Ciclos CPU: %lu cpu %lu PREDEC %lu RAW %lu WAW %lu BR %lu FLOAT %lu SALTOS\t(Total: %lu)\n",
	    ccpu, cpredec, craw, cwaw, cbr, cfloat, csalto, todocpu);

    fprintf(fpout, "%lu loads %lu stores %lu floats %lu saltos\n", loads, 
        stores, floats, saltos);
}
